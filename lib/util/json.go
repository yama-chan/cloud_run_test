package utils

import (
	"encoding/json"
	"net/url"
)

type Config struct {
	TemplatesURL url.URL
}

func (cfg *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"templates-url": cfg.TemplatesURL.String(),
	})
}

func (cfg *Config) UnmarshalJSON(data []byte) error {
	type Raw struct {
		TemplatesURL string `json:"templates-url"`
	}
	var raw Raw
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	u, err := url.Parse(raw.TemplatesURL)
	if err != nil {
		return err
	}
	cfg.TemplatesURL = *u
	return nil
}
