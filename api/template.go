package evolution

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

const (
	TemplateCreatePath = "/template/create"
	TemplateDeletePath = "/template/delete"
	TemplateEditPath   = "/template/edit"
	TemplateFindPath   = "/template/find"
)

type TemplateCategory string

const (
	TemplateCategoryMarketing      TemplateCategory = "MARKETING"
	TemplateCategoryUtility        TemplateCategory = "UTILITY"
	TemplateCategoryAuthentication TemplateCategory = "AUTHENTICATION"
)

type CreateTemplateRequest struct {
	Name       string           `json:"name"`
	Category   TemplateCategory `json:"category"`
	Language   string           `json:"language"`
	Components []map[string]any `json:"components"`
}

func (r CreateTemplateRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	switch r.Category {
	case TemplateCategoryMarketing, TemplateCategoryUtility, TemplateCategoryAuthentication:
	default:
		return fmt.Errorf("invalid category: %q", r.Category)
	}
	if r.Language == "" {
		return fmt.Errorf("language is required")
	}
	if len(r.Components) == 0 {
		return fmt.Errorf("components is required")
	}
	return nil
}

type EditTemplateRequest struct {
	Name       string           `json:"name"`
	Components []map[string]any `json:"components,omitempty"`
}

func (r EditTemplateRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type DeleteTemplateRequest struct {
	Name string `json:"name"`
}

func (r DeleteTemplateRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type TemplateResponse struct {
	Success  bool           `json:"success"`
	Template map[string]any `json:"template"`
}

type TemplateService interface {
	Create(ctx context.Context, instanceName string, req CreateTemplateRequest) (TemplateResponse, error)
	Edit(ctx context.Context, instanceName string, req EditTemplateRequest) (TemplateResponse, error)
	Delete(ctx context.Context, instanceName string, req DeleteTemplateRequest) (SuccessResponse, error)
	Find(ctx context.Context, instanceName string) ([]map[string]any, error)
}

type templateService struct {
	http   HttpProvider
	apiKey string
}

func NewTemplateService(http HttpProvider, apiKey string) TemplateService {
	return templateService{http: http, apiKey: apiKey}
}

func (s templateService) Create(ctx context.Context, instanceName string, req CreateTemplateRequest) (TemplateResponse, error) {
	if instanceName == "" {
		return TemplateResponse{}, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return TemplateResponse{}, fmt.Errorf("invalid create template request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return TemplateResponse{}, fmt.Errorf("failed to marshal create template request: %w", err)
	}

	path := buildInstancePath(TemplateCreatePath, instanceName)
	return executePost[TemplateResponse](ctx, s.http, s.apiKey, path, payload)
}

func (s templateService) Edit(ctx context.Context, instanceName string, req EditTemplateRequest) (TemplateResponse, error) {
	if instanceName == "" {
		return TemplateResponse{}, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return TemplateResponse{}, fmt.Errorf("invalid edit template request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return TemplateResponse{}, fmt.Errorf("failed to marshal edit template request: %w", err)
	}

	path := buildInstancePath(TemplateEditPath, instanceName)
	return executePost[TemplateResponse](ctx, s.http, s.apiKey, path, payload)
}

func (s templateService) Delete(ctx context.Context, instanceName string, req DeleteTemplateRequest) (SuccessResponse, error) {
	if instanceName == "" {
		return SuccessResponse{}, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return SuccessResponse{}, fmt.Errorf("invalid delete template request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return SuccessResponse{}, fmt.Errorf("failed to marshal delete template request: %w", err)
	}

	path := buildInstancePath(TemplateDeletePath, instanceName)
	return executeDelete[SuccessResponse](ctx, s.http, s.apiKey, path, payload)
}

func (s templateService) Find(ctx context.Context, instanceName string) ([]map[string]any, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}

	path := buildInstancePath(TemplateFindPath, instanceName)
	return executeGet[[]map[string]any](ctx, s.http, s.apiKey, path, nil)
}
