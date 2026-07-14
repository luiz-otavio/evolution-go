package evolution

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

const (
	LabelFindLabelsPath  = "/label/findLabels"
	LabelHandleLabelPath = "/label/handleLabel"
)

type Label struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type LabelType string

const (
	LabelTypeChat    LabelType = "chat"
	LabelTypeMessage LabelType = "message"
)

type LabelAction string

const (
	LabelActionAdd    LabelAction = "add"
	LabelActionRemove LabelAction = "remove"
)

type HandleLabelRequest struct {
	Name   string      `json:"name"`
	Type   LabelType   `json:"type"`
	ID     string      `json:"id"`
	Action LabelAction `json:"action,omitempty"`
}

func (r HandleLabelRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	switch r.Type {
	case LabelTypeChat, LabelTypeMessage:
	default:
		return fmt.Errorf("invalid type: %q", r.Type)
	}
	if r.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

type LabelService interface {
	FindLabels(ctx context.Context, instanceName string) ([]Label, error)
	HandleLabel(ctx context.Context, instanceName string, req HandleLabelRequest) (SuccessResponse, error)
}

type labelService struct {
	http   HttpProvider
	apiKey string
}

func NewLabelService(http HttpProvider, apiKey string) LabelService {
	return labelService{http: http, apiKey: apiKey}
}

func (s labelService) FindLabels(ctx context.Context, instanceName string) ([]Label, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}

	path := buildInstancePath(LabelFindLabelsPath, instanceName)
	return executeGet[[]Label](ctx, s.http, s.apiKey, path, nil)
}

func (s labelService) HandleLabel(ctx context.Context, instanceName string, req HandleLabelRequest) (SuccessResponse, error) {
	if instanceName == "" {
		return SuccessResponse{}, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return SuccessResponse{}, fmt.Errorf("invalid handle label request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return SuccessResponse{}, fmt.Errorf("failed to marshal handle label request: %w", err)
	}

	path := buildInstancePath(LabelHandleLabelPath, instanceName)
	return executePost[SuccessResponse](ctx, s.http, s.apiKey, path, payload)
}
