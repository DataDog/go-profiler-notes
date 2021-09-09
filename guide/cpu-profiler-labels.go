// +build ignore

package main

import (
	"context"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	pprof.StartCPUProfile(os.Stdout)
	defer pprof.StopCPUProfile()

	go work(context.Background(), "alice")
	go work(context.Background(), "bob")

	time.Sleep(50 * time.Millisecond)
}

func work(ctx context.Context, user string) {
	labels := pprof.Labels("user", user)
	pprof.Do(ctx, labels, func(_ context.Context) {
		go backgroundWork()
		directWork()
	})
}

func directWork() {
	for {
	}
}

func backgroundWork() {
	for {
	}
}
