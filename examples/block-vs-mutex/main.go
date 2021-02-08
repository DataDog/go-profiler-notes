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

func run() error {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	aquired := make(chan struct{})
	var m sync.Mutex
	m.Lock()
	go func() {
		<-aquired
		m.Lock()
		aquired <- struct{}{}
	}()
	aquired <- struct{}{}
	time.Sleep(time.Nanosecond)
	m.Unlock()
	<-aquired

	if err := writeProfile("block"); err != nil {
		return err
	} else if err := writeProfile("mutex"); err != nil {
		return err
	}
	return nil
}

func writeProfile(name string) error {
	f, err := os.Create(name + ".pb.gz")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := pprof.Lookup(name).WriteTo(f, 0); err != nil {
		return err
	}
	return nil
}
