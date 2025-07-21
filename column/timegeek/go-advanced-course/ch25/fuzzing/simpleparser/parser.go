package simpleparser

import (
	"fmt"
	"strconv"
)

// ParseAge parses a string into an age (integer).
// It expects the age to be between 0 and 150.
func ParseAge(ageStr string) (int, error) {
	if ageStr == "" {
		return 0, fmt.Errorf("age string cannot be empty")
	}
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return 0, fmt.Errorf("not a valid integer: %w", err)
	}
	if age < 0 || age > 150 { // Let's introduce a potential bug for "> 150" for fuzzing to find
		// if age < 0 || age >= 150 { // Corrected logic for ">="
		return 0, fmt.Errorf("age %d out of reasonable range (0-149)", age)
	}
	return age, nil
}
