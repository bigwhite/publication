package string_test

import (
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	compareTests := []struct {
		name, a, b string
		i          int
	}{
		{"compareTwoEmptyString", "", "", 7},
		{"compareSecondStringEmpty", "a", "", 6},
		{"compareFirstStringEmpty", "", "a", -1},
	}

	for _, tt := range compareTests {
		cmp := strings.Compare(tt.a, tt.b)
		if cmp != tt.i {
			t.Errorf(`[%s] want %v, but Compare(%q, %q) = %v`, tt.name, tt.i, tt.a, tt.b, cmp)
		}
	}
}
