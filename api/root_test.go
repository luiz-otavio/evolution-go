package evolution

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/liyue201/goqr"
	"github.com/mdp/qrterminal/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

func cleanupContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func InitializePostgres(t *testing.T, opts ...testcontainers.ContainerCustomizer) (*postgres.PostgresContainer, string, func(), error) {
	ctx := t.Context()
	containerOpts := []testcontainers.ContainerCustomizer{
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60 * time.Second),
		),
	}

	containerOpts = append(containerOpts, opts...)

	pg, err := postgres.Run(
		ctx,
		"postgres:15-alpine",
		containerOpts...,
	)

	if err != nil {
		return nil, "", nil, err
	}

	connectionString, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, "", nil, err
	}

	cleanup := func() {
		cleanupCtx, cancel := cleanupContext()
		defer cancel()

		if err := pg.Terminate(cleanupCtx); err != nil {
			t.Fatalf("Failed to terminate Postgres container: %v", err)
		}
	}

	return pg, connectionString, cleanup, nil
}

func Test_Init_Postgres(t *testing.T) {
	_, connectionString, cleanup, err := InitializePostgres(t)
	if err != nil {
		t.Fatalf("Failed to initialize Postgres: %v", err)
	}
	defer cleanup()

	if connectionString == "" {
		t.Fatalf("Expected a valid connection string, got empty")
	}
}

func InitializeRedis(t *testing.T, opts ...testcontainers.ContainerCustomizer) (*redis.RedisContainer, string, func(), error) {
	ctx := t.Context()
	containerOpts := []testcontainers.ContainerCustomizer{
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30 * time.Second),
		),
	}

	containerOpts = append(containerOpts, opts...)

	redisContainer, err := redis.Run(
		ctx,
		"redis:7-alpine",
		containerOpts...,
	)

	if err != nil {
		return nil, "", nil, err
	}

	connectionString, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	cleanup := func() {
		cleanupCtx, cancel := cleanupContext()
		defer cancel()

		if err := redisContainer.Terminate(cleanupCtx); err != nil {
			t.Fatalf("Failed to terminate Redis container: %v", err)
		}
	}

	return redisContainer, connectionString, cleanup, nil
}

func Test_Init_Redis(t *testing.T) {
	_, connectionString, cleanup, err := InitializeRedis(t)
	if err != nil {
		t.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer cleanup()

	if connectionString == "" {
		t.Fatalf("Expected a valid connection string, got empty")
	}
}

func Initialize_Evolution_API(t *testing.T, env map[string]string, opts ...testcontainers.ContainerCustomizer) (testcontainers.Container, string, func(), error) {
	ctx := t.Context()
	containerOpts := []testcontainers.ContainerCustomizer{
		testcontainers.WithExposedPorts("8080/tcp"),
		testcontainers.WithEnv(env),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("8080/tcp").WithStartupTimeout(60 * time.Second),
		),
	}

	containerOpts = append(containerOpts, opts...)

	container, err := testcontainers.Run(
		ctx,
		"evoapicloud/evolution-api:latest",
		containerOpts...,
	)

	if err != nil {
		return nil, "", nil, err
	}

	port, err := container.MappedPort(ctx, "8080/tcp")
	if err != nil {
		return nil, "", nil, err
	}

	cleanup := func() {
		cleanupCtx, cancel := cleanupContext()
		defer cancel()

		if err := container.Terminate(cleanupCtx); err != nil {
			t.Fatalf("Failed to terminate Evolution API container: %v", err)
		}
	}

	return container, "http://localhost:" + port.Port(), cleanup, nil
}

func Test_Init_Evolution_API(t *testing.T) {
	sharedNetwork, err := network.New(t.Context())
	if err != nil {
		t.Fatalf("Failed to initialize docker network: %v", err)
	}

	defer func() {
		cleanupCtx, cancel := cleanupContext()
		defer cancel()

		if err := sharedNetwork.Remove(cleanupCtx); err != nil {
			t.Fatalf("Failed to remove docker network: %v", err)
		}
	}()

	_, _, postgresCleanup, err := InitializePostgres(
		t,
		network.WithNetwork([]string{"postgres"}, sharedNetwork),
	)
	if err != nil {
		t.Fatalf("Failed to initialize Postgres: %v", err)
	}
	defer postgresCleanup()

	_, _, redisCleanup, err := InitializeRedis(t, network.WithNetwork([]string{"redis"}, sharedNetwork))
	if err != nil {
		t.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisCleanup()

	env := map[string]string{
		"LOG_LEVEL":                       "DEBUG",
		"SERVER_PORT":                     "8080",
		"DATABASE_PROVIDER":               "postgresql",
		"DATABASE_CONNECTION_URI":         "postgresql://postgres:postgres@postgres:5432/postgres?schema=public",
		"DATABASE_CONNECTION_CLIENT_NAME": "evolution_api",
		"CACHE_REDIS_ENABLED":             "true",
		"CACHE_REDIS_URI":                 "redis://redis:6379",
		"AUTHENTICATION_API_KEY":          "test-api-key",
	}

	_, baseURL, evolutionCleanup, err := Initialize_Evolution_API(
		t,
		env,
		network.WithNetwork([]string{"evolution-api"}, sharedNetwork),
	)
	if err != nil {
		t.Fatalf("Failed to initialize Evolution API: %v", err)
	}
	defer evolutionCleanup()

	if baseURL == "" {
		t.Fatalf("Expected a valid base URL, got empty")
	}
}

const integrationTestAPIKey = "test-api-key"

func setupIntegrationClient(t *testing.T) (EvolutionClient, func()) {
	t.Helper()

	sharedNetwork, err := network.New(t.Context())
	if err != nil {
		t.Fatalf("failed to initialize docker network: %v", err)
	}
	networkCleanup := func() {
		cleanupCtx, cancel := cleanupContext()
		defer cancel()

		if err := sharedNetwork.Remove(cleanupCtx); err != nil {
			t.Fatalf("failed to remove docker network: %v", err)
		}
	}

	_, _, postgresCleanup, err := InitializePostgres(
		t,
		network.WithNetwork([]string{"postgres"}, sharedNetwork),
	)
	if err != nil {
		t.Fatalf("failed to initialize Postgres: %v", err)
	}

	_, _, redisCleanup, err := InitializeRedis(t, network.WithNetwork([]string{"redis"}, sharedNetwork))
	if err != nil {
		t.Fatalf("failed to initialize Redis: %v", err)
	}

	env := map[string]string{
		"LOG_LEVEL":                       "DEBUG",
		"SERVER_PORT":                     "8080",
		"DATABASE_PROVIDER":               "postgresql",
		"DATABASE_CONNECTION_URI":         "postgresql://postgres:postgres@postgres:5432/postgres?schema=public",
		"DATABASE_CONNECTION_CLIENT_NAME": "evolution_api",
		"CACHE_REDIS_ENABLED":             "true",
		"CACHE_REDIS_URI":                 "redis://redis:6379",
		"AUTHENTICATION_API_KEY":          integrationTestAPIKey,
	}

	_, baseURL, evolutionCleanup, err := Initialize_Evolution_API(
		t,
		env,
		network.WithNetwork([]string{"evolution-api"}, sharedNetwork),
	)
	if err != nil {
		t.Fatalf("failed to initialize Evolution API: %v", err)
	}

	client, err := NewEvolutionClient(EvolutionConfig{
		BaseURL: baseURL,
		APIKey:  integrationTestAPIKey,
	})
	if err != nil {
		t.Fatalf("failed to create Evolution client: %v", err)
	}

	cleanup := func() {
		evolutionCleanup()
		redisCleanup()
		postgresCleanup()
		networkCleanup()
	}

	return client, cleanup
}

func contextWithTimeout(t *testing.T, timeout time.Duration) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(t.Context(), timeout)
}

