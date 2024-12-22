package add

import (
	"testing"
)

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5
	if got != want {
		t.Errorf("want %d, but got %d", want, got)
	}
}
