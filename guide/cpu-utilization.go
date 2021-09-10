// +build ignore

package main

import (
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	file, _ := os.Create("./cpu-utilization.pprof")
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	go cpuHog()
	go cpuHog()

	time.Sleep(time.Second)
}

func cpuHog() {
	for {
	}
}
