package pointers1

import "testing"

func TestDouble(t *testing.T) {
	x := 21
	Double(&x)
	if x != 42 {
		t.Errorf("x=%d; want 42", x)
	}
}
