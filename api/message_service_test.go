package evolution

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_MessageService_SendText_Validate(t *testing.T) {
	service := NewMessageService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.SendText(t.Context(), "inst", SendTextRequest{Number: "5511999999999"})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "text is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_MessageService_SendMedia_Validate(t *testing.T) {
	service := NewMessageService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.SendMedia(t.Context(), "inst", SendMediaRequest{
		Number:    "5511999999999",
		MediaType: "invalid",
		Media:     "https://example.com/file.png",
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "invalid mediatype") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_MessageService_SendText_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(MessageSendTextPath, "inst-a")
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if r.Header.Get("apikey") != "test-key" {
			t.Fatalf("expected apikey header")
		}

		var req SendTextRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Number != "5511999999999" || req.Text != "oi" {
			t.Fatalf("unexpected request payload: %+v", req)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"PENDING","key":{"id":"msg-1"}}`))
	}))
	defer server.Close()

	service := NewMessageService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.SendText(t.Context(), "inst-a", SendTextRequest{
		Number: "5511999999999",
		Text:   "oi",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "PENDING" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_MessageService_SendMedia_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(MessageSendMediaPath, "inst-a")
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req SendMediaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.MediaType != MediaTypeImage {
			t.Fatalf("unexpected media type: %+v", req)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"PENDING"}`))
	}))
	defer server.Close()

	service := NewMessageService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.SendMedia(t.Context(), "inst-a", SendMediaRequest{
		Number:    "5511999999999",
		MediaType: MediaTypeImage,
		Media:     "https://example.com/image.png",
		Caption:   "hello",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "PENDING" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_MessageService_SendButtons_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(MessageSendButtonsPath, "inst-a")
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req SendButtonsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Number != "5511999999999" {
			t.Fatalf("unexpected number: %+v", req)
		}
		if len(req.Buttons) != 1 || req.Buttons[0].Type != "reply" || req.Buttons[0].ID != "btn-1" {
			t.Fatalf("unexpected buttons payload: %+v", req)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"PENDING"}`))
	}))
	defer server.Close()

	service := NewMessageService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.SendButtons(t.Context(), "inst-a", SendButtonsRequest{
		Number: "5511999999999",
		Title:  "Titulo",
		Buttons: []Button{
			{Type: "reply", ID: "btn-1", DisplayText: "Opcao 1"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "PENDING" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_MessageService_SendTemplate_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(MessageSendTemplatePath, "inst-a")
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req SendTemplateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Name != "template_name" || req.Language != "pt_BR" {
			t.Fatalf("unexpected template payload: %+v", req)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"PENDING"}`))
	}))
	defer server.Close()

	service := NewMessageService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.SendTemplate(t.Context(), "inst-a", SendTemplateRequest{
		Number:   "5511999999999",
		Name:     "template_name",
		Language: "pt_BR",
		Components: []map[string]any{
			{"type": "body", "parameters": []map[string]any{{"type": "text", "text": "teste"}}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "PENDING" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_MessageService_SendContact_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(MessageSendContactPath, "inst-a")
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req SendContactRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if len(req.Contact) != 1 || req.Contact[0].FullName == "" {
			t.Fatalf("unexpected contact payload: %+v", req)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"PENDING"}`))
	}))
	defer server.Close()

	service := NewMessageService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.SendContact(t.Context(), "inst-a", SendContactRequest{
		Number: "5511999999999",
		Contact: []ContactMessage{
			{FullName: "Contato", PhoneNumber: "5511999999999"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "PENDING" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_MessageService_SendReaction_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(MessageSendReactionPath, "inst-a")
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req SendReactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Key["id"] == nil || req.Reaction != ":like:" {
			t.Fatalf("unexpected reaction payload: %+v", req)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"PENDING"}`))
	}))
	defer server.Close()

	service := NewMessageService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.SendReaction(t.Context(), "inst-a", SendReactionRequest{
		Key: map[string]any{
			"id":        "ABCD1234",
			"remoteJid": "5511999999999@s.whatsapp.net",
			"fromMe":    false,
		},
		Reaction: ":like:",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "PENDING" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_Message_Integration_SendText_EnvGated(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	if err := ReadLocalDotenv(); err != nil {
		t.Fatalf("failed to read .env file: %v", err)
	}

	targetNumber := strings.TrimSpace(os.Getenv("EVOLUTION_TEST_TARGET_NUMBER"))
	if targetNumber == "" {
		t.Skip("set EVOLUTION_TEST_TARGET_NUMBER to run message integration")
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

	resp, err := client.MessageService().SendText(t.Context(), instanceName, SendTextRequest{
		Number: targetNumber,
		Text:   "Lorenzo cuzudokkkkkkkkk bot testando",
	})
	if err != nil {
		t.Fatalf("failed to send text message: %v", err)
	}

	if resp.Status == "" {
		t.Fatalf("expected message status, got: %+v", resp)
	}
}

func Test_Message_Integration_SendButtons_EnvGated(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	if err := ReadLocalDotenv(); err != nil {
		t.Fatalf("failed to read .env file: %v", err)
	}

	targetNumber := strings.TrimSpace(os.Getenv("EVOLUTION_TEST_TARGET_NUMBER"))
	if targetNumber == "" {
		t.Skip("set EVOLUTION_TEST_TARGET_NUMBER to run message integration")
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

	resp, err := client.MessageService().SendButtons(t.Context(), instanceName, SendButtonsRequest{
		Number: targetNumber,
		Buttons: []Button{
			{Type: "reply", ID: "1", DisplayText: "Opcao 1"},
			{Type: "reply", ID: "2", DisplayText: "Opcao 2"},
		},
		Description: "aonde que fica a descrição",
		Footer:      "bumbum guloso",
		Delay:       5,
		Title:       "Lorenzo cuzudokkkkkkk",
	})

	if err != nil {
		t.Fatalf("failed to send text message: %v", err)
	}

	if resp.Status == "" {
		t.Fatalf("expected message status, got: %+v", resp)
	}
}
