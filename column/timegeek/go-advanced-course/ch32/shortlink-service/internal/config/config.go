// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 是应用的总配置结构体
type Config struct {
	AppName string        `mapstructure:"appName"`
	Server  ServerConfig  `mapstructure:"server"`
	Store   StoreConfig   `mapstructure:"store"`
	Tracing TracingConfig `mapstructure:"tracing"`
}

// ServerConfig 包含HTTP服务器相关的配置
type ServerConfig struct {
	Port            string        `mapstructure:"port"`
	LogLevel        string        `mapstructure:"logLevel"`
	LogFormat       string        `mapstructure:"logFormat"` // "text" or "json"
	ReadTimeout     time.Duration `mapstructure:"readTimeout"`
	WriteTimeout    time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout     time.Duration `mapstructure:"idleTimeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

// StoreConfig 包含与存储相关的配置
type StoreConfig struct {
	Type string `mapstructure:"type"`
	// DSN string `mapstructure:"dsn"` // Example for Postgres
}

// TracingConfig 包含与分布式追踪相关的配置
type TracingConfig struct {
	Enabled      bool    `mapstructure:"enabled"`
	OTELEndpoint string  `mapstructure:"otelEndpoint"`
	SampleRatio  float64 `mapstructure:"sampleRatio"`
}

// LoadConfig 从指定路径加载配置文件并允许环境变量覆盖。
func LoadConfig(configSearchPath string, configName string, configType string) (cfg Config, err error) {
	v := viper.New()

	// 1. 设置默认值
	v.SetDefault("appName", "ShortLinkServiceAppDefault")
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.logLevel", "info")
	v.SetDefault("server.logFormat", "json")
	v.SetDefault("server.readTimeout", "5s")
	v.SetDefault("server.writeTimeout", "10s")
	v.SetDefault("server.idleTimeout", "120s")
	v.SetDefault("server.shutdownTimeout", "15s")
	v.SetDefault("store.type", "memory")
	v.SetDefault("tracing.enabled", true)
	v.SetDefault("tracing.otelEndpoint", "localhost:4317") // OTel Collector gRPC
	v.SetDefault("tracing.sampleRatio", 1.0)

	// 2. 设置配置文件查找路径、名称和类型
	if configSearchPath != "" {
		v.AddConfigPath(configSearchPath)
	}
	v.AddConfigPath(".")
	v.SetConfigName(configName)
	v.SetConfigType(configType)

	// 3. 尝试读取配置文件
	if errRead := v.ReadInConfig(); errRead != nil {
		if _, ok := errRead.(viper.ConfigFileNotFoundError); !ok {
			return cfg, fmt.Errorf("config: failed to read config file: %w", errRead)
		}
		fmt.Fprintf(os.Stderr, "[ConfigLoader] Warning: Config file '%s.%s' not found. Using defaults/env vars.\n", configName, configType)
	} else {
		fmt.Fprintf(os.Stdout, "[ConfigLoader] Using config file: %s\n", v.ConfigFileUsed())
	}

	// 4. 启用环境变量覆盖
	v.SetEnvPrefix("SHORTLINK")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 5. Unmarshal到Config结构体
	if errUnmarshal := v.Unmarshal(&cfg); errUnmarshal != nil {
		return cfg, fmt.Errorf("config: unable to decode all configurations into struct: %w", errUnmarshal)
	}

	fmt.Fprintf(os.Stdout, "[ConfigLoader] Configuration loaded successfully. AppName: %s\n", cfg.AppName)
	return cfg, nil
}
