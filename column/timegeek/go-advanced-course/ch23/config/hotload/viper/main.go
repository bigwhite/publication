package main

import (
	"ch23/config/hotload/viper/hotloader" // Adjust import path
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Path to the directory containing the config file, and config file name/type
	// In a real app, these might come from flags or env vars
	configDir := "./"   // Assumes configs dir is relative to where binary is run
	configName := "app" // app.yaml
	configType := "yaml"

	// Initialize and start watching the configuration
	sharedAppConfig, _, err := hotloader.InitAndWatchConfig(configDir, configName, configType)
	if err != nil {
		log.Fatalf("Failed to initialize config loader: %v", err)
	}

	// Simulate an application using the configuration periodically
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Channel to listen for OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Application started. Waiting for config changes or OS signal...")
	log.Printf("Try modifying '%s/%s.%s' and save it.\n", configDir, configName, configType)

	for {
		select {
		case <-ticker.C:
			// Access configuration safely via the SharedConfig methods
			currentLogLevel := sharedAppConfig.GetLogLevel()
			isNewAuthEnabled := sharedAppConfig.IsFeatureEnabled("newAuth")

			log.Printf("Current LogLevel: %s, NewAuth Feature: %t (AppName: %s, Port: %d)\n",
				currentLogLevel,
				isNewAuthEnabled,
				sharedAppConfig.Get().AppName,     // Example of getting the whole struct copy
				sharedAppConfig.Get().Server.Port) // Then accessing fields

			// In a real application, components would use this sharedAppConfig
			// or be notified of changes to re-initialize or adjust behavior.

		case s := <-quit:
			log.Printf("Received OS signal: %s. Shutting down...", s)
			// Perform any cleanup قبل exiting
			log.Println("Application shut down gracefully.")
			return
		}
	}
}
