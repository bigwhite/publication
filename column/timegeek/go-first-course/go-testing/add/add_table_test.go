package add

import "testing"

func TestAddWithTable(t *testing.T) {
	cases := []struct {
		name string
		a    int
		b    int
		r    int
	}{
		{"2+3", 2, 3, 5},
		{"2+0", 2, 0, 2},
		{"2+(-2)", 2, -2, 0},
		//... ...
	}

	for _, caze := range cases {
		got := Add(caze.a, caze.b)
		if got != caze.r {
			t.Errorf("%s got %d, want %d", caze.name, got, caze.r)
		}
	}
}
