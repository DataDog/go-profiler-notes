// +build ignore

package main

import (
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func main() {
	runtime.SetCPUProfileRate(800)

	pprof.StartCPUProfile(os.Stdout)
	defer pprof.StopCPUProfile()

	go cpuHog()
	time.Sleep(time.Second)
}

func cpuHog() {
	for {
	}
}
