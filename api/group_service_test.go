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

func findGroupJIDByName(groups []GroupSummary, targetName string) string {
	target := strings.ToLower(strings.TrimSpace(targetName))
	if target == "" {
		return ""
	}

	for _, group := range groups {
		jid := strings.TrimSpace(group.ID)
		if !strings.HasSuffix(jid, "@g.us") {
			jid = strings.TrimSpace(group.JID)
		}
		if !strings.HasSuffix(jid, "@g.us") {
			jid = strings.TrimSpace(group.GroupJID)
		}
		if jid == "" {
			continue
		}

		for _, name := range []string{group.Subject, group.Name, group.PushName} {
			if strings.ToLower(strings.TrimSpace(name)) == target {
				return jid
			}
		}
	}

	return ""
}

func Test_GroupService_Create_Validate(t *testing.T) {
	service := NewGroupService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.Create(t.Context(), "inst", CreateGroupRequest{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "subject is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_GroupService_UpdateParticipant_Validate(t *testing.T) {
	service := NewGroupService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.UpdateParticipant(t.Context(), "inst", UpdateParticipantRequest{
		GroupJid:     "123@g.us",
		Action:       "invalid",
		Participants: []string{"5511999999999@s.whatsapp.net"},
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "invalid action") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_GroupService_FindGroupInfos_Validate(t *testing.T) {
	service := NewGroupService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.FindGroupInfos(t.Context(), "", "123@g.us")
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "instanceName is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_GroupService_FetchAllGroups_Validate(t *testing.T) {
	service := NewGroupService(NewHttpProvider("http://127.0.0.1:1"), "test-key")

	_, err := service.FetchAllGroups(t.Context(), "", false)
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "instanceName is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_GroupService_Create_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(GroupCreatePath, "inst-a")
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req CreateGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode payload: %v", err)
		}
		if req.Subject != "Test Group" || len(req.Participants) != 1 {
			t.Fatalf("unexpected payload: %+v", req)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"success":true,"group":{"id":"123@g.us"}}`))
	}))
	defer server.Close()

	service := NewGroupService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.Create(t.Context(), "inst-a", CreateGroupRequest{
		Subject:      "Test Group",
		Participants: []string{"5511999999999@s.whatsapp.net"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_GroupService_FindGroupInfos_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(GroupFindGroupInfosPath, "inst-a")
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		groupJid := r.URL.Query().Get("groupJid")
		if groupJid != "123@g.us" {
			t.Fatalf("expected query groupJid=123@g.us, got %s", groupJid)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"group":{"id":"123@g.us","subject":"Test"}}`))
	}))
	defer server.Close()

	service := NewGroupService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.FindGroupInfos(t.Context(), "inst-a", "123@g.us")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_GroupService_FetchAllGroups_RequestContract(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := buildInstancePath(GroupFetchAllGroupsPath, "inst-a")
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		getParticipants := r.URL.Query().Get("getParticipants")
		if getParticipants != "false" {
			t.Fatalf("expected query getParticipants=false, got %s", getParticipants)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":"123@g.us","subject":"Private"}]`))
	}))
	defer server.Close()

	service := NewGroupService(NewHttpProvider(server.URL), "test-key")

	resp, err := service.FetchAllGroups(t.Context(), "inst-a", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func Test_Group_Integration_FindGroupInfosAndParticipants_EnvGated(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	groupJid := strings.TrimSpace(os.Getenv("EVOLUTION_TEST_GROUP_JID"))
	if groupJid == "" {
		t.Skip("set EVOLUTION_TEST_GROUP_JID to run group integration")
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

	info, err := client.GroupService().FindGroupInfos(t.Context(), instanceName, groupJid)
	if err != nil {
		t.Fatalf("failed to find group infos: %v", err)
	}
	if !info.Success {
		t.Fatalf("expected success=true from FindGroupInfos, got: %+v", info)
	}

	participants, err := client.GroupService().Participants(t.Context(), instanceName, groupJid)
	if err != nil {
		t.Fatalf("failed to list participants: %v", err)
	}
	if !participants.Success {
		t.Fatalf("expected success=true from Participants, got: %+v", participants)
	}
}

func Test_Group_Integration_FindPrivateGroupAndSendMessage_EnvGated(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	if err := ReadLocalDotenv(); err != nil {
		t.Fatalf("failed to read .env file: %v", err)
	}

	groupName := strings.TrimSpace(os.Getenv("EVOLUTION_TEST_GROUP_NAME"))
	if groupName == "" {
		groupName = "Private"
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

	groupsResp, err := client.GroupService().FetchAllGroups(t.Context(), instanceName, false)
	if err != nil {
		t.Fatalf("failed to list groups: %v", err)
	}
	if len(groupsResp) == 0 {
		t.Skip("no groups found; ensure this account has groups and rerun")
	}

	groupJID := findGroupJIDByName(groupsResp, groupName)
	if groupJID == "" {
		t.Skipf("group %q not found among fetched groups; create/join the group and rerun", groupName)
	}

	resp, err := client.MessageService().SendText(t.Context(), instanceName, SendTextRequest{
		Number: groupJID,
		Text:   "integration test message to group via JID",
	})
	if err != nil {
		t.Fatalf("failed to send group message using JID %s: %v", groupJID, err)
	}
	if resp.Status == "" {
		t.Fatalf("expected message status, got: %+v", resp)
	}
}
