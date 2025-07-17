package config

// FeatureFlags holds boolean flags for features.
type FeatureFlags struct {
	NewAuth         bool `mapstructure:"newAuth"`
	ExperimentalAPI bool `mapstructure:"experimentalApi"`
}

// ServerConfig holds server specific configurations.
type ServerConfig struct {
	Port           int `mapstructure:"port"`
	TimeoutSeconds int `mapstructure:"timeoutSeconds"`
}

// AppConfig is the root configuration structure.
type AppConfig struct {
	AppName      string       `mapstructure:"appName"`
	LogLevel     string       `mapstructure:"logLevel"`
	FeatureFlags FeatureFlags `mapstructure:"featureFlags"`
	Server       ServerConfig `mapstructure:"server"`
}
