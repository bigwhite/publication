package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

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

var currentAppConfig *AppConfig // 包级私有变量

// Load loads configuration from the given file path.
func Load(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("stage2: failed to read config file %s: %w", filePath, err)
	}
	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("stage2: failed to unmarshal config data: %w", err)
	}
	currentAppConfig = &cfg
	fmt.Printf("[Stage2 Config] Loaded: AppName=%s, Port=%d\n", currentAppConfig.AppName, currentAppConfig.Server.Port)
	return nil
}

// GetByPath retrieves a configuration value by a dot-separated path.
// This is a simplified example using reflection and has performance implications.
func GetByPath(path string) (interface{}, bool) {
	if currentAppConfig == nil {
		return nil, false
	}

	parts := strings.Split(path, ".")
	v := reflect.ValueOf(*currentAppConfig) // Dereference pointer

	for _, part := range parts {
		if v.Kind() == reflect.Ptr { // Should not happen with currentAppConfig being value
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return nil, false
		}

		found := false
		// Case-insensitive field matching for flexibility, or use tags
		var matchedField reflect.Value
		for i := 0; i < v.NumField(); i++ {
			fieldName := v.Type().Field(i).Name
			yamlTag := v.Type().Field(i).Tag.Get("yaml")
			if strings.EqualFold(fieldName, part) || yamlTag == part {
				matchedField = v.Field(i)
				found = true
				break
			}
		}

		if !found {
			return nil, false
		}
		v = matchedField
	}

	if v.IsValid() && v.CanInterface() {
		return v.Interface(), true
	}
	return nil, false
}

// GetString provides a typed getter for string values.
func GetString(path string) (string, bool) {
	val, ok := GetByPath(path)
	if !ok {
		return "", false
	}
	s, ok := val.(string)
	return s, ok
}

// GetInt provides a typed getter for int values.
func GetInt(path string) (int, bool) {
	val, ok := GetByPath(path)
	if !ok {
		return 0, false
	}
	// Handle if it's already int or can be parsed from string
	switch v := val.(type) {
	case int:
		return v, true
	case int64: // YAML might unmarshal numbers as int64
		return int(v), true
	case float64: // YAML might unmarshal numbers as float64
		return int(v), true
	case string:
		i, err := strconv.Atoi(v)
		if err == nil {
			return i, true
		}
	}
	return 0, false
}
