package maps1

import "testing"

func TestWordCount(t *testing.T) {
	got := WordCount("go go gopher go")
	if got["go"] != 3 || got["gopher"] != 1 {
		t.Errorf("got=%v; want go=3, gopher=1", got)
	}
	if len(WordCount("")) != 0 {
		t.Errorf("empty input should give empty map")
	}
}