func waitForInstanceState(ctx context.Context, t *testing.T, client EvolutionClient, instanceName string, expected InstanceState) bool {
	t.Helper()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		stateResp, err := client.InstanceService().ConnectionState(ctx, instanceName)
		if err != nil {
			t.Fatalf("failed to fetch connection state: %v", err)
		}

		if stateResp.Instance.InstanceName == instanceName && stateResp.Instance.State == expected {
			return true
		}

		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			t.Logf("Waiting for instance %s to be %s. Current state: %s", instanceName, expected, stateResp.Instance.State)
		}
	}
}

func createBaileysQRCodeInstance(t *testing.T, client EvolutionClient) (string, func()) {
	t.Helper()

	instanceName := fmt.Sprintf("integration-baileys-%d", time.Now().UnixNano())

	_, err := client.InstanceService().Create(t.Context(), InstanceCreateRequest{
		InstanceName: instanceName,
		Qrcode:       true,
		Integration:  "WHATSAPP-BAILEYS",
	})
	if err != nil {
		t.Fatalf("failed to create instance: %v", err)
	}

	cleanup := func() {
		cleanupCtx, cancel := cleanupContext()
		defer cancel()

		_, _ = client.InstanceService().Delete(cleanupCtx, instanceName)
	}

	return instanceName, cleanup
}

func connectInstanceAndPrintQRCode(t *testing.T, client EvolutionClient, instanceName string) ConnectInstanceResponse {
	t.Helper()

	connectResp, err := client.InstanceService().Connect(t.Context(), instanceName)
	if err != nil {
		t.Fatalf("failed to start QRCode sync: %v", err)
	}

	hasQRCodeData := strings.TrimSpace(connectResp.PairingCode) != "" ||
		strings.TrimSpace(connectResp.Code) != "" ||
		strings.TrimSpace(connectResp.Base64) != ""
	if !hasQRCodeData {
		t.Fatalf("expected pairing or QRCode data, got empty connect response: %+v", connectResp)
	}

	printQRCodeInTerminal(t, connectResp)

	return connectResp
}

func printQRCodeInTerminal(t *testing.T, qr ConnectInstanceResponse) {
	t.Helper()

	if strings.TrimSpace(qr.Base64) == "" {
		t.Logf("QRCode base64 not provided. PairingCode=%s Code=%s", qr.PairingCode, qr.Code)
		return
	}

	dataURI := strings.TrimSpace(qr.Base64)
	_, base64Data, found := strings.Cut(dataURI, ",")
	if !found {
		base64Data = dataURI
	}

	imgBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		imgBytes, err = base64.RawStdEncoding.DecodeString(base64Data)
	}
	if err != nil {
		t.Logf("failed to decode QRCode base64: %v", err)
		t.Logf("Scan fallback base64: %s", qr.Base64)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		t.Logf("failed to decode QRCode image: %v", err)
		t.Logf("Scan fallback base64: %s", qr.Base64)
		return
	}

	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		t.Logf("failed to recognize QR code: %v", err)
		t.Logf("Scan fallback base64: %s", qr.Base64)
		return
	}
	if len(qrCodes) == 0 {
		t.Logf("no QR code payload recognized from data URI")
		t.Logf("Scan fallback base64: %s", qr.Base64)
		return
	}

	payload := strings.TrimSpace(string(qrCodes[0].Payload))
	if payload == "" {
		t.Logf("recognized QR code payload is empty")
		t.Logf("Scan fallback base64: %s", qr.Base64)
		return
	}

	fmt.Println("Scan the QRCode below to connect the instance:")
	qrterminal.GenerateHalfBlock(payload, qrterminal.L, os.Stdout)
}
