package main

import (
	"testing"
)

func Benchmark_compute(b *testing.B) {
	compute(b.N)
}
