package generics1

import (
	"strconv"
	"testing"
)

func TestMapInts(t *testing.T) {
	got := Map([]int{1, 2, 3}, func(n int) int { return n * n })
	if len(got) != 3 || got[0] != 1 || got[1] != 4 || got[2] != 9 {
		t.Errorf("got=%v; want [1 4 9]", got)
	}
}

func TestMapIntToString(t *testing.T) {
	// This will only compile once Map is generic.
	got := Map([]int{1, 2, 3}, strconv.Itoa)
	if len(got) != 3 || got[0] != "1" || got[1] != "2" || got[2] != "3" {
		t.Errorf("got=%v; want [\"1\" \"2\" \"3\"]", got)
	}
}
