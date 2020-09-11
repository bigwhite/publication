package string_test

import (
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	compareTests := map[string]struct {
		a, b string
		i    int
	}{
		`compareTwoEmptyString`:     {"", "", 0},
		`compareSecondParamIsEmpty`: {"a", "", 1},
		`compareFirstParamIsEmpty`:  {"", "a", -1},
	}

	for name, tt := range compareTests {
		t.Run(name, func(t *testing.T) {
			cmp := strings.Compare(tt.a, tt.b)
			if cmp != tt.i {
				t.Errorf(`want %v, but Compare(%q, %q) = %v`, tt.i, tt.a, tt.b, cmp)
			}
		})
	}
}
