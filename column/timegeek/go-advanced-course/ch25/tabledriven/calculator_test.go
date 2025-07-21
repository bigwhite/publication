package tabledriven

import (
	"fmt"
	"math"
	"testing"
	// "github.com/google/go-cmp/cmp"
)

// Function to be tested
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func TestDivide_TableDriven(t *testing.T) {
	testCases := []struct {
		name        string  // Test case name
		a, b        float64 // Inputs
		expectedVal float64 // Expected result value
		expectedErr string  // Expected error message substring (empty if no error)
	}{
		{
			name: "ValidDivision_PositiveNumbers",
			a:    10, b: 2,
			expectedVal: 5,
			expectedErr: "",
		},
		{
			name: "ValidDivision_NegativeResult",
			a:    -10, b: 2,
			expectedVal: -5,
			expectedErr: "",
		},
		{
			name: "DivisionByZero",
			a:    10, b: 0,
			expectedVal: 0,
			expectedErr: "division by zero",
		},
		{
			name: "ZeroDividedByNumber",
			a:    0, b: 5,
			expectedVal: 0,
			expectedErr: "",
		},
		{
			name: "FloatingPointPrecision",
			a:    1.0, b: 3.0,
			expectedVal: 0.3333333333333333,
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		currentTestCase := tc
		t.Run(currentTestCase.name, func(t *testing.T) {
			val, err := Divide(currentTestCase.a, currentTestCase.b)

			if currentTestCase.expectedErr != "" {
				if err == nil {
					t.Errorf("Divide(%f, %f): expected error containing '%s', got nil",
						currentTestCase.a, currentTestCase.b, currentTestCase.expectedErr)
				} else if err.Error() != currentTestCase.expectedErr {
					t.Errorf("Divide(%f, %f): unexpected error message: got '%v', want substring '%s'",
						currentTestCase.a, currentTestCase.b, err, currentTestCase.expectedErr)
				}
			} else {
				if err != nil {
					t.Errorf("Divide(%f, %f): unexpected error: %v", currentTestCase.a, currentTestCase.b, err)
				}
			}

			if currentTestCase.expectedErr == "" {
				const epsilon = 1e-9
				if math.Abs(val-currentTestCase.expectedVal) > epsilon {
					t.Errorf("Divide(%f, %f) = %f; want %f (within epsilon %e)",
						currentTestCase.a, currentTestCase.b, val, currentTestCase.expectedVal, epsilon)
				}
			}
		})
	}
}
