package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
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
	runtime.MemProfileRate = 1

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		addr := "localhost:6060"
		log.Printf("Listening on %s", addr)
		return http.ListenAndServe(addr, nil)
	})

	//g.Go(allocStuff)

	// Memory leak that leaks 1 kB / sec at a rate of 10 Hz.
	g.Go(func() error { return leakStuff(1024, 100*time.Millisecond) })
	g.Go(func() error { return allocStuff(1024, 100*time.Millisecond) })
	g.Go(func() error { return forceGc(time.Second) })

	return g.Wait()
}

var leak []*Data

// leakStuff leaks ~bytes every interval.
func leakStuff(bytes int, interval time.Duration) error {
	for {
		leak = append(leak, newData(bytes))
		time.Sleep(interval)
	}
}

// allocStuff is allocating things but not leaking them.
func allocStuff(bytes int, interval time.Duration) error {
	for {
		newData(bytes)
		time.Sleep(interval)
	}
}

// forceGc forces a GC to occur every interval.
func forceGc(interval time.Duration) error {
	for {
		time.Sleep(interval)
		runtime.GC()
	}
}

func newData(size int) *Data {
	return &Data{data: make([]byte, size)}
}

type Data struct {
	data []byte
}
