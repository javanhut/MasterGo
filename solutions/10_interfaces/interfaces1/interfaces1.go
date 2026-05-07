package interfaces1

import "fmt"

type Greeter interface {
	Greet() string
}

type English struct{ Name string }

func (e English) Greet() string {
	return fmt.Sprintf("Hello, %s!", e.Name)
}

type Spanish struct{ Name string }

func (s Spanish) Greet() string {
	return fmt.Sprintf("Hola, %s!", s.Name)
}

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
