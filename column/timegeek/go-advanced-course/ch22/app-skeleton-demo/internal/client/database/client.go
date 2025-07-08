package database

import (
	"app-skeleton-demo/internal/foundational/config"
	"app-skeleton-demo/internal/foundational/logger"
	"app-skeleton-demo/pkg/lifecycle" // Assuming lifecycle is in pkg
	"context"
	"fmt"
	"time"
)

// Client simulates a database client and implements lifecycle.Component.
type Client struct {
	dsn    string
	logger *logger.Logger
	// In a real app: *sql.DB or other DB driver instance
}

// New creates a new simulated database client.
func New(cfg config.DBConfig, logger *logger.Logger) (lifecycle.Component, error) {
	logger.Infof("BasicServiceClient: Initializing DB client with DSN: %s (simulated)", cfg.DSN)
	// Simulate some setup latency
	time.Sleep(50 * time.Millisecond)
	return &Client{dsn: cfg.DSN, logger: logger}, nil
}

// Name returns the component name.
func (c *Client) Name() string { return "DatabaseClient" }

// Start simulates starting the database client (e.g., establishing connections).
func (c *Client) Start(ctx context.Context) error {
	c.logger.Infof("DatabaseClient: Starting...")
	// Simulate a check or a connection ping
	select {
	case <-time.After(100 * time.Millisecond):
		c.logger.Infof("DatabaseClient: Started and connection pool ready (simulated).")
		return nil
	case <-ctx.Done():
		c.logger.Errorf("DatabaseClient: Start cancelled: %v", ctx.Err())
		return ctx.Err()
	}
}

// Stop simulates closing database connections.
func (c *Client) Stop(ctx context.Context) error {
	c.logger.Infof("DatabaseClient: Stopping...")
	// Simulate closing connections
	select {
	case <-time.After(200 * time.Millisecond):
		c.logger.Infof("DatabaseClient: Stopped and connections closed (simulated).")
		return nil
	case <-ctx.Done():
		c.logger.Errorf("DatabaseClient: Stop timed out or cancelled: %v", ctx.Err())
		return ctx.Err()
	}
}

// Query simulates executing a database query.
func (c *Client) Query(ctx context.Context, query string) (string, error) {
	c.logger.Debugf("DatabaseClient: Executing query: %s", query)
	// Simulate query latency and context cancellation
	select {
	case <-time.After(150 * time.Millisecond):
		return fmt.Sprintf("result for '%s' from DB", query), nil
	case <-ctx.Done():
		c.logger.Errorf("DatabaseClient: Query context cancelled for: %s: %v", query, ctx.Err())
		return "", ctx.Err()
	}
}
