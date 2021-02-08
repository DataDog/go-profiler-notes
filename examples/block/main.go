package main

import (
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

var foo []string

func main() {
	demonstrateSleep()

	f, err := os.Create("block.pb.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
		panic(err)
	}
}

func demonstrateSleep() {
	runtime.SetBlockProfileRate(1)
	<-time.After(time.Millisecond)
}

func demonstrateSelect() {
	runtime.SetBlockProfileRate(1)

	ch1 := make(chan struct{}, 0)
	ch2 := make(chan struct{}, 1)
	ch3 := make(chan struct{}, 0)

	go func() {
		ch2 <- struct{}{}
	}()

	time.Sleep(20 * time.Millisecond)

	select {
	case <-ch1:
	case <-ch2:
	case <-ch3:
	}
}

func demonstrateSampling() {
	runtime.SetBlockProfileRate(int(40 * time.Microsecond.Nanoseconds()))
	for i := 0; i < 10000; i++ {
		blockMutex(10 * time.Microsecond)
	}
}

func blockMutex(d time.Duration) {
	m := &sync.Mutex{}
	m.Lock()
	go func() {
		spinSleep(d)
		m.Unlock()
	}()
	m.Lock()
}

// spinSleep is a more accurate version of time.Sleep() for short sleep
// durations. Accuracy seems to be ~35ns.
func spinSleep(d time.Duration) {
	start := time.Now()
	n := 0
	for time.Since(start) < d {
		n++
	}
}
