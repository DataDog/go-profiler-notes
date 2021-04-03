package main

import (
	"runtime"
	"testing"
)

var callers = make([]uintptr, 32)

func BenchmarkCallers(b *testing.B) {
	b.Run("short", func(b *testing.B) {
		short(func() {
			for i := 0; i < b.N; i++ {
				runtime.Callers(0, callers)
			}
		})
	})

	b.Run("loop", func(b *testing.B) {
		loop(func() {
			for i := 0; i < b.N; i++ {
				runtime.Callers(0, callers)
			}
		})
	})

	b.Run("long", func(b *testing.B) {
		long(func() {
			for i := 0; i < b.N; i++ {
				runtime.Callers(0, callers)
			}
		})
	})
}

func short(fn func()) {
	other()
	fn()
}

func loop(fn func()) {
	for i := 0; i < 1000; i++ {
		other()
	}
	fn()
}

func long(fn func()) {
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	other()
	fn()
}

func other() int {
	m := map[int]string{}
	m[0] = "foo"
	m[100] = "bar"
	return len(m)
}
