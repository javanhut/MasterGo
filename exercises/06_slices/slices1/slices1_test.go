package slices1

import (
	"reflect"
	"testing"
)

func TestEvens(t *testing.T) {
	got := Evens([]int{1, 2, 3, 4, 5, 6})
	want := []int{2, 4, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Evens=%v; want %v", got, want)
	}
	if e := Evens([]int{1, 3, 5}); len(e) != 0 {
		t.Errorf("Evens(odds)=%v; want empty", e)
	}
}
