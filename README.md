# go-profiler-notes

Hey there 👋🏻, I'm [felixge](https://github.com/felixge) and I've just started a new job at [Datadog](https://www.datadoghq.com/) to work on [Continuous Profiling](https://www.datadoghq.com/product/code-profiling/) for Go.

I found that Go has a lot of profilers and there are many tools for looking at the data, but that there is very little information on what any of it means. So in order to make sure that I know what I'm talking about, I've started to research the existing profilers and how they work. This repository is my attempt to summarize my findings in the hope that it might be useful to others.

- [pprof tool & format](./pprof.md): Describes the pprof tool and the binary data format for storing profiles.
- [Goroutine Profiling](./goroutine.md): Allows you to get a list of all active goroutines and what they're currently doing.
- [Block Profiling](./block.md): Understand how much time your code spends waiting on channels and locks.
- CPU Profiling (🚧 coming soon!)
- Heap Profiling (🚧 coming soon!)
- Mutex Profiling (🚧 coming soon!)
- Wallclock Profiling  (🚧 coming soon!)

## External Links

- Go Docs
  - [Diagnostics](https://golang.org/doc/diagnostics.html): Has a very good overview over the available profiling and tracing facilities but doesn't go into a lot of depth.
  - [runtime/pprof](https://golang.org/pkg/runtime/pprof/#Profile): Lists the available profiles and has a little more explanation about what kind of data they produce.
  - [runtime](https://golang.org/pkg/runtime/): Has documentation on the various control knobs and pprof facilities, e.g. `MemProfileRate`.
  - [net/http/pprof](https://golang.org/src/net/http/pprof/pprof.go): Not a lot of docs, but diving into the code from there shows how the various profilers can be started/stopped on demand.
- JDB
  - [Profiler labels in Go](https://rakyll.org/profiler-labels/): An introduction to using pprof labels and how they allow you to add additional context to your profiles.
  - [Custom pprof profiles](https://rakyll.org/custom-profiles/): Example for using custom profiles, shows tracking open/close events of a blob store and how to figure out how many blobs are open at a given time.
  - [Mutex profile](https://rakyll.org/mutexprofile/): Brief intro to the mutex profile.
  - [Using Instruments to profile Go programs](https://rakyll.org/instruments/): How to use the macOS Instruments app (I think it's built on dtrace) to profile Go programs. Not clear what the benfits are, if any.
- [Profiling Go programs with pprof](https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/) by Julia Evans: A nice tour with a focus on heap profiling and the pprof output format. 

Got great links to recommend? Open an issue or PR, I'd happy to add your suggestions : ).

## License

The markdown files in this repository are licensed under the [CC BY-SA 4.0 license](https://creativecommons.org/licenses/by-sa/4.0/).

## Disclaimers

I'm [felixge](https://github.com/felixge) and work at [Datadog](https://www.datadoghq.com/) on [Continuous Profiling](https://www.datadoghq.com/product/code-profiling/) for Go. You should check it out. We're also [hiring](https://www.datadoghq.com/jobs-engineering/#all&all_locations) : ).

The information on this page is believed to be correct, but no warranty is provided. Feedback is welcome!
