package evolution

import "context"

var _ BusinessService = (*MockBusinessService)(nil)

type MockBusinessService struct {
	GetCatalogFn     func(ctx context.Context, instanceName string, req BusinessNumberRequest) (CatalogResponse, error)
	GetCollectionsFn func(ctx context.Context, instanceName string, req BusinessNumberRequest) (CollectionsResponse, error)
}

func (m *MockBusinessService) GetCatalog(ctx context.Context, instanceName string, req BusinessNumberRequest) (CatalogResponse, error) {
	return m.GetCatalogFn(ctx, instanceName, req)
}

func (m *MockBusinessService) GetCollections(ctx context.Context, instanceName string, req BusinessNumberRequest) (CollectionsResponse, error) {
	return m.GetCollectionsFn(ctx, instanceName, req)
}
