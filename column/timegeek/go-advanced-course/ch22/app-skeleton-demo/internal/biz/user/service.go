package user

import (
	// This import path implies concrete type, adjust if using interface from elsewhere
	"app-skeleton-demo/internal/foundational/logger"
	"context"
	"fmt"
)

// DB represents the database operations needed by UserService.
// This is an example of an interface that 'database.Client' should satisfy.
type DB interface {
	Query(ctx context.Context, query string) (string, error)
	// Add other methods your service needs from the DB client
}

// Service encapsulates business logic for users.
type Service struct {
	db     DB // Dependency is an interface, promoting loose coupling.
	logger *logger.Logger
}

// NewService creates a new UserService.
// Note: The 'db' parameter is now of type DB (interface).
func NewService(db DB, logger *logger.Logger) *Service {
	logger.Infof("BusinessComponent: User service initialized.")
	return &Service{db: db, logger: logger}
}

// GetUser retrieves user information.
func (s *Service) GetUser(ctx context.Context, id int) (string, error) {
	s.logger.Infof("UserService: Getting user with ID %d", id)
	// Use the db client (via interface) to fetch user data
	query := fmt.Sprintf("SELECT name FROM users WHERE id = %d (simulated)", id)
	userData, err := s.db.Query(ctx, query) // Call the method on the interface
	if err != nil {
		s.logger.Errorf("UserService: Failed to get user %d from DB: %v", id, err)
		return "", fmt.Errorf("user service: failed to get user %d: %w", id, err)
	}
	s.logger.Debugf("UserService: Successfully retrieved data for user %d: %s", id, userData)
	return fmt.Sprintf("User-%d [%s]", id, userData), nil
}
