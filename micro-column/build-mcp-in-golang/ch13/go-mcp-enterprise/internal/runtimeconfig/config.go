package runtimeconfig

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config 保存 MCP HTTP 服务的运行配置。
type Config struct {
	Environment     string
	Addr            string
	MCPPath         string
	BearerToken     string
	AllowedHosts    []string
	AllowedOrigins  []string
	ShutdownTimeout time.Duration
}

// Load 从环境变量读取配置，并进行启动前校验。
func Load() (Config, error) {
	cfg := Config{
		Environment:     envOrDefault("APP_ENV", "development"),
		Addr:            envOrDefault("MCP_ADDR", "127.0.0.1:8080"),
		MCPPath:         envOrDefault("MCP_PATH", "/mcp"),
		BearerToken:     os.Getenv("MCP_AUTH_TOKEN"),
		AllowedHosts:    splitList(os.Getenv("MCP_ALLOWED_HOSTS")),
		AllowedOrigins:  splitList(os.Getenv("MCP_ALLOWED_ORIGINS")),
		ShutdownTimeout: 15 * time.Second,
	}

	if raw := os.Getenv("MCP_SHUTDOWN_TIMEOUT"); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return Config{}, fmt.Errorf(
				"invalid MCP_SHUTDOWN_TIMEOUT: %w", err,
			)
		}
		if d <= 0 {
			return Config{}, fmt.Errorf(
				"MCP_SHUTDOWN_TIMEOUT must be positive",
			)
		}
		cfg.ShutdownTimeout = d
	}

	switch cfg.Environment {
	case "development", "testing", "production":
	default:
		return Config{}, fmt.Errorf(
			"unsupported APP_ENV %q", cfg.Environment,
		)
	}

	if !strings.HasPrefix(cfg.MCPPath, "/") {
		return Config{}, fmt.Errorf(
			"MCP_PATH must start with /",
		)
	}

	if cfg.Environment == "production" {
		if cfg.BearerToken == "" {
			return Config{}, fmt.Errorf(
				"MCP_AUTH_TOKEN is required in production",
			)
		}
		if len(cfg.AllowedHosts) == 0 {
			return Config{}, fmt.Errorf(
				"MCP_ALLOWED_HOSTS is required in production",
			)
		}
	}

	// 开发环境默认只允许本机访问。
	// 生产环境必须由部署配置明确指定允许的 Host。
	if cfg.Environment != "production" &&
		len(cfg.AllowedHosts) == 0 {
		cfg.AllowedHosts = []string{
			"127.0.0.1:8080",
			"localhost:8080",
		}
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func splitList(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	var result []string
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}
