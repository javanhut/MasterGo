package if1

import "testing"

func TestClassify(t *testing.T) {
	cases := map[int]string{-3: "negative", 0: "zero", 5: "positive"}
	for in, want := range cases {
		if got := Classify(in); got != want {
			t.Errorf("Classify(%d)=%q; want %q", in, got, want)
		}
	}
}
