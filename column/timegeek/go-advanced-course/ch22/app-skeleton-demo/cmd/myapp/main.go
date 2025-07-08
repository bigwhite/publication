package main

import (
	"app-skeleton-demo/internal/appcore"
	"app-skeleton-demo/internal/foundational/config"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	configPath string
)

func init() {
	// Define -config flag for specifying the configuration file path.
	flag.StringVar(&configPath, "config", "configs/config.yaml", "path to config file (e.g., ./configs/config.yaml)")
}

func main() {
	flag.Parse() // Parse the command-line flags.

	// 1. Load Configuration (Foundational Component's responsibility)
	cfg, err := config.Load(configPath)
	if err != nil {
		// Use fmt for fatal startup errors as logger might not be initialized yet.
		fmt.Fprintf(os.Stderr, "FATAL: Failed to load configuration from '%s': %v\n", configPath, err)
		os.Exit(1) // Exit with a non-zero code to indicate failure.
	}

	// 2. Create Application Core instance (Dependency Injection happens within appcore.New)
	app, err := appcore.New(cfg)
	if err != nil {
		// If app creation fails, the app's logger might not be available.
		fmt.Fprintf(os.Stderr, "FATAL: Failed to create application: %v\n", err)
		os.Exit(1)
	}

	// 3. Run the Application (This is a blocking call until shutdown)
	// Errors during run (e.g., critical component failure triggering shutdown) will be handled.
	if err := app.Run(); err != nil {
		// The app.Run() method should ideally handle its own logging for runtime errors
		// or errors during shutdown. This is a final catch-all.
		// We can use the app's logger if available, otherwise fallback to fmt.
		finalLogMsg := fmt.Sprintf("FATAL: Application terminated with error: %v\n", err)
		if app != nil && app.Logger() != nil {
			app.Logger().Errorf(strings.TrimSuffix(finalLogMsg, "\n")) // Use Errorf for consistency
		} else {
			fmt.Fprint(os.Stderr, finalLogMsg)
		}
		os.Exit(1)
	}

	// If app.Run() returns nil, it implies a graceful shutdown was completed successfully.
	// The app's logger (within app.Run) should have logged "Application stopped gracefully."
	// No further action needed here for a successful exit.
	// os.Exit(0) is implicit on normal main return.
}
