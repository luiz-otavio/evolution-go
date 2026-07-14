package evolution

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

// EvolutionErrorResponse is the standard Evolution API error envelope.
// Example: {"success": false, "error": {"code": "...", "message": "..."}, "meta": {...}}
type EvolutionErrorResponse struct {
	Success bool            `json:"success"`
	Error   *EvolutionError `json:"error,omitempty"`
	Meta    *EvolutionMeta  `json:"meta,omitempty"`
}

type EvolutionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e EvolutionError) Error() string {
	if e.Code != "" {
		return e.Code + ": " + e.Message
	}
	return e.Message
}

type EvolutionMeta struct {
	Timestamp string `json:"timestamp,omitempty"`
	Path      string `json:"path,omitempty"`
	Method    string `json:"method,omitempty"`
}

// ParseEvolutionError extracts an error from a non-2xx response body.
func ParseEvolutionError(body []byte, statusCode int) error {
	var envelope EvolutionErrorResponse
	if err := jsoniter.Unmarshal(body, &envelope); err == nil && envelope.Error != nil {
		return envelope.Error
	}

	var singleError EvolutionError
	if err := jsoniter.Unmarshal(body, &singleError); err == nil && singleError.Message != "" {
		return &singleError
	}

	return fmt.Errorf("unexpected error, status code: %d, response: %s", statusCode, body)
}
