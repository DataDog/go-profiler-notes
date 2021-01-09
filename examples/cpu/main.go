package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
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
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		addr := "localhost:6060"
		log.Printf("Listening on %s", addr)
		return http.ListenAndServe(addr, nil)
	})

	g.Go(func() error { return computeSleepLoop(1000000, 10*time.Millisecond) })

	return g.Wait()
}

func computeSleepLoop(n int, sleep time.Duration) error {
	for {
		compute(n)
		time.Sleep(sleep)
	}
	return nil
}

func compute(n int) int64 {
	var sum int64
	for i := 0; i < n; i++ {
		sum += int64(i) / 2
		sum += int64(i) / 3
		sum += int64(i) / 4
		sum += int64(i) / 5
		sum += int64(i) / 6
		sum += int64(i) / 7
		sum += int64(i) / 8
	}
	return sum
}
