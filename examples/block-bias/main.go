package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	fastEventDuration = 1 * time.Millisecond
	slowEventDuration = 1 * fastEventDuration
)

func run() error {
	runtime.SetBlockProfileRate(int(slowEventDuration.Nanoseconds()))

	var (
		done = make(chan struct{})
		wg   = &sync.WaitGroup{}
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		slowEvent(done)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fastEvent(done)
	}()
	time.Sleep(1 * time.Second)
	close(done)
	wg.Wait()

	f, err := os.Create("block.pb.gz")
	if err != nil {
		return err
	}
	defer f.Close()
	if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
		return err
	}
	return nil
}

func slowEvent(done chan struct{}) error {
	return simulateBlockEvents(slowEventDuration, done)
}

func fastEvent(done chan struct{}) error {
	return simulateBlockEvents(fastEventDuration, done)
}

func simulateBlockEvents(duration time.Duration, done chan struct{}) error {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// do nothing
		case <-done:
			return nil
		}
	}
}
