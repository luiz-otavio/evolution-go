package evolution

import "context"

var _ ChatService = (*MockChatService)(nil)

type MockChatService struct {
	CheckWhatsAppNumbersFn func(ctx context.Context, instanceName string, req WhatsAppNumbersRequest) (WhatsAppNumbersResponse, error)
	MarkMessageAsReadFn    func(ctx context.Context, instanceName string, req MarkMessageAsReadRequest) (SuccessResponse, error)
	ArchiveChatFn          func(ctx context.Context, instanceName string, req ArchiveChatRequest) (SuccessResponse, error)
	FindChatsFn            func(ctx context.Context, instanceName string, query Query) ([]ChatSummary, error)
	FindContactsFn         func(ctx context.Context, instanceName string, query Query) ([]Contact, error)
	FindMessagesFn         func(ctx context.Context, instanceName string, query Query) (FindMessagesResponse, error)
	UpdateProfileNameFn    func(ctx context.Context, instanceName string, req UpdateProfileNameRequest) (SuccessResponse, error)
	UpdateProfilePictureFn func(ctx context.Context, instanceName string, req UpdateProfilePictureRequest) (SuccessResponse, error)
	UpdateProfileStatusFn  func(ctx context.Context, instanceName string, req UpdateProfileStatusRequest) (SuccessResponse, error)
}

func (m *MockChatService) CheckWhatsAppNumbers(ctx context.Context, instanceName string, req WhatsAppNumbersRequest) (WhatsAppNumbersResponse, error) {
	return m.CheckWhatsAppNumbersFn(ctx, instanceName, req)
}

func (m *MockChatService) MarkMessageAsRead(ctx context.Context, instanceName string, req MarkMessageAsReadRequest) (SuccessResponse, error) {
	return m.MarkMessageAsReadFn(ctx, instanceName, req)
}

func (m *MockChatService) ArchiveChat(ctx context.Context, instanceName string, req ArchiveChatRequest) (SuccessResponse, error) {
	return m.ArchiveChatFn(ctx, instanceName, req)
}

func (m *MockChatService) FindChats(ctx context.Context, instanceName string, query Query) ([]ChatSummary, error) {
	return m.FindChatsFn(ctx, instanceName, query)
}

func (m *MockChatService) FindContacts(ctx context.Context, instanceName string, query Query) ([]Contact, error) {
	return m.FindContactsFn(ctx, instanceName, query)
}

func (m *MockChatService) FindMessages(ctx context.Context, instanceName string, query Query) (FindMessagesResponse, error) {
	return m.FindMessagesFn(ctx, instanceName, query)
}

func (m *MockChatService) UpdateProfileName(ctx context.Context, instanceName string, req UpdateProfileNameRequest) (SuccessResponse, error) {
	return m.UpdateProfileNameFn(ctx, instanceName, req)
}

func (m *MockChatService) UpdateProfilePicture(ctx context.Context, instanceName string, req UpdateProfilePictureRequest) (SuccessResponse, error) {
	return m.UpdateProfilePictureFn(ctx, instanceName, req)
}

func (m *MockChatService) UpdateProfileStatus(ctx context.Context, instanceName string, req UpdateProfileStatusRequest) (SuccessResponse, error) {
	return m.UpdateProfileStatusFn(ctx, instanceName, req)
}
