package slices1

func Evens(nums []int) []int {
	var out []int
	for _, n := range nums {
		if n%2 == 0 {
			out = append(out, n)
		}
	}
	return out
}
