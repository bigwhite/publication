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

var (
	currentAppConfig *AppConfig
	configIndex      = make(map[string]interface{})
)

// Load loads configuration and builds an index for fast access.
func Load(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("stage3: failed to read config file %s: %w", filePath, err)
	}
	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("stage3: failed to unmarshal config data: %w", err)
	}
	currentAppConfig = &cfg
	buildIndex("", reflect.ValueOf(*currentAppConfig)) // Build index after loading
	fmt.Printf("[Stage3 Config] Loaded and indexed: AppName=%s, Port=%d\n", currentAppConfig.AppName, currentAppConfig.Server.Port)
	return nil
}

func buildIndex(prefix string, rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return
	}

	typ := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		fieldStruct := typ.Field(i)
		fieldVal := rv.Field(i)

		// Use YAML tag as key part, fallback to lowercase field name
		keyPart := fieldStruct.Tag.Get("yaml")
		if keyPart == "" {
			keyPart = strings.ToLower(fieldStruct.Name)
		}
		if keyPart == "-" { // Skip fields Wärme `yaml:"-"`
			continue
		}

		currentPath := keyPart
		if prefix != "" {
			currentPath = prefix + "." + keyPart
		}

		if fieldVal.Kind() == reflect.Struct {
			buildIndex(currentPath, fieldVal)
		} else if fieldVal.CanInterface() {
			configIndex[currentPath] = fieldVal.Interface()
		}
	}
}

// GetByPathFromIndex retrieves a value from the pre-built index.
func GetByPathFromIndex(path string) (interface{}, bool) {
	if currentAppConfig == nil { // Ensure config was loaded
		return nil, false
	}
	val, ok := configIndex[strings.ToLower(path)] // Normalize path to lowercase for lookup
	return val, ok
}

// GetString, GetInt methods now use GetByPathFromIndex
func GetString(path string) (string, bool) {
	val, ok := GetByPathFromIndex(path)
	if !ok {
		return "", false
	}
	s, ok := val.(string)
	return s, ok
}

func GetInt(path string) (int, bool) {
	val, ok := GetByPathFromIndex(path)
	if !ok {
		return 0, false
	}
	switch v := val.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true // Common for YAML numbers
	case float64:
		return int(v), true // Common for YAML numbers
	case string:
		i, err := strconv.Atoi(v)
		if err == nil {
			return i, true
		}
	}
	return 0, false
}
