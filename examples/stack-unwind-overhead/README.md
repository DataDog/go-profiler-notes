# stack-unwind-overhead

This directory contains benchmarks to explore which factors impact stack unwinding in Go. It's informed by an analysis of the `gopclntab` unwinding implementation.

## Result 1: O(N) Stack Depth

The benchmark below shows that stack unwinding has O(N) complexity with regard to the number of call frames on the stack:

```
BenchmarkStackDepth/8-12    	2208600	      547.2 ns/op
BenchmarkStackDepth/16-12   	1447922	      810.8 ns/op
BenchmarkStackDepth/32-12   	 915291	     1338 ns/op
BenchmarkStackDepth/64-12   	 488719	     2366 ns/op
BenchmarkStackDepth/128-12  	 264735	     4462 ns/op
BenchmarkStackDepth/256-12  	 137575	     8643 ns/op
BenchmarkStackDepth/512-12  	  68355	    17316 ns/op
BenchmarkStackDepth/1024-12 	  34710	    34810 ns/op
```

## Result 2: O(N) Function Size

Perhaps suprisingly, stack unwinding is also O(N) with regard to the size of the generated machine code for the function:

```
BenchmarkFunctionSize/1-12  	2562176	      462.8 ns/op
BenchmarkFunctionSize/2-12  	2509465	      484.7 ns/op
BenchmarkFunctionSize/4-12  	2356609	      504.6 ns/op
BenchmarkFunctionSize/8-12  	2095870	      568.3 ns/op
BenchmarkFunctionSize/16-12 	1778889	      669.7 ns/op
BenchmarkFunctionSize/32-12 	1396009	      856.0 ns/op
BenchmarkFunctionSize/64-12 	 943807	     1269 ns/op
BenchmarkFunctionSize/128-12         	 516487	     2271 ns/op
BenchmarkFunctionSize/256-12         	 277821	     4490 ns/op
```

## Disclaimer

YMMV, and especially the function size also depends on the program counter at which the function is being unwound. All tests were done on my local machine:

```
go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/felixge/go-profiler-notes/examples/stack-unwind-overhead
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
```
