package evolution

import "context"

var _ CallService = (*MockCallService)(nil)

type MockCallService struct {
	OfferFn func(ctx context.Context, instanceName string, req OfferCallRequest) (SuccessResponse, error)
}

func (m *MockCallService) Offer(ctx context.Context, instanceName string, req OfferCallRequest) (SuccessResponse, error) {
	return m.OfferFn(ctx, instanceName, req)
}
