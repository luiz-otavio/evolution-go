package evolution

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/valyala/fasthttp"
)

func FillIfNotEmpty(params map[string]string, key, value string) {
	if value != "" {
		params[key] = value
	}
}

func FillIfNotZero(params map[string]string, key string, value int) {
	if value != 0 {
		params[key] = strconv.Itoa(value)
	}
}

// buildInstancePath appends an URL-escaped instance name to a base path.
func buildInstancePath(base, instanceName string) string {
	return fmt.Sprintf("%s/%s", base, url.PathEscape(instanceName))
}

// isSuccess reports whether the HTTP status code indicates success.
func isSuccess(statusCode int) bool {
	return statusCode >= fasthttp.StatusOK && statusCode < fasthttp.StatusMultipleChoices
}
