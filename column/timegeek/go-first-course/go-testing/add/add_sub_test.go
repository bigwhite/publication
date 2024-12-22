package add

import "testing"

func TestAddWithSubtest(t *testing.T) {
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
		t.Run(caze.name, func(t *testing.T) {
			got := Add(caze.a, caze.b)
			if got != caze.r {
				t.Errorf("got %d, want %d", got, caze.r)
			}
		})
	}
}
