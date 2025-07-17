package hotloader

import (
	"ch23/config/hotload/viper/config" // Adjust import path if your module name is different
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// SharedConfig holds the current application configuration.
// It's protected by a RWMutex for concurrent access.
type SharedConfig struct {
	mu  sync.RWMutex
	cfg *config.AppConfig
}

// Get returns a copy of the current config to avoid race conditions on the caller's side
// if they hold onto it while it's being updated. Or, caller can use its methods.
func (sc *SharedConfig) Get() config.AppConfig {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	if sc.cfg == nil { // Should not happen if initialized properly
		return config.AppConfig{}
	}
	return *sc.cfg // Return a copy
}

// Update atomically updates the shared configuration.
func (sc *SharedConfig) Update(newCfg *config.AppConfig) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.cfg = newCfg
	log.Printf("[HotLoader] Configuration updated: %+v\n", *sc.cfg)
}

// GetLogLevel is an example of a type-safe getter.
func (sc *SharedConfig) GetLogLevel() string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	if sc.cfg == nil {
		return "info" // Default
	}
	return sc.cfg.LogLevel
}

// IsFeatureEnabled is another example.
func (sc *SharedConfig) IsFeatureEnabled(featureKey string) bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	if sc.cfg == nil {
		return false
	}
	// This is a simplified check; a real app might have more robust feature flag access
	switch featureKey {
	case "newAuth":
		return sc.cfg.FeatureFlags.NewAuth
	case "experimentalApi":
		return sc.cfg.FeatureFlags.ExperimentalAPI
	default:
		return false
	}
}

// InitAndWatchConfig initializes Viper, loads initial config, and starts watching for changes.
// It returns a SharedConfig instance that can be safely accessed by the application.
func InitAndWatchConfig(configDir string, configName string, configType string) (*SharedConfig, *viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath(configDir)
	v.SetConfigName(configName)
	v.SetConfigType(configType)

	// Initial read of the config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, nil, fmt.Errorf("failed to read config: %w", err)
	}

	var initialCfg config.AppConfig
	if err := v.Unmarshal(&initialCfg); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal initial config: %w", err)
	}
	log.Printf("[HotLoader] Initial configuration loaded: %+v\n", initialCfg)

	sharedCfg := &SharedConfig{
		cfg: &initialCfg,
	}

	// Start watching for config changes in a separate goroutine
	go func() {
		v.WatchConfig() // This blocks internally or uses a goroutine, check Viper docs
		log.Println("[HotLoader] Viper WatchConfig started.")
		v.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("[HotLoader] Config file changed: %s (Op: %s)\n", e.Name, e.Op)

			// It's crucial to re-read and re-unmarshal the config
			// as v.ReadInConfig() is needed to refresh Viper's internal state from the file.
			if err := v.ReadInConfig(); err != nil {
				log.Printf("[HotLoader] Error re-reading config after change: %v", err)
				// Decide on error handling: revert, keep old, or panic?
				// For simplicity, we keep the old config.
				return
			}

			var newCfgInstance config.AppConfig
			if err := v.Unmarshal(&newCfgInstance); err != nil {
				log.Printf("[HotLoader] Error unmarshaling new config: %v", err)
				// Keep the old config if unmarshaling fails
				return
			}
			sharedCfg.Update(&newCfgInstance)

			// Here, you would typically notify other parts of your application
			// that the configuration has changed, e.g., via channels or callbacks.
			// For example: configUpdateChannel <- newCfgInstance
		})
	}()

	return sharedCfg, v, nil // Return viper instance if needed for direct access, otherwise just sharedCfg
}
