package strings1

import "testing"

func TestShout(t *testing.T) {
	cases := map[string]string{
		"hello":   "HELLO!",
		"Go":      "GO!",
		"golings": "GOLINGS!",
	}
	for in, want := range cases {
		if got := Shout(in); got != want {
			t.Errorf("Shout(%q)=%q; want %q", in, got, want)
		}
	}
}
