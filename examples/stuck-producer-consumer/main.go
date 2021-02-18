package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
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
	if err := os.Setenv("DD_PROFILING_WAIT_PROFILE", "yes"); err != nil {
		return err
	} else if err := profiler.Start(
		profiler.WithService("stuck-producer-consumer"),
		profiler.WithVersion(time.Now().String()),
		profiler.WithPeriod(60*time.Second),
		profiler.WithProfileTypes(
			profiler.GoroutineProfile,
		),
	); err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())

	workCh := make(chan struct{})
	go consumer(workCh)
	go producer(workCh)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// force gc on a regular basis to make sure g.waitsince gets populated.
			fmt.Printf("gc\n")
			runtime.GC()
		case sig := <-c:
			fmt.Printf("sig: %s\n", sig)
			return nil
		}
	}
}

func producer(workCh chan<- struct{}) {
	for {
		select {
		case workCh <- struct{}{}:
			if rand.Int63n(10) == 0 {
				takeNap()
			}
		default:
		}
	}
}

func consumer(workCh <-chan struct{}) {
	for {
		<-workCh
		if rand.Int63n(10) == 0 {
			takeNap()
		}
	}
}

func takeNap() {
	fmt.Printf("taking a nap\n")
	var forever chan struct{}
	<-forever
}
