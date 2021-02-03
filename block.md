This document was last updated for `go1.15.7` but probably still applies to older/newer versions for the most parts.

# Block Profiling in Go

## Description

The block profile in Go lets you analyze how much time your program spends waiting on the blocking operations listed below:

- [select](https://github.com/golang/go/blob/go1.15.7/src/runtime/select.go#L511)
- [chan send](https://github.com/golang/go/blob/go1.15.7/src/runtime/chan.go#L279) (except sending to a nil chan which blocks forever but is [not tracked](https://github.com/golang/go/blob/go1.15.7/src/runtime/chan.go#L163))
- [chan receive](https://github.com/golang/go/blob/go1.15.7/src/runtime/chan.go#L586) (except receiving from a nil chan which blocks forever but is [not tracked](https://github.com/golang/go/blob/go1.15.7/src/runtime/chan.go#L466))
- [semacquire](https://github.com/golang/go/blob/go1.15.7/src/runtime/sema.go#L150) ( [`Mutex.Lock`](https://golang.org/pkg/sync/#Mutex.Lock), [`RWMutex.RLock`](https://golang.org/pkg/sync/#RWMutex.RLock) , [`RWMutex.Lock`](https://golang.org/pkg/sync/#RWMutex.Lock), [`WaitGroup.Wait`](https://golang.org/pkg/sync/#WaitGroup.Wait))
- [notifyListWait](https://github.com/golang/go/blob/go1.15.7/src/runtime/sema.go#L515) ( [`Cond.Wait`](https://golang.org/pkg/sync/#Cond.Wait))

Time is only tracked when Go has to suspend the goroutine's execution by parking it into a [waiting](https://github.com/golang/go/blob/go1.15.7/src/runtime/runtime2.go#L51-L59) state. So for example a `Mutex.Lock()` operation will not show up in your profile if the lock can be immediately aquired.

The operations above are a subset of the [waiting states](https://github.com/golang/go/blob/go1.15.7/src/runtime/runtime2.go#L996-L1024) used by the Go runtime, i.e. the operations below **will not** show up in a block profile:

- [`time.Sleep`](https://golang.org/pkg/time/#Sleep) (but [`time.After`](https://golang.org/pkg/time/#After), [`time.Tick`](https://golang.org/pkg/time/#Tick) and other channel based wrappers will show up)
- GC
- Syscalls
- Internal Locks (e.g. for [stopTheWorld](https://github.com/golang/go/blob/go1.15.7/src/runtime/proc.go#L900))
- Blocking in [cgo](https://golang.org/cmd/cgo/) calls

## Usage

The block profiler is disabled by default. To record all blocking events regardless of their duration, simply call:

```
runtime.SetBlockProfileRate(1)
```

You should read the [Accuracy](#accuracy) and [Overhead](#overhead) sections below to figure out if passing `1` as the sampling rate is a good idea for you or not.

Block durations are aggregated over the lifetime of the program (while the profiling is enabled). To get a [pprof formated](./pprof.md) snapshot of the current stack traces that lead to blocking events and their cumulative time duration, you can call:

```go
pprof.Lookup("block").WriteTo(myFile, 0)
```

Alternatively you can use [github.com/pkg/profile](https://pkg.go.dev/github.com/pkg/profile) for convenience, or [net/http/pprof](net/http/pprof) to expose profiling via http, or use a [continious profiler](https://www.datadoghq.com/product/code-profiling/) to collect the data automatically in production.

Last but not least you can use the [`runtime.BlockProfile`](https://golang.org/pkg/runtime/#BlockProfile) API to get programmatically get the same information.

## Overhead

**tl;dr:** A blockprofile rate `>= 10000` (10µs) should have negligable impact on production apps, including those suffering from extreme contention.

### Implementation Details

Block profiling is essentially implemented like this inside of the Go runtime (see the links in the [Description](#description) above for real code):

```go
func chansend(...) {
  var t0 int64
  if blockprofilerate > 0 {
    t0 = cputicks()
  }
  // ... park goroutine in waiting state while blocked ...
  if blockprofilerate > 0 {
    cycles := cputicks() - t0
    if blocksampled(cycles) {
      saveblockevent(cycles)
    }
  }
}
```

This means that unless you enable block profiling, the overhead should be effectively zero thanks to CPU branch prediction.

When block profiling is enabled, every blocking operation will pay the overhead of two `cputicks()` calls. On `amd64` this is done via [optimized assembly](https://github.com/golang/go/blob/go1.15.7/src/runtime/asm_amd64.s#L874-L887) using the [RDTSC instruction](https://en.wikipedia.org/wiki/Time_Stamp_Counter) and takes a negligible `~10ns/op` on [my machine](https://github.com/felixge/dump/tree/master/cputicks). On other platforms various alternative clock sources are used which may have higher overheads and lower accuracy.

Depending on the configured `blockprofilerate` (more about this in the [Accuracy](#accuracy) section) the block event may end up getting saved. This means a stack trace is collected which takes `~1µs` on [my machine](https://github.com/felixge/dump/tree/master/go-callers-bench) (stackdepth=16). The stack is then used as a key to update an [internal hashmap](https://github.com/golang/go/blob/go1.15.7/src/runtime/mprof.go#L144) by incrementing the corresponding [`blockRecord`](https://github.com/golang/go/blob/go1.15.7/src/runtime/mprof.go#L133-L138) count and cycles.

```go
type blockRecord struct {
	count  int64
	cycles int64
}
```

The costs of updating the hash map is probably similar to collecting the stack traces, but I haven't measured it yet.

### Benchmarks

Anyway, what does all of this mean in terms of overhead for your application? It means that block profiling is **low overhead**. Unless your application spends literally all of its time parking/unparking goroutines due to contention, you probably won't be able to see a measurable impact even when sampling every block event.

That being said, the benchmark results below (see [Methodology](./bench/)) should give you an idea of the **theoretical worst case** overhead block profiling could have. The graph `chan(cap=0)` shows that setting `blockprofilerate` from  `1` to `1000` on a workload that consists entirely in sending tiny messages across unbuffered channels decreases throughput significantly. Using a buffered channel as in graph `chan(cap=128)` greatly reduces the problem to the point that it probably won't matter for real applications that don't spend all of their time on channel communication overheads.

It's also interesting to note that I was unable to see significant overheads for `mutex` based workloads. I believe this is due to the fact that mutexes employe spin locks before parking a goroutine when there is contention. If somebody has a good idea for a workload that exhibits high non-spinning mutex contention in Go, please let me know!

![block_linux_x86_64](./bench/block_linux_x86_64.png)

### Initialization Costs 

The first call to `runtime.SetBlockProfileRate()` takes `100ms` because it tries to [measure](https://github.com/golang/go/blob/go1.15.7/src/runtime/runtime.go#L22-L47) the speed difference between the wall clock and the [TSC](https://en.wikipedia.org/wiki/Time_Stamp_Counter) clock. However, recent changes around async preemption have [broken](https://github.com/golang/go/issues/40653#issuecomment-766340860) this code, so the call is taking only `~10ms` right now.

### Memory Usage

Block profiling utilizes a shared hash map that takes up 1.4 MiB even when empty. Unless you explicitly [disable heap profiling](https://twitter.com/felixge/status/1355846360562589696) in your application, this map will get allocated regardless of whether you use the block profiler or not.

Addtionally each unique stack trace will take up some additional memory. The `BuckHashSys` field of [`runtime.MemStats`](https://golang.org/pkg/runtime/#MemStats) allows you to inspect this usage at runtime. In the future I might try to provide additional information about this along with real world data.

## Accuracy

🚧

- `rate` concept [same as memory profiler](https://codereview.appspot.com/6443115#msg38)? 
- sampling (should cycles < rate events be multiplied?)
- multi socket cpus / tsc

- max stack depth 32
- spin locks
- [runtime: problems with rdtsc in VMs moving backward · Issue #8976 · golang/go](https://github.com/golang/go/issues/8976)
- https://github.com/golang/go/issues/16755#issuecomment-332279965 ?

## Relationship with Time

🚧

## Relationship with Mutex Profiling

🚧

## pprof Output

🚧

## pprof Labels

🚧

## History

Block profiling was [implemented](https://codereview.appspot.com/6443115) by [Dmitry Vyukov](https://github.com/dvyukov) and first appeared in the [go1.1](https://golang.org/doc/go1.1) release (2013-05-13).

## Disclaimers

I work at [Datadog](https://www.datadoghq.com/) on [Continuous Profiling](https://www.datadoghq.com/product/code-profiling/) for Go (you should check it out) and they generously allowed me to do all this research and publish it.

The information on this page is believed to be correct, but no warranty is provided. Feedback is welcome!