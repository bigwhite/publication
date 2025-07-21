package subtests

import (
	"fmt"
	"testing"
)

// ParseInput is a function we want to test with subtests.
func ParseInput(input string, strictMode bool) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input cannot be empty")
	}
	if strictMode && len(input) > 10 {
		return "", fmt.Errorf("input too long in strict mode (max 10 chars)")
	}
	return "parsed: " + input, nil
}

func TestParseInput_WithSubtests(t *testing.T) {
	t.Log("Setting up for ParseInput tests...")
	defer t.Log("Tearing down after ParseInput tests...") // Example of shared teardown

	t.Run("EmptyInput", func(t *testing.T) {
		// t.Parallel() // This subtest could run in parallel if independent
		_, err := ParseInput("", false)
		if err == nil {
			t.Error("Expected error for empty input, got nil")
		} else {
			t.Logf("Got expected error for empty input: %v", err)
		}
	})

	t.Run("ValidInputNonStrict", func(t *testing.T) {
		t.Parallel()
		input := "hello world" // More than 10 chars, but non-strict
		expected := "parsed: " + input
		result, err := ParseInput(input, false)
		if err != nil {
			t.Errorf("Expected no error for valid input in non-strict mode, got %v", err)
		}
		if result != expected {
			t.Errorf("Expected result '%s', got '%s'", expected, result)
		}
	})

	t.Run("StrictChecks", func(t *testing.T) { // A group for strict mode tests
		t.Parallel() // This group itself can be parallel with other top-level t.Run

		t.Run("InputTooLongInStrictMode", func(t *testing.T) {
			t.Parallel()
			input := "thisisareallylonginput" // More than 10 chars
			_, err := ParseInput(input, true) // strictMode = true
			if err == nil {
				t.Error("Expected error for too long input in strict mode, got nil")
			} else {
				expectedErrorMsg := "input too long in strict mode (max 10 chars)"
				if err.Error() != expectedErrorMsg {
					t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
				}
				t.Logf("Got expected error for long input in strict mode: %v", err)
			}
		})

		t.Run("ValidShortInputInStrictMode", func(t *testing.T) {
			t.Parallel()
			input := "short"
			expected := "parsed: " + input
			result, err := ParseInput(input, true) // strictMode = true
			if err != nil {
				t.Errorf("Expected no error for short input in strict mode, got %v", err)
			}
			if result != expected {
				t.Errorf("Expected result '%s', got '%s'", expected, result)
			}
		})
	})
}
