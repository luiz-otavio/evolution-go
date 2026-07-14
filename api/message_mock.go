package evolution

import "context"

var _ MessageService = (*MockMessageService)(nil)

type MockMessageService struct {
	SendTextFn     func(ctx context.Context, instanceName string, req SendTextRequest) (MessageResponse, error)
	SendMediaFn    func(ctx context.Context, instanceName string, req SendMediaRequest) (MessageResponse, error)
	SendButtonsFn  func(ctx context.Context, instanceName string, req SendButtonsRequest) (MessageResponse, error)
	SendListFn     func(ctx context.Context, instanceName string, req SendListRequest) (MessageResponse, error)
	SendContactFn  func(ctx context.Context, instanceName string, req SendContactRequest) (MessageResponse, error)
	SendLocationFn func(ctx context.Context, instanceName string, req SendLocationRequest) (MessageResponse, error)
	SendPollFn     func(ctx context.Context, instanceName string, req SendPollRequest) (MessageResponse, error)
	SendReactionFn func(ctx context.Context, instanceName string, req SendReactionRequest) (MessageResponse, error)
	SendTemplateFn func(ctx context.Context, instanceName string, req SendTemplateRequest) (MessageResponse, error)
}

func (m *MockMessageService) SendText(ctx context.Context, instanceName string, req SendTextRequest) (MessageResponse, error) {
	return m.SendTextFn(ctx, instanceName, req)
}

func (m *MockMessageService) SendMedia(ctx context.Context, instanceName string, req SendMediaRequest) (MessageResponse, error) {
	return m.SendMediaFn(ctx, instanceName, req)
}

func (m *MockMessageService) SendButtons(ctx context.Context, instanceName string, req SendButtonsRequest) (MessageResponse, error) {
	return m.SendButtonsFn(ctx, instanceName, req)
}

func (m *MockMessageService) SendList(ctx context.Context, instanceName string, req SendListRequest) (MessageResponse, error) {
	return m.SendListFn(ctx, instanceName, req)
}

func (m *MockMessageService) SendContact(ctx context.Context, instanceName string, req SendContactRequest) (MessageResponse, error) {
	return m.SendContactFn(ctx, instanceName, req)
}

func (m *MockMessageService) SendLocation(ctx context.Context, instanceName string, req SendLocationRequest) (MessageResponse, error) {
	return m.SendLocationFn(ctx, instanceName, req)
}

func (m *MockMessageService) SendPoll(ctx context.Context, instanceName string, req SendPollRequest) (MessageResponse, error) {
	return m.SendPollFn(ctx, instanceName, req)
}

func (m *MockMessageService) SendReaction(ctx context.Context, instanceName string, req SendReactionRequest) (MessageResponse, error) {
	return m.SendReactionFn(ctx, instanceName, req)
}

func (m *MockMessageService) SendTemplate(ctx context.Context, instanceName string, req SendTemplateRequest) (MessageResponse, error) {
	return m.SendTemplateFn(ctx, instanceName, req)
}
