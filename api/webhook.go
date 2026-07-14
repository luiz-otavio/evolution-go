package evolution

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

const (
	WebhookFindPath  = "/webhook/find"
	WebhookSetPath   = "/webhook/set"
	WebSocketSetPath = "/websocket/set"
)

type Webhook struct {
	Enabled bool     `json:"enabled"`
	URL     string   `json:"url,omitempty"`
	Events  []string `json:"events,omitempty"`
}

type SetWebhookRequest struct {
	Enabled bool              `json:"enabled"`
	URL     string            `json:"url,omitempty"`
	Events  []string          `json:"events,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Base64  bool              `json:"base64,omitempty"`
}

type SetWebSocketRequest struct {
	Enabled bool     `json:"enabled"`
	Events  []string `json:"events,omitempty"`
}

type WebhookService interface {
	// FindWebhook returns the webhook configuration, or nil when none is configured.
	FindWebhook(ctx context.Context, instanceName string) (*Webhook, error)
	SetWebhook(ctx context.Context, instanceName string, req SetWebhookRequest) (SuccessResponse, error)
	SetWebSocket(ctx context.Context, instanceName string, req SetWebSocketRequest) (SuccessResponse, error)
}

type webhookService struct {
	http   HttpProvider
	apiKey string
}

func NewWebhookService(http HttpProvider, apiKey string) WebhookService {
	return webhookService{http: http, apiKey: apiKey}
}

func (s webhookService) FindWebhook(ctx context.Context, instanceName string) (*Webhook, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}

	path := buildInstancePath(WebhookFindPath, instanceName)
	return executeGet[*Webhook](ctx, s.http, s.apiKey, path, nil)
}

func (s webhookService) SetWebhook(ctx context.Context, instanceName string, req SetWebhookRequest) (SuccessResponse, error) {
	if instanceName == "" {
		return SuccessResponse{}, fmt.Errorf("instanceName is required")
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return SuccessResponse{}, fmt.Errorf("failed to marshal webhook request: %w", err)
	}

	path := buildInstancePath(WebhookSetPath, instanceName)
	return executePost[SuccessResponse](ctx, s.http, s.apiKey, path, payload)
}

func (s webhookService) SetWebSocket(ctx context.Context, instanceName string, req SetWebSocketRequest) (SuccessResponse, error) {
	if instanceName == "" {
		return SuccessResponse{}, fmt.Errorf("instanceName is required")
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return SuccessResponse{}, fmt.Errorf("failed to marshal websocket request: %w", err)
	}

	path := buildInstancePath(WebSocketSetPath, instanceName)
	return executePost[SuccessResponse](ctx, s.http, s.apiKey, path, payload)
}
