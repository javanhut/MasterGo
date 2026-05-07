package main

import "fmt"

func main() {
	var i int = 10
	var f float64 = 3.5

	sum := float64(i) + f
	rounded := int(sum)

	fmt.Println(sum, rounded)
}
