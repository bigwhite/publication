package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Error("DefaultConfig() returned nil")
	}

	if cfg.Output.Format != "markdown" {
		t.Errorf("Expected format 'markdown', got %s", cfg.Output.Format)
	}

	if cfg.Output.Destination != "output" {
		t.Errorf("Expected destination 'output', got %s", cfg.Output.Destination)
	}

	if cfg.Output.Overwrite != false {
		t.Errorf("Expected overwrite false, got %v", cfg.Output.Overwrite)
	}
}

func TestGetEnvironment(t *testing.T) {
	// 保存原始环境变量
	originalToken := os.Getenv("GITHUB_TOKEN")
	originalDebug := os.Getenv("DEBUG")
	originalNoColor := os.Getenv("NO_COLOR")

	// 测试清理
	defer func() {
		os.Setenv("GITHUB_TOKEN", originalToken)
		os.Setenv("DEBUG", originalDebug)
		os.Setenv("NO_COLOR", originalNoColor)
	}()

	// 设置测试环境变量
	os.Setenv("GITHUB_TOKEN", "test-token")
	os.Setenv("DEBUG", "true")
	os.Setenv("NO_COLOR", "1")

	env := GetEnvironment()

	if env.GitHubToken != "test-token" {
		t.Errorf("Expected GitHubToken 'test-token', got %s", env.GitHubToken)
	}

	if env.Debug != true {
		t.Errorf("Expected Debug true, got %v", env.Debug)
	}

	if env.NoColor != true {
		t.Errorf("Expected NoColor true, got %v", env.NoColor)
	}
}

func TestGetBoolEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue bool
		setValue     string
		want         bool
	}{
		{
			name:         "unset env var",
			key:          "UNSET_VAR",
			defaultValue: false,
			want:         false,
		},
		{
			name:         "unset env var with true default",
			key:          "UNSET_VAR",
			defaultValue: true,
			want:         true,
		},
		{
			name:         "true string",
			key:          "TEST_VAR",
			defaultValue: false,
			setValue:     "true",
			want:         true,
		},
		{
			name:         "false string",
			key:          "TEST_VAR",
			defaultValue: true,
			setValue:     "false",
			want:         false,
		},
		{
			name:         "number 1",
			key:          "TEST_VAR",
			defaultValue: false,
			setValue:     "1",
			want:         true,
		},
		{
			name:         "number 0",
			key:          "TEST_VAR",
			defaultValue: true,
			setValue:     "0",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理环境变量
			os.Unsetenv(tt.key)
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
			}
			defer os.Unsetenv(tt.key)

			got := getBoolEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getBoolEnv(%s, %v) = %v, want %v", tt.key, tt.defaultValue, got, tt.want)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				GitHubToken: "test-token",
				Output: OutputConfig{
					Format: "markdown",
				},
			},
			wantErr: false,
		},
		{
			name: "missing github token",
			cfg: &Config{
				Output: OutputConfig{
					Format: "markdown",
				},
			},
			wantErr: true,
		},
		{
			name: "missing output format",
			cfg: &Config{
				GitHubToken: "test-token",
				Output: OutputConfig{
					Format: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "test_field",
		Message: "test message",
	}

	if err.Error() != "test message" {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), "test message")
	}

	if err.Field != "test_field" {
		t.Errorf("ValidationError.Field = %v, want %v", err.Field, "test_field")
	}
}