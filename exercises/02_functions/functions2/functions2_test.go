package functions2

import "testing"

func TestDivMod(t *testing.T) {
	cases := []struct{ a, b, q, r int }{
		{10, 3, 3, 1},
		{20, 5, 4, 0},
		{7, 2, 3, 1},
		{0, 9, 0, 0},
	}
	for _, c := range cases {
		q, r := DivMod(c.a, c.b)
		if q != c.q || r != c.r {
			t.Errorf("DivMod(%d,%d) = (%d,%d); want (%d,%d)", c.a, c.b, q, r, c.q, c.r)
		}
	}
}
