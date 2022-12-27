//go:build ignore

package main

import (
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

func main() {
	file, _ := os.Create("./trace.trace")
	trace.Start(file)
	defer trace.Stop()

	pprof.StartCPUProfile(io.Discard)
	defer pprof.StopCPUProfile()

	go cpuHog()
	time.Sleep(1 * time.Second)
}

func cpuHog() {
	for i := 0; ; i++ {
		fmt.Fprintf(io.Discard, "%d", i)
	}
}
