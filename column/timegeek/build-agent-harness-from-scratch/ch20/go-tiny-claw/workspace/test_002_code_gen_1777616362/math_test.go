package math

import "testing"

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        3,
			b:        4,
			expected: 12,
		},
		{
			name:     "negative numbers",
			a:        -2,
			b:        -5,
			expected: 10,
		},
		{
			name:     "positive and negative",
			a:        6,
			b:        -3,
			expected: -18,
		},
		{
			name:     "zero with positive",
			a:        0,
			b:        7,
			expected: 0,
		},
		{
			name:     "zero with negative",
			a:        0,
			b:        -4,
			expected: 0,
		},
		{
			name:     "zero with zero",
			a:        0,
			b:        0,
			expected: 0,
		},
		{
			name:     "negative and positive",
			a:        -8,
			b:        2,
			expected: -16,
		},
		{
			name:     "large numbers",
			a:        1000,
			b:        1000,
			expected: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Multiply(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}