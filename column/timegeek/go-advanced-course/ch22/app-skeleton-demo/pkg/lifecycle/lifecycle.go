package lifecycle

import "context"

// Component defines the interface for a manageable application component
// that has a distinct start and stop lifecycle.
type Component interface {
	Start(ctx context.Context) error // Starts the component.
	Stop(ctx context.Context) error  // Stops the component gracefully.
	Name() string                    // Returns the name of the component for logging.
}
