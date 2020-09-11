package string_test

import (
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	compareTests := []struct {
		a, b string
		i    int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", -1},
		{"abc", "abc", 0},
		{"ab", "abc", -1},
		{"abc", "ab", 1},
		{"x", "ab", 1},
		{"ab", "x", -1},
		{"x", "a", 1},
		{"b", "x", -1},
		// test runtime·memeq's chunked implementation
		{"abcdefgh", "abcdefgh", 0},
		{"abcdefghi", "abcdefghi", 0},
		{"abcdefghi", "abcdefghj", -1},
	}

	for _, tt := range compareTests {
		cmp := strings.Compare(tt.a, tt.b)
		if cmp != tt.i {
			t.Errorf(`want %v, but Compare(%q, %q) = %v`, tt.i, tt.a, tt.b, cmp)
		}
	}
}
