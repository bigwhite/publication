package http

import (
	// Alias for clarity
	"app-skeleton-demo/internal/foundational/config"
	"app-skeleton-demo/internal/foundational/logger"
	"app-skeleton-demo/pkg/lifecycle"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// UserServiceProvider defines the interface for user-related business logic
// that the HTTP server depends on.
type UserServiceProvider interface {
	GetUser(ctx context.Context, id int) (string, error)
}

// Server is an HTTP server that implements lifecycle.Component.
type Server struct {
	httpServer  *http.Server
	logger      *logger.Logger
	userService UserServiceProvider // Dependency is an interface
	cfg         config.HTTPServerConfig
}

// New creates a new HTTP Server lifecycle component.
func New(cfg config.HTTPServerConfig, userService UserServiceProvider, logger *logger.Logger) lifecycle.Component {
	mux := http.NewServeMux() // Using standard mux for simplicity
	s := &Server{
		logger:      logger,
		userService: userService,
		cfg:         cfg,
	}
	mux.HandleFunc("/api/user", s.handleGetUser) // Register a simple endpoint

	s.httpServer = &http.Server{
		Addr:    cfg.Addr,
		Handler: mux, // In a real app, use a router like chi, gin, or echo
		// Good practice to set Read/Write timeouts
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	logger.Infof("LifecycleComponent: HTTP server configured for address: %s", cfg.Addr)
	return s
}

// Name returns the component name.
func (s *Server) Name() string { return "HTTPServer" }

// handleGetUser is an example HTTP handler.
func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	s.logger.Debugf("HTTPServer: Received request for %s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' query parameter, must be an integer", http.StatusBadRequest)
		return
	}

	userName, err := s.userService.GetUser(r.Context(), id)
	if err != nil {
		s.logger.Errorf("HTTPServer: Error getting user data: %v", err)
		// In a real app, you might have more sophisticated error mapping to HTTP status codes
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s.logger.Infof("HTTPServer: Successfully processed request for user %d", id)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hello, %s!", userName)
}

// Start begins listening for HTTP requests.
func (s *Server) Start(ctx context.Context) error {
	s.logger.Infof("HTTPServer: Starting to listen on %s...", s.cfg.Addr)
	// Run ListenAndServe in a goroutine so Start doesn't block.
	// The context passed to Start can be used to signal an early shutdown
	// if something goes wrong during the startup phase of other components.
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("HTTPServer: ListenAndServe error: %v", err) // Fatal if server fails catastrophically
		}
	}()
	s.logger.Infof("HTTPServer: Started and listening on %s.", s.cfg.Addr)

	// It's good practice for Start to be non-blocking or to have a way
	// to signal readiness. For this demo, we assume it's ready quickly.
	// One could use the passed 'ctx' to listen for cancellation during a lengthy startup.
	// For example:
	// select {
	// case <-time.After(1 * time.Second): // Simulate readiness check
	//  	return nil
	// case <-ctx.Done():
	// 		s.logger.Errorf("HTTPServer: Start cancelled during initialization: %v", ctx.Err())
	// 		// Attempt to clean up anything started
	// 		_ = s.httpServer.Close() // or Shutdown if ListenAndServe was already called
	// 		return ctx.Err()
	// }
	return nil
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Infof("HTTPServer: Stopping...")

	// Use the context passed to Stop, which should have a timeout managed by AppCore.
	// If AppCore doesn't provide a timeout, we can add one here,
	// but it's better if the caller (AppCore) controls the overall shutdown timeout.
	// For example: shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()

	err := s.httpServer.Shutdown(ctx) // ctx here should ideally have a deadline
	if err == nil {
		s.logger.Infof("HTTPServer: Stopped gracefully.")
	} else {
		s.logger.Errorf("HTTPServer: Shutdown error: %v", err)
	}
	return err // Return the error from Shutdown
}
