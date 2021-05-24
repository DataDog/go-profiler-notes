package main

import (
	"os"
	"runtime/pprof"
)

func main() {
	foo()
}

func foo() {
	bar()
}
func bar() {
	debug2, _ := os.Create("frames.txt")
	debug0, _ := os.Create("frames.pb.gz")

	pprof.Lookup("goroutine").WriteTo(debug2, 2)
	pprof.Lookup("goroutine").WriteTo(debug0, 0)
}
