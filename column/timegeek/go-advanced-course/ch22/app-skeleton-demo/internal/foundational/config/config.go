package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoggerConfig holds configuration for the logger.
type LoggerConfig struct {
	Level string `yaml:"level"`
}

// HTTPServerConfig holds configuration for the HTTP server.
type HTTPServerConfig struct {
	Addr string `yaml:"addr"`
}

// DBConfig holds configuration for the database client.
type DBConfig struct {
	DSN string `yaml:"dsn"`
}

// Config is the root configuration structure for the application.
type Config struct {
	AppName    string           `yaml:"appName"`
	Logger     LoggerConfig     `yaml:"logger"`
	HTTPServer HTTPServerConfig `yaml:"httpServer"`
	DB         DBConfig         `yaml:"db"`
}

// Load reads configuration from the given file path.
func Load(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("config: failed to read file %s: %w", filePath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: failed to unmarshal data: %w", err)
	}

	fmt.Printf("[ConfigLoader] Loaded configuration from %s for app: %s\n", filePath, cfg.AppName)
	return &cfg, nil
}
