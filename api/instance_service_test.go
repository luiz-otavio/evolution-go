package evolution

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_InstanceService_Create_Validate(t *testing.T) {
	service := NewInstanceService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.Create(t.Context(), InstanceCreateRequest{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "instanceName is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_InstanceService_Connect_Validate(t *testing.T) {
	service := NewInstanceService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.Connect(t.Context(), "")
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "instanceName is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_InstanceService_SetPresence_Validate(t *testing.T) {
	service := NewInstanceService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.SetPresence(t.Context(), "inst", SetPresenceRequest{Presence: "invalid"})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "invalid presence") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_InstanceService_Create_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != InstanceCreatePath {
			t.Fatalf("expected path %s, got %s", InstanceCreatePath, r.URL.Path)
		}
		if r.Header.Get("apikey") != "test-key" {
			t.Fatalf("expected apikey header")
		}

		var req InstanceCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.InstanceName != "inst-a" {
			t.Fatalf("unexpected instance name: %s", req.InstanceName)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"instance":{"instanceName":"inst-a"},"qrcode":{"code":"ABC","base64":"data:image/png;base64,xyz"}}`))
	}))
	defer server.Close()

	service := NewInstanceService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.Create(t.Context(), InstanceCreateRequest{
		InstanceName: "inst-a",
		Qrcode:       true,
		Integration:  "WHATSAPP-BAILEYS",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Qrcode == nil || resp.Qrcode.Code != "ABC" {
		t.Fatalf("expected qrcode payload, got: %+v", resp)
	}
}

func Test_InstanceService_Connect_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(InstanceConnectPath, "my instance")
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.EscapedPath() != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.EscapedPath())
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"pairingCode":"123-456","code":"qrcode-code","base64":"data:image/png;base64,abc","count":1}`))
	}))
	defer server.Close()

	service := NewInstanceService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.Connect(t.Context(), "my instance")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PairingCode == "" || resp.Base64 == "" {
		t.Fatalf("expected connect payload, got: %+v", resp)
	}
}
func Test_Instance_Integration_CreateAndStartQRCodeSync(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	client, cleanup := setupIntegrationClient(t)
	defer cleanup()

	instanceName, instanceCleanup := createBaileysQRCodeInstance(t, client)
	defer instanceCleanup()

	_ = connectInstanceAndPrintQRCode(t, client, instanceName)

	waitCtx, cancel := contextWithTimeout(t, 2*time.Minute)
	defer cancel()

	if !waitForInstanceState(waitCtx, t, client, instanceName, InstanceStateOpen) {
		t.Skip("instance did not reach open state in time; pair the QRCode and rerun")
	}
}
