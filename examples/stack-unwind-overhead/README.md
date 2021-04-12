# stack-unwind-overhead

This directory contains benchmarks to explore which factors impact stack unwinding in Go. It's informed by an analysis of the `gopclntab` unwinding implementation.

## Result 1: O(N) Stack Depth

The benchmark below shows that stack unwinding has O(N) complexity with regard to the number of call frames on the stack:

```
BenchmarkStackDepth/8-12  	1968214	      612.2 ns/op
BenchmarkStackDepth/16-12 	 975457	     1184 ns/op
BenchmarkStackDepth/32-12 	 572706	     2101 ns/op
BenchmarkStackDepth/64-12 	 333598	     3596 ns/op
BenchmarkStackDepth/128-12         	 182450	     6561 ns/op
BenchmarkStackDepth/256-12         	  94783	    12548 ns/op
BenchmarkStackDepth/512-12         	  48439	    24471 ns/op
BenchmarkStackDepth/1024-12        	  24884	    48310 ns/op
```

Tests were done on my local machine:

```
go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/felixge/go-profiler-notes/examples/stack-unwind-overhead
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
```
