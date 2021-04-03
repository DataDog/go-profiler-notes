# callers-func-length

A simple benchmark that shows that the costs of stack unwinding in Go increase for larger functions. This is due to the way `gopclntab` based unwinding is implemented.

```
go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/felixge/go-profiler-notes/examples/callers-func-length
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkCallers/short-12         	2623002	      454.1 ns/op
BenchmarkCallers/loop-12          	2590384	      466.8 ns/op
BenchmarkCallers/long-12          	 638096	     1862 ns/op
```
