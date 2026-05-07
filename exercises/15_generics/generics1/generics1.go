package generics1

// I AM NOT DONE

// Map applies f to each element of s and returns a new slice with the results.
// Make this generic over the input element type T and output element type U.

func Map(s []int, f func(int) int) []int {
	out := make([]int, len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}
