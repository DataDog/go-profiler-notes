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

	g.Go(func() error { return computeSum(1000000, 10*time.Millisecond) })

	return g.Wait()
}

func computeSum(to int, sleep time.Duration) error {
	for {
		var sum int64
		for i := 0; i < to; i++ {
			sum += int64(i) / 2
			sum += int64(i) / 3
			sum += int64(i) / 4
			sum += int64(i) / 5
			sum += int64(i) / 6
			sum += int64(i) / 7
			sum += int64(i) / 8
		}
		time.Sleep(sleep)
	}

	return nil
}
