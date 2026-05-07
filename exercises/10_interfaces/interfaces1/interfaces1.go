package interfaces1

// I AM NOT DONE

import "fmt"

// Greeter is implemented by anything that can produce a greeting.
type Greeter interface {
	Greet() string
}

// English is a Greeter; its Greet method should return "Hello, <Name>!".
type English struct {
	Name string
}

// TODO: implement English.Greet so it satisfies Greeter.

// Spanish is also a Greeter; its Greet method should return "Hola, <Name>!".
type Spanish struct {
	Name string
}

// TODO: implement Spanish.Greet so it satisfies Greeter.

// GreetAll returns the greetings of everyone in the slice, joined by "\n".
func GreetAll(gs []Greeter) string {
	out := ""
	for i, g := range gs {
		if i > 0 {
			out += "\n"
		}
		out += g.Greet()
	}
	return out
}

var _ = fmt.Sprintf
