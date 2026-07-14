package evolution

import "fmt"

// EvolutionConfig holds the configuration required to build an EvolutionClient.
type EvolutionConfig struct {
	BaseURL string
	APIKey  string
}

func (c EvolutionConfig) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("base URL is not set, please set the base URL in the configuration")
	}
	if c.APIKey == "" {
		return fmt.Errorf("api key is not set, please set the API key in the configuration")
	}

	return nil
}
