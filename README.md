# go-profiler-notes

I've just started a new job at [Datadog](https://www.datadoghq.com/) to work on [Continuous Profiling](https://www.datadoghq.com/product/code-profiling/) for Go. To make sure that I know what I'm talking about, I'm planning to do an in-depth study of the existing profilers and how they work. I'll try to summarize what I learned in this repository as it might be useful to others.

- [pprof tool & format](./pprof.md): Describes the pprof tool and it's binary data format.

## Todo

- CPU Profiler
- Heap Profiler
- Mutex Profiler
- Goroutine Profiler
- Goroutine Profiler
- Wallclock Profiler (fgprof)

## External Links

- Go Docs
  - [Diagnostics](https://golang.org/doc/diagnostics.html): Has a very good overview over the available profiling and tracing facilities but doesn't go into a lot of depth.
  - [runtime/pprof](https://golang.org/pkg/runtime/pprof/#Profile): Lists the available profiles and has a little more explanation about what kind of data they produce.
  - [runtime](https://golang.org/pkg/runtime/): Has documentation on the various control knobs and pprof facilities, e.g. `MemProfileRate`.
  - [net/http/pprof](net/http/pprof): Not a lot of docs, but diving into the code from there shows how the various profilers can be started/stopped on demand.
- JDB
  - [Profiler labels in Go](https://rakyll.org/profiler-labels/): An introduction to using pprof labels and how they allow you to add additional context to your profiles.
  - [Custom pprof profiles](https://rakyll.org/custom-profiles/): Example for using custom profiles, shows tracking open/close events of a blob store and how to figure out how many blobs are open at a given time.
  - [Mutex profile](https://rakyll.org/mutexprofile/): Brief intro to the mutex profile.
  - [Using Instruments to profile Go programs](https://rakyll.org/instruments/): How to use the macOS Instruments app (I think it's built on dtrace) to profile Go programs. Not clear what the benfits are, if any.
- [Profiling Go programs with pprof](https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/) by Julia Evans: A nice tour with a focus on heap profiling and the pprof output format. 

Got great links to recommend? Open an issue or PR, I'd happy to add your suggestions : ).