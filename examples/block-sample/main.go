package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	runtime.SetBlockProfileRate(int((40 * time.Microsecond).Nanoseconds()))
	done := make(chan struct{})
	g := errgroup.Group{}
	g.Go(func() error {
		return eventA(done)
	})
	g.Go(func() error {
		return eventB(done)
	})
	time.Sleep(time.Second)
	close(done)
	if err := g.Wait(); err != nil {
		return err
	}

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

func eventA(done chan struct{}) error {
	return simulateBlockEvents(20*time.Microsecond, done)
}

func eventB(done chan struct{}) error {
	return simulateBlockEvents(40*time.Microsecond, done)
}

const tolerance = 1.1

func simulateBlockEvents(meanDuration time.Duration, done chan struct{}) error {
	var (
		prev   time.Time
		sum    time.Duration
		count  int
		ticker = time.NewTicker(meanDuration)
	)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			if !prev.IsZero() {
				sum += now.Sub(prev)
				count += 1
				if count > 1000 {
					actualMean := float64(sum) / float64(count)
					max := tolerance * float64(meanDuration)
					min := float64(meanDuration) / tolerance
					if actualMean <= min || actualMean >= max {
						return fmt.Errorf("low clock accuracy: got=%s want=%s", time.Duration(actualMean), meanDuration)
					}
				}
			}
			prev = now
		case <-done:
			return nil
		}
	}
}
