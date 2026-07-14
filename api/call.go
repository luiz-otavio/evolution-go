package evolution

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

const CallOfferPath = "/call/offer"

type OfferCallRequest struct {
	Number string         `json:"number"`
	Offer  map[string]any `json:"offer"`
}

func (r OfferCallRequest) Validate() error {
	if r.Number == "" {
		return fmt.Errorf("number is required")
	}
	if len(r.Offer) == 0 {
		return fmt.Errorf("offer is required")
	}
	return nil
}

type CallService interface {
	Offer(ctx context.Context, instanceName string, req OfferCallRequest) (SuccessResponse, error)
}

type callService struct {
	http   HttpProvider
	apiKey string
}

func NewCallService(http HttpProvider, apiKey string) CallService {
	return callService{http: http, apiKey: apiKey}
}

func (s callService) Offer(ctx context.Context, instanceName string, req OfferCallRequest) (SuccessResponse, error) {
	if instanceName == "" {
		return SuccessResponse{}, fmt.Errorf("instanceName is required")
	}
	if err := req.Validate(); err != nil {
		return SuccessResponse{}, fmt.Errorf("invalid offer call request: %w", err)
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return SuccessResponse{}, fmt.Errorf("failed to marshal offer call request: %w", err)
	}

	path := buildInstancePath(CallOfferPath, instanceName)
	return executePost[SuccessResponse](ctx, s.http, s.apiKey, path, payload)
}
