// +build ignore

package main

import (
	"fmt"
)

func main() {
	fmt.Println(*add(23, 42))
}

func add(a, b int) *int {
	sum := a + b
	return &sum
}
