package main

// I AM NOT DONE

// Pi should be a constant. Make this compile without changing the
// printed value, and without making Pi a `var`.

import "fmt"

var Pi = 3.14159

func main() {
	Pi = 3.0 // <- this line must be removed; Pi is meant to be constant
	fmt.Println(Pi)
}
