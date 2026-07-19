package evolution

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_SettingsService_Find_Validate(t *testing.T) {
	service := NewSettingsService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.Find(t.Context(), "")
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "instanceName is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_SettingsService_Set_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(SettingsSetPath, "inst-a")
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req Settings
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode payload: %v", err)
		}
		if !req.AlwaysOnline {
			t.Fatalf("expected alwaysOnline=true payload, got: %+v", req)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"settings":{"instanceName":"inst-a","settings":{"alwaysOnline":true}}}`))
	}))
	defer server.Close()

	service := NewSettingsService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.Set(t.Context(), "inst-a", Settings{AlwaysOnline: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Settings.Settings.AlwaysOnline {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_SettingsService_Find_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(SettingsFindPath, "inst-a")
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"alwaysOnline":true,"readMessages":true}`))
	}))
	defer server.Close()

	service := NewSettingsService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.Find(t.Context(), "inst-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.AlwaysOnline || !resp.ReadMessages {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_Settings_Integration_SetAndFind(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	client, cleanup := setupIntegrationClient(t)
	defer cleanup()

	instanceName, instanceCleanup := createBaileysQRCodeInstance(t, client)
	defer instanceCleanup()

	setReq := Settings{
		RejectCall:      true,
		MsgCall:         "No calls, please",
		GroupsIgnore:    false,
		AlwaysOnline:    true,
		ReadMessages:    true,
		ReadStatus:      true,
		SyncFullHistory: false,
	}

	setResp, err := client.SettingsService().Set(t.Context(), instanceName, setReq)
	if err != nil {
		t.Fatalf("failed to set settings: %v", err)
	}
	if setResp.Settings.InstanceName == "" {
		t.Fatalf("expected instanceName in set response, got: %+v", setResp)
	}

	found, err := client.SettingsService().Find(t.Context(), instanceName)
	if err != nil {
		t.Fatalf("failed to find settings: %v", err)
	}

	if found.MsgCall != setReq.MsgCall || found.AlwaysOnline != setReq.AlwaysOnline || found.RejectCall != setReq.RejectCall {
		t.Fatalf("settings mismatch, expected %+v, got %+v", setReq, found)
	}
}
