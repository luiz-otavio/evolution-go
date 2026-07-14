package evolution

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

const (
	ProxyFindPath = "/proxy/find"
	ProxySetPath  = "/proxy/set"
)

type ProxyProtocol string

const (
	ProxyProtocolHTTP   ProxyProtocol = "http"
	ProxyProtocolHTTPS  ProxyProtocol = "https"
	ProxyProtocolSocks4 ProxyProtocol = "socks4"
	ProxyProtocolSocks5 ProxyProtocol = "socks5"
)

type Proxy struct {
	Enabled       bool          `json:"enabled"`
	ProxyHost     string        `json:"proxyHost,omitempty"`
	ProxyPort     string        `json:"proxyPort,omitempty"`
	ProxyProtocol ProxyProtocol `json:"proxyProtocol,omitempty"`
	ProxyUsername string        `json:"proxyUsername,omitempty"`
	ProxyPassword string        `json:"proxyPassword,omitempty"`
}

type ProxyService interface {
	// Find returns the proxy configuration, or nil when no proxy is configured.
	Find(ctx context.Context, instanceName string) (*Proxy, error)
	Set(ctx context.Context, instanceName string, req Proxy) (SuccessResponse, error)
}

type proxyService struct {
	http   HttpProvider
	apiKey string
}

func NewProxyService(http HttpProvider, apiKey string) ProxyService {
	return proxyService{http: http, apiKey: apiKey}
}

func (s proxyService) Find(ctx context.Context, instanceName string) (*Proxy, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}

	path := buildInstancePath(ProxyFindPath, instanceName)
	return executeGet[*Proxy](ctx, s.http, s.apiKey, path, nil)
}

func (s proxyService) Set(ctx context.Context, instanceName string, req Proxy) (SuccessResponse, error) {
	if instanceName == "" {
		return SuccessResponse{}, fmt.Errorf("instanceName is required")
	}

	payload, err := jsoniter.Marshal(req)
	if err != nil {
		return SuccessResponse{}, fmt.Errorf("failed to marshal proxy request: %w", err)
	}

	path := buildInstancePath(ProxySetPath, instanceName)
	return executePost[SuccessResponse](ctx, s.http, s.apiKey, path, payload)
}
