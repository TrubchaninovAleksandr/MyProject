package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port                        string `yaml:"port"`
	URL                         string `yaml:"url"`
	MaxConnections              int    `yaml:"max_connections"`
	MinConnections              int    `yaml:"min_connections"`
	ExternalAPICoindeskBaseURL  string `yaml:"external_api_coindesk_base_url"`
	ExternalAPICoindeskTimeout  int    `yaml:"external_api_coindesk_timeout"`
	ExternalAPICoindeskCurrency string `yaml:"external_api_coindesk_currency"`
	UpdateInterval              string `yaml:"update_interval"`
	LogLevel                    string `yaml:"log_level"`
	LogFormat                   string `yaml:"log_format"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
