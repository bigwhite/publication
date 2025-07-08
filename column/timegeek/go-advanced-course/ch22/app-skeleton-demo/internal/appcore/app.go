package appcore

import (
	bizUser "app-skeleton-demo/internal/biz/user"
	clientDB "app-skeleton-demo/internal/client/database"
	foundConf "app-skeleton-demo/internal/foundational/config"
	foundLogger "app-skeleton-demo/internal/foundational/logger"
	transportHTTP "app-skeleton-demo/internal/transport/http"
	"app-skeleton-demo/pkg/lifecycle"
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// App is the core application structure that orchestrates components.
type App struct {
	appName    string
	logger     *foundLogger.Logger
	cfg        *foundConf.Config
	components []lifecycle.Component // Components that need lifecycle management
	// wg to wait for component goroutines (e.g., their Start methods) to complete,
	// especially if they might trigger an early shutdown.
	runWg sync.WaitGroup
}

// New creates and initializes the application core.
// This is where manual Dependency Injection happens.
func New(cfg *foundConf.Config) (*App, error) {
	// 1. Initialize Foundational Components (Config is already loaded by main)
	appLogger := foundLogger.New(cfg.Logger)
	// ... Initialize Metrics, Tracer here if they are also foundational and needed by others.

	// 2. Initialize Basic Service Clients
	// Note: The New function for database.Client now returns (lifecycle.Component, error)
	dbComp, err := clientDB.New(cfg.DB, appLogger)
	if err != nil {
		return nil, fmt.Errorf("appcore: failed to initialize database client: %w", err)
	}
	// We need the concrete type for DI if biz layer expects it, or cast to interface.
	// For this demo, biz.user.Service expects the DB interface.
	// database.Client must implement bizUser.DB interface.
	var dbForBiz bizUser.DB
	if concreteDB, ok := dbComp.(*clientDB.Client); ok {
		dbForBiz = concreteDB
	} else {
		return nil, fmt.Errorf("appcore: database client does not conform to expected type for biz layer")
	}

	// ... Initialize CacheClient, MQClient, etc. if they exist and are lifecycle components.

	// 3. Initialize Business Components
	userSvc := bizUser.NewService(dbForBiz, appLogger)
	// ... Initialize OrderService, etc.

	// 4. Initialize API Servers / Lifecycle Components
	// transportHTTP.New now returns lifecycle.Component
	httpSrvComp := transportHTTP.New(cfg.HTTPServer, userSvc, appLogger)

	// Collect all components that need lifecycle management
	componentsToManage := []lifecycle.Component{
		dbComp,      // Database client is now a lifecycle component
		httpSrvComp, // HTTP server is a lifecycle component
		// Add other lifecycle components like MQ consumers here
	}

	appLogger.Infof("AppCore: All components initialized and dependencies injected.")
	return &App{
		appName:    cfg.AppName,
		logger:     appLogger,
		cfg:        cfg,
		components: componentsToManage,
	}, nil
}

// Logger provides access to the app's logger.
func (a *App) Logger() *foundLogger.Logger {
	return a.logger
}

// Run starts the application and blocks until a shutdown signal is received
// or a critical component fails.
func (a *App) Run() error {
	a.logger.Infof("AppCore: Starting application: %s", a.appName)

	// Create a context that gets cancelled on OS signal (SIGINT, SIGTERM)
	// or if any critical component's Start method calls the 'appStopCauseError' function.
	appCtx, appStopCauseError := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appStopCauseError() // Important to release resources

	// Start all registered lifecycle components concurrently
	for _, comp := range a.components {
		c := comp // Capture range variable for goroutine
		a.runWg.Add(1)
		go func() {
			defer a.runWg.Done()
			a.logger.Infof("AppCore: Starting component: %s", c.Name())
			if err := c.Start(appCtx); err != nil { // Pass the app's cancellable context
				a.logger.Errorf("AppCore: Error starting component %s: %v. Initiating application shutdown.", c.Name(), err)
				appStopCauseError() // Trigger shutdown for all other components
			}
			// If Start is blocking and succeeds, it will run until appCtx is cancelled.
			// If Start is non-blocking, this goroutine will exit after Start returns.
			// We need to ensure Start methods handle appCtx cancellation properly if they are long-running.
		}()
	}
	a.logger.Infof("AppCore: All components initiated for start. Application %s is running.", a.appName)

	// Block here until the appCtx is cancelled
	<-appCtx.Done()
	errCause := context.Cause(appCtx)
	if errCause != nil && errCause != context.Canceled && errCause != context.DeadlineExceeded { // context.Canceled from signal is normal
		a.logger.Infof("AppCore: Shutdown initiated due to: %v", errCause)
	} else {
		a.logger.Infof("AppCore: Shutdown signal received or context cancelled normally.")
	}

	// --- Graceful Shutdown Procedure ---
	a.logger.Infof("AppCore: Initiating graceful stop of application %s...", a.appName)

	// Create a new context for the shutdown procedure itself, with a timeout.
	// This timeout is for the *entire* shutdown sequence of all components.
	shutdownOverallTimeout := 20 * time.Second // Configurable
	stopCtx, cancelStopCtx := context.WithTimeout(context.Background(), shutdownOverallTimeout)
	defer cancelStopCtx()

	// Stop components in reverse order of registration (LIFO).
	// This assumes a simple dependency order; more complex apps might need a dependency graph.
	for i := len(a.components) - 1; i >= 0; i-- {
		comp := a.components[i]
		a.logger.Infof("AppCore: Attempting to stop component: %s", comp.Name())

		// We pass stopCtx to each component's Stop method.
		// The component's Stop method should respect this context's deadline.
		if err := comp.Stop(stopCtx); err != nil {
			a.logger.Errorf("AppCore: Error stopping component %s: %v", comp.Name(), err)
		} else {
			a.logger.Infof("AppCore: Component %s stopped successfully.", comp.Name())
		}
	}

	// Wait for all goroutines launched by `go func() { c.Start(appCtx) }` to complete.
	// This is crucial because if a component's Start() errors out and calls appStopCauseError(),
	// its goroutine might still be running. We need to ensure all these initial goroutines
	// have finished before declaring the app fully stopped.
	a.logger.Infof("AppCore: Waiting for initial component start goroutines to complete...")
	a.runWg.Wait()
	a.logger.Infof("AppCore: All initial component start goroutines have completed.")

	a.logger.Infof("AppCore: Application %s stopped gracefully.", a.appName)

	// Check if shutdown was due to an error from context.Cause or normal signal.
	// signal.NotifyContext cancels with context.Canceled on signal.
	// If it was another error, that might be the one to return from Run().
	if errCause != nil && errCause != context.Canceled {
		return fmt.Errorf("application shutdown due to error: %w", errCause)
	}
	return nil // Graceful shutdown completed
}
