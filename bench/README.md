# Go Profiler Overheads

This page is documenting the benchmark methodology used to analyze the performance overhead of the various go profilers. The results are discussed in the documents for each individual profiler.

Benchmarking is done by invoking the Go program included in this directory. You can look at [run.sh](./run.sh) to see the current arguments that are being used, but here is an example for block profiling with two workloads and various profiling rates:

```
go run . \
  -workloads mutex,chan \
  -ops 100000 \
  -blockprofilerates 0,1,10,100,1000,10000,100000,1000000 \
  -runs 20 \
  -depths 16 \
  > "result.csv"
```

The benchmark works by spawning a new child process for the given number of `-runs` and every unique combination of parameters. The child reports the results to the parent process which then combines all the results in a CSV file. The hope is that using a new child process for every config/run eliminates scheduler, GC and other runtime state building up as a source of errors.

Workloads are defined in the [workloads.go](./workloads.go) file. For now the workloads are designed to be **pathological**, i.e. they try to show the worst performance impact the profiler might have on applications that are not doing anything useful other than stressing the profiler. The numbers are not intended to scare you away from profiling in production, but to guide you towards universally **safe profiling rates** as a starting point.

The CSV files are visualized using the [analysis.ipynb](./analysis.ipynb) notebook that's included in this directory.

For now the data is only collected from my local MacBook Pro machine (using docker for mac), but more realistic environments will be included in the future. But it's probably a good setup for finding pathological scenarios : ).

## Disclaimers

I work at [Datadog](https://www.datadoghq.com/) on [Continuous Profiling](https://www.datadoghq.com/product/code-profiling/) for Go (you should check it out) and they generously allowed me to do all this research and publish it.

The information on this page is believed to be correct, but no warranty is provided. Feedback is welcome!