package main

// I AM NOT DONE

// Go does not auto-convert between numeric types. Make this compile
// by inserting the right explicit conversions. Don't change the
// declared types of i and f.

import "fmt"

func main() {
	var i int = 10
	var f float64 = 3.5

	sum := i + f          // <- type mismatch
	rounded := int(sum)   // ok once sum is fixed

	fmt.Println(sum, rounded)
}
