package evolution

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

const (
	InstanceCreatePath          = "/instance/create"
	InstanceConnectPath         = "/instance/connect"
	InstanceConnectionStatePath = "/instance/connectionState"
	InstanceFetchPath           = "/instance/fetchInstances"
	InstanceRestartPath         = "/instance/restart"
	InstanceLogoutPath          = "/instance/logout"
	InstanceDeletePath          = "/instance/delete"
	InstanceSetPresencePath     = "/instance/setPresence"
)

// InstanceWebhook is the optional webhook configuration accepted when creating an instance.
type InstanceWebhook struct {
	Enabled bool     `json:"enabled,omitempty"`
	URL     string   `json:"url,omitempty"`
	Events  []string `json:"events,omitempty"`
}

type InstanceCreateRequest struct {
	InstanceName string           `json:"instanceName"`
	Qrcode       bool             `json:"qrcode,omitempty"`
	Integration  string           `json:"integration,omitempty"`
	Token        string           `json:"token,omitempty"`
	Number       string           `json:"number,omitempty"`
	Webhook      *InstanceWebhook `json:"webhook,omitempty"`
}

func (r InstanceCreateRequest) Validate() error {
	if r.InstanceName == "" {
		return fmt.Errorf("instanceName is required")
	}
	return nil
}

type InstanceCreateResponse struct {
	Instance map[string]any `json:"instance"`
	Hash     any            `json:"hash,omitempty"`
	Webhook  map[string]any `json:"webhook,omitempty"`
	Settings map[string]any `json:"settings,omitempty"`
	Qrcode   *QRCode        `json:"qrcode,omitempty"`
}

type QRCode struct {
	PairingCode string `json:"pairingCode,omitempty"`
	Code        string `json:"code,omitempty"`
	Base64      string `json:"base64,omitempty"`
	Count       int    `json:"count,omitempty"`
}

// ConnectInstanceResponse is returned by the Connect endpoint (QR code / pairing data).
type ConnectInstanceResponse struct {
	PairingCode string `json:"pairingCode"`
	Code        string `json:"code"`
	Base64      string `json:"base64"`
	Count       int    `json:"count"`
}

type ConnectionStateResponse struct {
	Instance ConnectionStateInstance `json:"instance"`
}

type ConnectionStateInstance struct {
	InstanceName string `json:"instanceName"`
	State        string `json:"state"`
}

type InstanceResponse struct {
	Instance map[string]any `json:"instance"`
}

type RestartInstanceResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
}

type SetPresenceRequest struct {
	Presence Presence `json:"presence"`
}

func (r SetPresenceRequest) Validate() error {
	switch r.Presence {
	case PresenceAvailable, PresenceUnavailable, PresenceComposing, PresenceRecording:
		return nil
	default:
		return fmt.Errorf("invalid presence: %q", r.Presence)
	}
}

type InstanceService interface {
	Create(ctx context.Context, req InstanceCreateRequest) (InstanceCreateResponse, error)
	Connect(ctx context.Context, instanceName string) (ConnectInstanceResponse, error)
	ConnectionState(ctx context.Context, instanceName string) (ConnectionStateResponse, error)
	FetchInstances(ctx context.Context) ([]InstanceResponse, error)
	Restart(ctx context.Context, instanceName string) (RestartInstanceResponse, error)
	Logout(ctx context.Context, instanceName string) (SuccessResponse, error)
	Delete(ctx context.Context, instanceName string) (SuccessResponse, error)
	SetPresence(ctx context.Context, instanceName string, req SetPresenceRequest) (SuccessResponse, error)
}

type instanceService struct {
	http   HttpProvider
	apiKey string
}

func NewInstanceService(http HttpProvider, apiKey string) InstanceService {
	return instanceService{http: http, apiKey: apiKey}
}

func (s instanceService) Create(ctx context.Context, req InstanceCreateRequest) (InstanceCreateResponse, error) {
	if err := req.Validate(); err != nil {
		return InstanceCreateResponse{}, fmt.Errorf("invalid create instance request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return InstanceCreateResponse{}, fmt.Errorf("failed to marshal create instance request: %w", err)
	}

	return executePost[InstanceCreateResponse](ctx, s.http, s.apiKey, InstanceCreatePath, payload)
}

func (s instanceService) Connect(ctx context.Context, instanceName string) (ConnectInstanceResponse, error) {
	if instanceName == "" {
		return ConnectInstanceResponse{}, fmt.Errorf("instanceName is required")
	}

	path := buildInstancePath(InstanceConnectPath, instanceName)
	return executeGet[ConnectInstanceResponse](ctx, s.http, s.apiKey, path, nil)
}

func (s instanceService) ConnectionState(ctx context.Context, instanceName string) (ConnectionStateResponse, error) {
	if instanceName == "" {
		return ConnectionStateResponse{}, fmt.Errorf("instanceName is required")
	}

	path := buildInstancePath(InstanceConnectionStatePath, instanceName)
	return executeGet[ConnectionStateResponse](ctx, s.http, s.apiKey, path, nil)
}

func (s instanceService) FetchInstances(ctx context.Context) ([]InstanceResponse, error) {
	return executeGet[[]InstanceResponse](ctx, s.http, s.apiKey, InstanceFetchPath, nil)
}

func (s instanceService) Restart(ctx context.Context, instanceName string) (RestartInstanceResponse, error) {
	if instanceName == "" {
		return RestartInstanceResponse{}, fmt.Errorf("instanceName is required")
	}

	path := buildInstancePath(InstanceRestartPath, instanceName)
	return executePost[RestartInstanceResponse](ctx, s.http, s.apiKey, path, nil)
}

func (s instanceService) Logout(ctx context.Context, instanceName string) (SuccessResponse, error) {
	if instanceName == "" {
		return SuccessResponse{}, fmt.Errorf("instanceName is required")
	}

	path := buildInstancePath(InstanceLogoutPath, instanceName)
	return executeDelete[SuccessResponse](ctx, s.http, s.apiKey, path, nil)
}

func (s instanceService) Delete(ctx context.Context, instanceName string) (SuccessResponse, error) {
	if instanceName == "" {
		return SuccessResponse{}, fmt.Errorf("instanceName is required")
	}

	path := buildInstancePath(InstanceDeletePath, instanceName)
	return executeDelete[SuccessResponse](ctx, s.http, s.apiKey, path, nil)
}

func (s instanceService) SetPresence(ctx context.Context, instanceName string, req SetPresenceRequest) (SuccessResponse, error) {
	if instanceName == "" {
		return SuccessResponse{}, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return SuccessResponse{}, fmt.Errorf("invalid set presence request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return SuccessResponse{}, fmt.Errorf("failed to marshal set presence request: %w", err)
	}

	path := buildInstancePath(InstanceSetPresencePath, instanceName)
	return executePost[SuccessResponse](ctx, s.http, s.apiKey, path, payload)
}
