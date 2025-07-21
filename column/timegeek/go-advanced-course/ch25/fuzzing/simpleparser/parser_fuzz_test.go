package simpleparser

import (
	"testing"
	"unicode/utf8" // We might check for valid strings if ParseAge expects it
)

func FuzzParseAge(f *testing.F) {
	// Add seed corpus: valid ages, edge cases, invalid inputs
	f.Add("0")
	f.Add("1")
	f.Add("149") // Edge case for upper bound (based on "> 150" bug, this is valid)
	f.Add("150") // This input should ideally be valid, but might fail due to the bug
	f.Add("-1")
	f.Add("abc")      // Not an integer
	f.Add("")         // Empty string
	f.Add("1000")     // Out of range
	f.Add(" 77 ")     // String with spaces (Atoi handles this)
	f.Add("\x80test") // Invalid UTF-8 prefix - strconv.Atoi might handle or error early

	// The Fuzzing execution function
	f.Fuzz(func(t *testing.T, ageStr string) {
		// Call the function being fuzzed
		age, err := ParseAge(ageStr)

		// Define our expectations / invariants
		if err != nil {
			t.Logf("ParseAge(%q) returned error: %v (this might be expected for fuzzed inputs)", ageStr, err)
			return
		}

		if age < -1000 || age > 1000 { // Arbitrary broad check for successfully parsed ages
			t.Errorf("ParseAge(%q) resulted in an unexpected age %d without error", ageStr, age)
		}

		if utf8.ValidString(ageStr) {
			if age < 0 || age >= 150 { // This should not happen if err == nil
				t.Errorf("Successfully parsed age %d for input %q is out of the *absolute* expected range 0-150", age, ageStr)
			}
		}
	})
}
