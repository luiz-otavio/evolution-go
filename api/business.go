package evolution

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

const (
	BusinessGetCatalogPath     = "/business/getCatalog"
	BusinessGetCollectionsPath = "/business/getCollections"
)

type BusinessNumberRequest struct {
	Number string `json:"number"`
}

func (r BusinessNumberRequest) Validate() error {
	if r.Number == "" {
		return fmt.Errorf("number is required")
	}
	return nil
}

type CatalogResponse struct {
	Success bool           `json:"success"`
	Catalog map[string]any `json:"catalog"`
}

type CollectionsResponse struct {
	Success     bool             `json:"success"`
	Collections []map[string]any `json:"collections"`
}

type BusinessService interface {
	GetCatalog(ctx context.Context, instanceName string, req BusinessNumberRequest) (CatalogResponse, error)
	GetCollections(ctx context.Context, instanceName string, req BusinessNumberRequest) (CollectionsResponse, error)
}

type businessService struct {
	http   HttpProvider
	apiKey string
}

func NewBusinessService(http HttpProvider, apiKey string) BusinessService {
	return businessService{http: http, apiKey: apiKey}
}

func (s businessService) GetCatalog(ctx context.Context, instanceName string, req BusinessNumberRequest) (CatalogResponse, error) {
	if instanceName == "" {
		return CatalogResponse{}, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return CatalogResponse{}, fmt.Errorf("invalid get catalog request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return CatalogResponse{}, fmt.Errorf("failed to marshal get catalog request: %w", err)
	}

	path := buildInstancePath(BusinessGetCatalogPath, instanceName)
	return executePost[CatalogResponse](ctx, s.http, s.apiKey, path, payload)
}

func (s businessService) GetCollections(ctx context.Context, instanceName string, req BusinessNumberRequest) (CollectionsResponse, error) {
	if instanceName == "" {
		return CollectionsResponse{}, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return CollectionsResponse{}, fmt.Errorf("invalid get collections request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return CollectionsResponse{}, fmt.Errorf("failed to marshal get collections request: %w", err)
	}

	path := buildInstancePath(BusinessGetCollectionsPath, instanceName)
	return executePost[CollectionsResponse](ctx, s.http, s.apiKey, path, payload)
}
