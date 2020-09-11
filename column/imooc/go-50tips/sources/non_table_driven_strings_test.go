package string_test

import (
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	var a, b string
	var i int

	a, b = "", ""
	i = 0
	cmp := strings.Compare(a, b)
	if cmp != i {
		t.Errorf(`want %v, but Compare(%q, %q) = %v`, i, a, b, cmp)
	}

	a, b = "a", ""
	i = 1
	cmp = strings.Compare(a, b)
	if cmp != i {
		t.Errorf(`want %v, but Compare(%q, %q) = %v`, i, a, b, cmp)
	}

	a, b = "", "a"
	i = -1
	cmp = strings.Compare(a, b)
	if cmp != i {
		t.Errorf(`want %v, but Compare(%q, %q) = %v`, i, a, b, cmp)
	}
}
