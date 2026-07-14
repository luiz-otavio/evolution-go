package evolution

import "context"

var _ ProxyService = (*MockProxyService)(nil)

type MockProxyService struct {
	FindFn func(ctx context.Context, instanceName string) (*Proxy, error)
	SetFn  func(ctx context.Context, instanceName string, req Proxy) (SuccessResponse, error)
}

func (m *MockProxyService) Find(ctx context.Context, instanceName string) (*Proxy, error) {
	return m.FindFn(ctx, instanceName)
}

func (m *MockProxyService) Set(ctx context.Context, instanceName string, req Proxy) (SuccessResponse, error) {
	return m.SetFn(ctx, instanceName, req)
}
