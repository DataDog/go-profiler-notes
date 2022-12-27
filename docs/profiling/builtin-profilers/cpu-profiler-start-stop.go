//go:build ignore

package main

import (
	"os"
	"runtime/pprof"
)

func main() {
	file, _ := os.Create("./cpu.pprof")
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()
}
