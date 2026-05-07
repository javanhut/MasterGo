package structs1

import "testing"

func TestArea(t *testing.T) {
	r := Rectangle{Width: 3, Height: 4}
	if got := Area(r); got != 12 {
		t.Errorf("Area=%v; want 12", got)
	}
}
