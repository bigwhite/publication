package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port    int    `yaml:"port"`
	Timeout string `yaml:"timeout"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type AppConfig struct {
	AppName  string         `yaml:"appName"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// GlobalAppConfig holds the global application configuration.
// 注意：在现代Go应用中，推荐通过依赖注入传递配置，而非使用全局变量。
// 这里使用全局变量是为了演示早期简单阶段。
var GlobalAppConfig *AppConfig

// LoadGlobalConfig loads configuration from the given file path into GlobalAppConfig.
func LoadGlobalConfig(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("stage1: failed to read config file %s: %w", filePath, err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("stage1: failed to unmarshal config data: %w", err)
	}

	GlobalAppConfig = &cfg
	fmt.Printf("[Stage1 Config] Loaded: AppName=%s, Port=%d\n", GlobalAppConfig.AppName, GlobalAppConfig.Server.Port)
	return nil
}
