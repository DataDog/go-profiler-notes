package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if err := startProfiler(); err != nil {
		return err
	}
	defer profiler.Stop()

	for {
		time.Sleep(10 * time.Millisecond)
	}
}

func startProfiler() error {
	return profiler.Start(
		profiler.WithService("datadog-example"),
		profiler.WithEnv("localhost"),
		profiler.WithVersion("1.0"),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			profiler.MutexProfile,
			profiler.GoroutineProfile,
			profiler.BlockProfile,
		),
	)
}
