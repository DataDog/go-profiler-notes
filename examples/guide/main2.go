package main

import (
	"os"
	"runtime/trace"
	"time"
)

func main() {
	trace.Start(os.Stdout)
	defer trace.Stop()

	go b()
	a()
}

func a() {
	start := time.Now()
	for time.Since(start) < time.Second {
	}
}

func b() {
	start := time.Now()
	for time.Since(start) < time.Second {
	}
}
