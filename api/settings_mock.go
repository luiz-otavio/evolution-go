package evolution

import "context"

var _ SettingsService = (*MockSettingsService)(nil)

type MockSettingsService struct {
	FindFn func(ctx context.Context, instanceName string) (Settings, error)
	SetFn  func(ctx context.Context, instanceName string, req Settings) (SetSettingsResponse, error)
}

func (m *MockSettingsService) Find(ctx context.Context, instanceName string) (Settings, error) {
	return m.FindFn(ctx, instanceName)
}

func (m *MockSettingsService) Set(ctx context.Context, instanceName string, req Settings) (SetSettingsResponse, error) {
	return m.SetFn(ctx, instanceName, req)
}
