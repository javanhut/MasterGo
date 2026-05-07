package goroutines1

// I AM NOT DONE

// SumConcurrent should split nums into two halves, sum each half in its
// own goroutine, and return the total. Use sync.WaitGroup (or a channel)
// to wait for both goroutines before returning.

func SumConcurrent(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
