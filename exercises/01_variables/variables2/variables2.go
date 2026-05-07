package main

// I AM NOT DONE

// `:=` only works inside a function. Fix the package-level declaration
// so this compiles. The value should remain "golings".

import "fmt"

name := "golings"

func main() {
	fmt.Println("Hello,", name)
}
