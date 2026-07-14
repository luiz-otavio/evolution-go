package evolution

import "context"

var _ InstanceService = (*MockInstanceService)(nil)

type MockInstanceService struct {
	CreateFn          func(ctx context.Context, req InstanceCreateRequest) (InstanceCreateResponse, error)
	ConnectFn         func(ctx context.Context, instanceName string) (ConnectInstanceResponse, error)
	ConnectionStateFn func(ctx context.Context, instanceName string) (ConnectionStateResponse, error)
	FetchInstancesFn  func(ctx context.Context) ([]InstanceResponse, error)
	RestartFn         func(ctx context.Context, instanceName string) (RestartInstanceResponse, error)
	LogoutFn          func(ctx context.Context, instanceName string) (SuccessResponse, error)
	DeleteFn          func(ctx context.Context, instanceName string) (SuccessResponse, error)
	SetPresenceFn     func(ctx context.Context, instanceName string, req SetPresenceRequest) (SuccessResponse, error)
}

func (m *MockInstanceService) Create(ctx context.Context, req InstanceCreateRequest) (InstanceCreateResponse, error) {
	return m.CreateFn(ctx, req)
}

func (m *MockInstanceService) Connect(ctx context.Context, instanceName string) (ConnectInstanceResponse, error) {
	return m.ConnectFn(ctx, instanceName)
}

func (m *MockInstanceService) ConnectionState(ctx context.Context, instanceName string) (ConnectionStateResponse, error) {
	return m.ConnectionStateFn(ctx, instanceName)
}

func (m *MockInstanceService) FetchInstances(ctx context.Context) ([]InstanceResponse, error) {
	return m.FetchInstancesFn(ctx)
}

func (m *MockInstanceService) Restart(ctx context.Context, instanceName string) (RestartInstanceResponse, error) {
	return m.RestartFn(ctx, instanceName)
}

func (m *MockInstanceService) Logout(ctx context.Context, instanceName string) (SuccessResponse, error) {
	return m.LogoutFn(ctx, instanceName)
}

func (m *MockInstanceService) Delete(ctx context.Context, instanceName string) (SuccessResponse, error) {
	return m.DeleteFn(ctx, instanceName)
}

func (m *MockInstanceService) SetPresence(ctx context.Context, instanceName string, req SetPresenceRequest) (SuccessResponse, error) {
	return m.SetPresenceFn(ctx, instanceName, req)
}
