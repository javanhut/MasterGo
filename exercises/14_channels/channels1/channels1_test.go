package channels1

import "testing"

func TestProduce(t *testing.T) {
	var got []int
	for v := range Produce(5) {
		got = append(got, v)
	}
	want := []int{1, 2, 3, 4, 5}
	if len(got) != len(want) {
		t.Fatalf("got=%v; want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got=%v; want %v", got, want)
		}
	}
}
