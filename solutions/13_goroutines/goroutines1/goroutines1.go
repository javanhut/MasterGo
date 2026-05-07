package goroutines1

import "sync"

func SumConcurrent(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	mid := len(nums) / 2
	var wg sync.WaitGroup
	var left, right int

	wg.Add(2)
	go func() {
		defer wg.Done()
		for _, n := range nums[:mid] {
			left += n
		}
	}()
	go func() {
		defer wg.Done()
		for _, n := range nums[mid:] {
			right += n
		}
	}()
	wg.Wait()
	return left + right
}
