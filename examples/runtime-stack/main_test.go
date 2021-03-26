package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func BenchmarkRuntimeStack(b *testing.B) {
	buf := make([]byte, 1024*1024*64)
	for g := 1; g <= 1024*1024; g = g * 2 {
		g := g
		name := fmt.Sprintf("%d goroutines", g)

		b.Run(name, func(b *testing.B) {
			initalRoutines := runtime.NumGoroutine()

			readyCh := make(chan struct{})
			stopCh := make(chan struct{})
			for i := 0; i < g; i++ {
				go atStackDepth(16, func() {
					defer func() { stopCh <- struct{}{} }()
					readyCh <- struct{}{}
				})
				<-readyCh
			}

			gotRoutines := runtime.NumGoroutine() - initalRoutines
			if gotRoutines != g {
				b.Logf("want %d goroutines, but got %d", g, gotRoutines)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				runtime.Stack(buf, true)
			}
			b.StopTimer()
			for i := 0; i < g; i++ {
				<-stopCh
			}
			start := time.Now()
			for i := 0; ; i++ {
				if runtime.NumGoroutine() == initalRoutines {
					break
				}
				time.Sleep(20 * time.Millisecond)
				if time.Since(start) > 10*time.Second {
					b.Fatalf("%d goroutines still running, want %d", runtime.NumGoroutine(), initalRoutines)
				}
			}
		})
	}
}

func atStackDepth(depth int, fn func()) {
	pcs := make([]uintptr, depth*10)
	n := runtime.Callers(1, pcs)
	if n > depth {
		panic("depth exceeded")
	} else if n < depth {
		atStackDepth(depth, fn)
		return
	}

	fn()
}
