package evolution

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/valyala/fasthttp"
)

type HttpProvider interface {
	Post(ctx context.Context, path string, headers map[string]string, payload []byte) (*fasthttp.Response, error)
	Get(ctx context.Context, path string, headers map[string]string, queryParams map[string]string) (*fasthttp.Response, error)
	Put(ctx context.Context, path string, headers map[string]string, payload []byte) (*fasthttp.Response, error)
	Delete(ctx context.Context, path string, headers map[string]string, payload []byte) (*fasthttp.Response, error)
}

type httpProvider struct {
	baseUrl string
	client  *fasthttp.Client
}

func NewHttpProvider(baseUrl string) HttpProvider {
	return httpProvider{
		baseUrl: baseUrl,
		client:  &fasthttp.Client{},
	}
}

func (h httpProvider) Post(ctx context.Context, path string, headers map[string]string, payload []byte) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()

	if payload != nil {
		req.SetBody(payload)
	}

	req.SetRequestURI(h._buildUrlByPath(path))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if err := h.client.DoTimeout(req, resp, h._getTimeoutFromCtx(ctx)); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, fmt.Errorf("failed to execute POST request: %w", err)
	}

	return resp, nil
}

func (h httpProvider) Get(ctx context.Context, path string, headers map[string]string, queryParams map[string]string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(h._buildUrlByPath(buildPathWithQuery(path, queryParams)))
	req.Header.SetMethod(fasthttp.MethodGet)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if err := h.client.DoTimeout(req, resp, h._getTimeoutFromCtx(ctx)); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, fmt.Errorf("failed to execute GET request: %w", err)
	}

	return resp, nil
}

func (h httpProvider) Put(ctx context.Context, path string, headers map[string]string, payload []byte) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(h._buildUrlByPath(path))
	req.Header.SetMethod(fasthttp.MethodPut)
	req.Header.SetContentType("application/json")
	req.SetBody(payload)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if err := h.client.DoTimeout(req, resp, h._getTimeoutFromCtx(ctx)); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, fmt.Errorf("failed to execute PUT request: %w", err)
	}

	return resp, nil
}

func (h httpProvider) Delete(ctx context.Context, path string, headers map[string]string, payload []byte) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(h._buildUrlByPath(path))
	req.Header.SetMethod(fasthttp.MethodDelete)
	req.Header.SetContentType("application/json")

	if payload != nil {
		req.SetBody(payload)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if err := h.client.DoTimeout(req, resp, h._getTimeoutFromCtx(ctx)); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, fmt.Errorf("failed to execute DELETE request: %w", err)
	}

	return resp, nil
}

func (h httpProvider) _getTimeoutFromCtx(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return 60 * time.Second
	}
	timeout := time.Until(deadline)
	if timeout <= 0 {
		return time.Millisecond
	}
	return timeout
}

func (h httpProvider) _buildUrlByPath(path string) string {
	return fmt.Sprintf("%s%s", h.baseUrl, path)
}

func buildPathWithQuery(path string, queryParams map[string]string) string {
	if len(queryParams) == 0 {
		return path
	}

	keys := make([]string, 0, len(queryParams))
	for key := range queryParams {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	for _, key := range keys {
		args.Set(key, queryParams[key])
	}

	return fmt.Sprintf("%s?%s", path, args.String())
}
