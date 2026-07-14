package evolution

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

// executeGet performs a GET request and decodes the response body into T.
func executeGet[T any](ctx context.Context, http HttpProvider, apiKey, path string, query map[string]string) (T, error) {
	var out T

	resp, err := http.Get(ctx, path, buildAuthHeaders(apiKey), query)
	if err != nil {
		return out, fmt.Errorf("failed to execute GET %s: %w", path, err)
	}
	defer fasthttp.ReleaseResponse(resp)

	return decodeResponse[T](resp)
}

// executePost performs a POST request with a JSON payload and decodes the response into T.
func executePost[T any](ctx context.Context, http HttpProvider, apiKey, path string, payload []byte) (T, error) {
	var out T

	resp, err := http.Post(ctx, path, buildJSONHeaders(apiKey), payload)
	if err != nil {
		return out, fmt.Errorf("failed to execute POST %s: %w", path, err)
	}
	defer fasthttp.ReleaseResponse(resp)

	return decodeResponse[T](resp)
}

// executeDelete performs a DELETE request with an optional JSON payload and decodes the response into T.
func executeDelete[T any](ctx context.Context, http HttpProvider, apiKey, path string, payload []byte) (T, error) {
	var out T

	resp, err := http.Delete(ctx, path, buildJSONHeaders(apiKey), payload)
	if err != nil {
		return out, fmt.Errorf("failed to execute DELETE %s: %w", path, err)
	}
	defer fasthttp.ReleaseResponse(resp)

	return decodeResponse[T](resp)
}

// decodeResponse validates the status code and unmarshals a successful body into T.
func decodeResponse[T any](resp *fasthttp.Response) (T, error) {
	var out T

	if !isSuccess(resp.StatusCode()) {
		return out, ParseEvolutionError(resp.Body(), resp.StatusCode())
	}

	body := resp.Body()
	if len(body) == 0 {
		return out, nil
	}

	if err := jsoniter.Unmarshal(body, &out); err != nil {
		return out, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return out, nil
}
