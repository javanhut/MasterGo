package maps1

import "strings"

func WordCount(s string) map[string]int {
	out := make(map[string]int)
	for _, w := range strings.Fields(s) {
		out[w]++
	}
	return out
}
