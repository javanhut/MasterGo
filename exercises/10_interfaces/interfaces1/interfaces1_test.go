package interfaces1

import "testing"

func TestGreetAll(t *testing.T) {
	got := GreetAll([]Greeter{English{Name: "Ada"}, Spanish{Name: "Ada"}})
	want := "Hello, Ada!\nHola, Ada!"
	if got != want {
		t.Errorf("got=%q; want %q", got, want)
	}
}
