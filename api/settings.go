package evolution

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

const (
	SettingsFindPath = "/settings/find"
	SettingsSetPath  = "/settings/set"
)

// Settings represents the per-instance behaviour configuration.
type Settings struct {
	RejectCall      bool   `json:"rejectCall"`
	MsgCall         string `json:"msgCall"`
	GroupsIgnore    bool   `json:"groupsIgnore"`
	AlwaysOnline    bool   `json:"alwaysOnline"`
	ReadMessages    bool   `json:"readMessages"`
	ReadStatus      bool   `json:"readStatus"`
	SyncFullHistory bool   `json:"syncFullHistory"`
	WavoipToken     string `json:"wavoipToken,omitempty"`
}

type SetSettingsResponse struct {
	Settings SetSettingsData `json:"settings"`
}

type SetSettingsData struct {
	InstanceName string   `json:"instanceName"`
	Settings     Settings `json:"settings"`
}

type SettingsService interface {
	Find(ctx context.Context, instanceName string) (Settings, error)
	Set(ctx context.Context, instanceName string, req Settings) (SetSettingsResponse, error)
}

type settingsService struct {
	http   HttpProvider
	apiKey string
}

func NewSettingsService(http HttpProvider, apiKey string) SettingsService {
	return settingsService{http: http, apiKey: apiKey}
}

func (s settingsService) Find(ctx context.Context, instanceName string) (Settings, error) {
	if instanceName == "" {
		return Settings{}, fmt.Errorf("instanceName is required")
	}

	path := buildInstancePath(SettingsFindPath, instanceName)
	return executeGet[Settings](ctx, s.http, s.apiKey, path, nil)
}

func (s settingsService) Set(ctx context.Context, instanceName string, req Settings) (SetSettingsResponse, error) {
	if instanceName == "" {
		return SetSettingsResponse{}, fmt.Errorf("instanceName is required")
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return SetSettingsResponse{}, fmt.Errorf("failed to marshal settings request: %w", err)
	}

	path := buildInstancePath(SettingsSetPath, instanceName)
	return executePost[SetSettingsResponse](ctx, s.http, s.apiKey, path, payload)
}
