package goldenfiles

import (
	"fmt"
	"strings"
)

// TransformText applies a simple transformation to the input string.
func TransformText(input string) string {
	// Example transformation: Convert to uppercase and add a prefix.
	// In a real scenario, this could be much more complex.
	transformed := strings.ToUpper(input)
	return fmt.Sprintf("TRANSFORMED: %s", transformed)
}
