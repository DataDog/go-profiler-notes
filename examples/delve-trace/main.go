package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	for i := 0; i < 10; i++ {
		foo(i)
	}
	return nil
}

func foo(i int) int {
	return i * 2
}
