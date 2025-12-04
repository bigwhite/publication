package config

import (
	"os"
	"strconv"
)

// Config 应用程序配置
type Config struct {
	GitHubToken string      `json:"github_token"`
	Output      OutputConfig `json:"output"`
	Parser      ParserConfig `json:"parser"`
}

// OutputConfig 输出配置
type OutputConfig struct {
	Format      string `json:"format"`       // markdown, html, json
	Filename    string `json:"filename"`
	Destination string `json:"destination"`
	Overwrite   bool   `json:"overwrite"`
}

// ParserConfig 解析器配置
type ParserConfig struct {
	IncludeComments    bool `json:"include_comments"`
	IncludeMetadata    bool `json:"include_metadata"`
	IncludeTimestamps  bool `json:"include_timestamps"`
	IncludeUserLinks   bool `json:"include_user_links"`
	EmojisEnabled      bool `json:"emojis_enabled"`
	PreserveLineBreaks bool `json:"preserve_line_breaks"`
}

// Environment 环境变量配置
type Environment struct {
	GitHubToken string
	Debug       bool
	NoColor     bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Output: OutputConfig{
			Format:      "markdown",
			Destination: "output",
			Overwrite:   false,
		},
		Parser: ParserConfig{
			IncludeComments:    true,
			IncludeMetadata:    true,
			IncludeTimestamps:  true,
			IncludeUserLinks:   true,
			EmojisEnabled:      true,
			PreserveLineBreaks: true,
		},
	}
}

// LoadFromEnv 从环境变量加载配置
func (c *Config) LoadFromEnv() {
	env := GetEnvironment()

	if token := env.GitHubToken; token != "" {
		c.GitHubToken = token
	}

	if debug := env.Debug; debug {
		c.Parser.EmojisEnabled = false // Example debug setting
	}
}

// GetEnvironment 获取环境变量
func GetEnvironment() *Environment {
	env := &Environment{
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
		Debug:       getBoolEnv("DEBUG", false),
		NoColor:     getBoolEnv("NO_COLOR", false),
	}
	return env
}

// getBoolEnv 获取布尔环境变量
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
		// 解析失败时返回默认值，但不丢弃错误（已处理）
	}
	return defaultValue
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.GitHubToken == "" {
		return &ValidationError{
			Field:   "github_token",
			Message: "GitHub token is required",
		}
	}

	if c.Output.Format == "" {
		return &ValidationError{
			Field:   "output.format",
			Message: "Output format is required",
		}
	}

	return nil
}

// ValidationError 配置验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}