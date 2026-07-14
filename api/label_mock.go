package evolution

import "context"

var _ LabelService = (*MockLabelService)(nil)

type MockLabelService struct {
	FindLabelsFn  func(ctx context.Context, instanceName string) ([]Label, error)
	HandleLabelFn func(ctx context.Context, instanceName string, req HandleLabelRequest) (SuccessResponse, error)
}

func (m *MockLabelService) FindLabels(ctx context.Context, instanceName string) ([]Label, error) {
	return m.FindLabelsFn(ctx, instanceName)
}

func (m *MockLabelService) HandleLabel(ctx context.Context, instanceName string, req HandleLabelRequest) (SuccessResponse, error) {
	return m.HandleLabelFn(ctx, instanceName, req)
}
