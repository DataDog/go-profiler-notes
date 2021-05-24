package main

import (
	"fmt"
)

func main() {
	foo()
	foobar()
}

func foo() {
	bar()
}

func bar() {
	fmt.Printf("Hello from bar\n")
	panic("oh no")
}

func foobar() {
	fmt.Printf("Hello from foobar\n")
}
