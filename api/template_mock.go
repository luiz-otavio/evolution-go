package evolution

import "context"

var _ TemplateService = (*MockTemplateService)(nil)

type MockTemplateService struct {
	CreateFn func(ctx context.Context, instanceName string, req CreateTemplateRequest) (TemplateResponse, error)
	EditFn   func(ctx context.Context, instanceName string, req EditTemplateRequest) (TemplateResponse, error)
	DeleteFn func(ctx context.Context, instanceName string, req DeleteTemplateRequest) (SuccessResponse, error)
	FindFn   func(ctx context.Context, instanceName string) ([]map[string]any, error)
}

func (m *MockTemplateService) Create(ctx context.Context, instanceName string, req CreateTemplateRequest) (TemplateResponse, error) {
	return m.CreateFn(ctx, instanceName, req)
}

func (m *MockTemplateService) Edit(ctx context.Context, instanceName string, req EditTemplateRequest) (TemplateResponse, error) {
	return m.EditFn(ctx, instanceName, req)
}

func (m *MockTemplateService) Delete(ctx context.Context, instanceName string, req DeleteTemplateRequest) (SuccessResponse, error) {
	return m.DeleteFn(ctx, instanceName, req)
}

func (m *MockTemplateService) Find(ctx context.Context, instanceName string) ([]map[string]any, error) {
	return m.FindFn(ctx, instanceName)
}
