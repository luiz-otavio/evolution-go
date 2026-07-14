package evolution

// buildAuthHeaders builds the headers required to authenticate against the
// Evolution API. Authentication is performed via the `apikey` header, which can
// hold either the global API key or an instance-specific token.
func buildAuthHeaders(apiKey string) map[string]string {
	return map[string]string{
		"apikey": apiKey,
	}
}

// buildJSONHeaders builds the headers for a JSON request, including auth.
func buildJSONHeaders(apiKey string) map[string]string {
	return map[string]string{
		"apikey":       apiKey,
		"Content-Type": "application/json",
	}
}
