package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
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
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Ошибка чтения конфигурации: %v", err)
		return nil,
			fmt.Errorf("ошибка чтения конфига: %w", err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Ошибка распаковки в структуру: %v", err)
		return nil,
			fmt.Errorf("ошибка распаковки конфига: %w", err)
	}
	return &cfg, nil
}
