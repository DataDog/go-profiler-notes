// +build ignore

package main

import (
	"os"
	"runtime/pprof"
)

func main() {
	pprof.StartCPUProfile(os.Stdout)
	defer pprof.StopCPUProfile()

	belowLimit()
	aboveLimit()
}

func belowLimit() {
	atDepth(32, cpuHog)
}

func aboveLimit() {
	atDepth(64, cpuHog)
}

func cpuHog() {
	for i := 0; i < 5000*1000*1000; i++ {
	}
}

func atDepth(n int, fn func()) {
	if n > 0 {
		atDepth(n-1, fn)
	} else {
		fn()
	}
}
