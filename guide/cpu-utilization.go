// +build ignore

package main

import (
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	pprof.StartCPUProfile(os.Stdout)
	defer pprof.StopCPUProfile()

	go cpuHog()
	go cpuHog()

	time.Sleep(time.Second)
}

func cpuHog() {
	for {
	}
}
