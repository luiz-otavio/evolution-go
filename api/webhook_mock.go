package evolution

import "context"

var _ WebhookService = (*MockWebhookService)(nil)

type MockWebhookService struct {
	FindWebhookFn  func(ctx context.Context, instanceName string) (*Webhook, error)
	SetWebhookFn   func(ctx context.Context, instanceName string, req SetWebhookRequest) (SuccessResponse, error)
	SetWebSocketFn func(ctx context.Context, instanceName string, req SetWebSocketRequest) (SuccessResponse, error)
}

func (m *MockWebhookService) FindWebhook(ctx context.Context, instanceName string) (*Webhook, error) {
	return m.FindWebhookFn(ctx, instanceName)
}

func (m *MockWebhookService) SetWebhook(ctx context.Context, instanceName string, req SetWebhookRequest) (SuccessResponse, error) {
	return m.SetWebhookFn(ctx, instanceName, req)
}

func (m *MockWebhookService) SetWebSocket(ctx context.Context, instanceName string, req SetWebSocketRequest) (SuccessResponse, error) {
	return m.SetWebSocketFn(ctx, instanceName, req)
}
