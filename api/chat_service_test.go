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

func Test_ChatService_CheckWhatsAppNumbers_Validate(t *testing.T) {
	service := NewChatService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.CheckWhatsAppNumbers(t.Context(), "inst", WhatsAppNumbersRequest{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "numbers is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_ChatService_MarkMessageAsRead_Validate(t *testing.T) {
	service := NewChatService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.MarkMessageAsRead(t.Context(), "inst", MarkMessageAsReadRequest{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "readMessages is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_ChatService_FindMessages_ValidateInstance(t *testing.T) {
	service := NewChatService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.FindMessages(t.Context(), "", Query{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "instanceName is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_ChatService_CheckWhatsAppNumbers_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(ChatWhatsAppNumbersPath, "inst-a")
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req WhatsAppNumbersRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode payload: %v", err)
		}
		if len(req.Numbers) != 1 || req.Numbers[0] != "5511999999999" {
			t.Fatalf("unexpected payload: %+v", req)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"numbers":[{"number":"5511999999999","exists":true}]}`))
	}))
	defer server.Close()

	service := NewChatService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.CheckWhatsAppNumbers(t.Context(), "inst-a", WhatsAppNumbersRequest{Numbers: []string{"5511999999999"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Numbers) != 1 || !resp.Numbers[0].Exists {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_ChatService_FindMessages_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(ChatFindMessagesPath, "inst-a")
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"messages":{"total":1,"pages":1,"currentPage":1,"records":[{"id":"m1"}]}}`))
	}))
	defer server.Close()

	service := NewChatService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.FindMessages(t.Context(), "inst-a", Query{Take: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Messages.Total != 1 || len(resp.Messages.Records) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_Chat_Integration_CheckWhatsAppNumbers_EnvGated(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	targetNumber := strings.TrimSpace(os.Getenv("EVOLUTION_TEST_TARGET_NUMBER"))
	if targetNumber == "" {
		t.Skip("set EVOLUTION_TEST_TARGET_NUMBER to run chat integration")
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

	resp, err := client.ChatService().CheckWhatsAppNumbers(t.Context(), instanceName, WhatsAppNumbersRequest{
		Numbers: []string{targetNumber},
	})
	if err != nil {
		t.Fatalf("failed to check whatsapp numbers: %v", err)
	}
	if len(resp.Numbers) == 0 {
		t.Fatalf("expected at least one number status, got: %+v", resp)
	}
}
