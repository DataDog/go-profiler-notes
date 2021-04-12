package main

import (
	"fmt"
	"runtime"
	"testing"
)

// BenchmarkStackDepth measures the impact of stack depth on the time it takes
// to create a stack trace.
func BenchmarkStackDepth(b *testing.B) {
	for d := 8; d <= 1024; d = d * 2 {
		b.Run(fmt.Sprintf("%d", d), func(b *testing.B) {
			callers := make([]uintptr, d*2)
			atStackDepth(d, func() {
				for i := 0; i < b.N; i++ {
					if n := runtime.Callers(1, callers); n != d {
						b.Fatalf("got=%d want=%d", n, d)
					}
				}
			})
		})
	}
}

// atStackDepth calls functions that call each other until there are depth-1
// functions on the stack and then calls fn.
//go:noinline
func atStackDepth(depth int, fn func()) {
	pcs := make([]uintptr, depth+10)
	remaining := depth - runtime.Callers(1, pcs) - 1
	if remaining < 1 {
		panic("can't simulate desired stack depth: too low")
	} else if remaining == 1 {
		fn()
	} else if f, ok := stackdepth[remaining]; !ok {
		panic("can't simulate desired stack depth: no map entry")
	} else {
		f(fn)
	}
}

// disabled: turned out to be a deadend
//func BenchmarkFunctionSize(b *testing.B) {
//var callers = make([]uintptr, 32)
//for s := 1; s <= 1024; s = s * 2 {
//b.Run(fmt.Sprintf("%d", s), func(b *testing.B) {
//called := false
//funcsize[s](0, 0, func() {
//for i := 0; i < b.N; i++ {
//runtime.Callers(0, callers)
//}
//called = true
//})
//if !called {
//b.Fatal("not called")
//}
//})
//}
//}
