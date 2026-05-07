package goroutines1

import "testing"

func TestSumConcurrent(t *testing.T) {
	got := SumConcurrent([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	if got != 55 {
		t.Errorf("sum=%d; want 55", got)
	}
	if SumConcurrent(nil) != 0 {
		t.Errorf("empty slice should sum to 0")
	}
}
