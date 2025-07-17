package main

import (
	"ch23/config/stage4/config"
	"fmt"
	"log"
	"os" // For setting environment variables for testing

	"github.com/spf13/pflag"
)

var (
	// Define flags using pflag for Viper integration
	port    = pflag.IntP("port", "p", 0, "HTTP server port (from pflag)")
	appName = pflag.String("appname", "", "Application name (from pflag)")
	// It's better to use a config file path flag if viper needs to load a specific file
	// rather than individual config item flags, but this shows direct flag binding.
)

func main() {
	pflag.Parse() // Must be called before Viper binds flags

	// Simulate setting environment variables for testing
	os.Setenv("MYAPP_SERVER_TIMEOUT", "60s") // This will be mapped to server.timeout if struct has it
	os.Setenv("MYAPP_DATABASE_DSN", "env_user:env_pass@tcp(env_host:3306)/env_db")

	// For Viper to find "configs/app.yaml", ensure "configs" is a search path
	// and "app" is the ConfigName, "yaml" is ConfigType.
	// The configPath here is one of the directories Viper will search.
	cfg, err := config.LoadConfigWithViper("./", "app", "yaml")
	if err != nil {
		log.Fatalf("Failed to load config with Viper: %v", err)
	}

	fmt.Println("\n--- Final Configuration ---")
	fmt.Printf("AppName: %s\n", cfg.AppName)
	fmt.Printf("Server Port: %d\n", cfg.Server.Port)
	fmt.Printf("Server Timeout: %s\n", cfg.Server.Timeout)
	fmt.Printf("Database DSN: %s\n", cfg.Database.DSN)

	fmt.Println("\n--- Accessing directly from Viper instance ---")
	// Note: Viper keys are case-insensitive by default for Get, but exact for Unmarshal via mapstructure tags
	fmt.Printf("Viper AppName: %s\n", config.ViperInstance.GetString("appName"))
	fmt.Printf("Viper Server Port: %d\n", config.ViperInstance.GetInt("server.port"))
	fmt.Printf("Viper DB DSN: %s\n", config.ViperInstance.GetString("database.dsn"))
	fmt.Printf("Viper Server Timeout (from env): %s\n", config.ViperInstance.GetString("server.timeout"))

	// Example of how pflag overrides if a value was set
	if *port > 0 { // pflag gives 0 if not set for IntP
		fmt.Printf("Pflag --port was set, it has high priority: %d (reflected in cfg.Server.Port if names match or bound)\n", *port)
		// To make pflag 'port' directly affect 'cfg.Server.Port',
		// you'd typically bind it via v.BindPFlag("server.port", pflag.Lookup("port"))
		// or ensure viper keys and pflag names match what mapstructure expects (if not using Set explicitly)
		// For simplicity in LoadConfigWithViper, we used v.BindPFlags(pflag.CommandLine) which binds all parsed flags.
		// Viper then tries to match them. If a pflag name is 'port', it might set top-level 'port' in Viper.
		// The Unmarshal to struct relies on mapstructure tags.
	}
}
