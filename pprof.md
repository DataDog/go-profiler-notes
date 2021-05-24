⬅ [Index of all go-profiler-notes](./README.md)

# Go's pprof tool & format

The various profilers built into Go are designed to work with the pprof visualization tool. [pprof](https://github.com/google/pprof) itself is an inofficial Google project that is designed to analyze profiling data from C++, Java and Go programs. The project defines a protocol buffer format that is used by all Go profilers and described in this document.

The Go project itself [bundles](https://github.com/golang/go/tree/master/src/cmd/pprof) a version of pprof that can be invoked via  `go tool pprof`. It's largely the same as the upstream tool, except for a few tweaks. Go recommends to always use `go tool pprof` instead of the upstream tool for working with Go profiles.

## Features

The pprof tool features an interactive command line interface, but also a Web UI, as well as various other output format options.

## File Format

### Description

The pprof format is defined in the [profile.proto](https://github.com/google/pprof/blob/master/proto/profile.proto) protocol buffer definition which has good comments. Additionally there is an official [README](https://github.com/google/pprof/blob/master/proto/README.md) for it. pprof files are always stored using gzip-compression on disk.

A picture is worth a thousand words, so below is an automatically [generated](https://github.com/seamia/protodot) visualization of the format. Please note that fields such as `filename` are pointers into the `string_table` which are not visualized, improvements for this would be welcome!

<!-- transparent to white bg: convert -flatten profile.png profile.png -->

<img src="./profile.png" alt="profile.proto visualized" style="zoom: 80%;" />

pprof's data format appears to be designed to for efficency, multiple languages and different profile types (CPU, Heap, etc.), but because of this it's very abstract and full of indirection. If you want all the details, follow the links above. If you want the **tl;dr**, keep reading:

A pprof file contains a list of **stack traces** called *samples* that have one or more numeric **value** associated with them. For a CPU profile the value might be the CPU time duration in nanoseonds that the stack trace was observed for during profiling. For a heap profile it might be the number of bytes allocated. The **value types** themselves are described in the beginning of the file and used to populate the "SAMPLE" drop down in the pprof UI. In addition to the values, each stack trace can also include a set of **labels**. The labels are key-value pairs and can even include a unit. In Go those labels are used for [profiler labels](https://rakyll.org/profiler-labels/).

The profile also includes the **time** (in UTC) that the profile was recorded, and the **duration** of the recording.

Additionally the format allows for **drop/keep** regexes for excluding/including certain stack traces, but they're [not used](https://github.com/golang/go/blob/go1.15.6/src/runtime/pprof/proto.go#L375-L376) by Go. There is also room for a list of **comments** ([not used](https://github.com/golang/go/search?q=tagProfile_Comment) either), as well as describing the **periodic** interval at which samples were taken.

The code for generating pprof output in Go can be found in: [runtime/pprof/proto.go](https://github.com/golang/go/blob/go1.15.6/src/runtime/pprof/proto.go).

### Decoding

Below are a few tools for decoding pprof files into human readable text output. They are ordered by the complexity of their output format, with tools showing simplified output being listed first:

#### Using `pprofutils`

[pprofutils](https://github.com/felixge/pprofutils) is a small tool for converting between pprof files and Brendan Gregg's [folded text](https://github.com/brendangregg/FlameGraph#2-fold-stacks) format. You can use it like this:

```
$ pprof2text < examples/cpu/pprof.samples.cpu.001.pb.gz

golang.org/x/sync/errgroup.(*Group).Go.func1;main.run.func2;main.computeSum 19
golang.org/x/sync/errgroup.(*Group).Go.func1;main.run.func2;main.computeSum;runtime.asyncPreempt 5
runtime.mcall;runtime.gopreempt_m;runtime.goschedImpl;runtime.schedule;runtime.findrunnable;runtime.stopm;runtime.notesleep;runtime.semasleep;runtime.pthread_cond_wait 1
runtime.mcall;runtime.park_m;runtime.schedule;runtime.findrunnable;runtime.checkTimers;runtime.nanotime;runtime.nanotime1 1
runtime.mcall;runtime.park_m;runtime.schedule;runtime.findrunnable;runtime.stopm;runtime.notesleep;runtime.semasleep;runtime.pthread_cond_wait 2
runtime.mcall;runtime.park_m;runtime.resetForSleep;runtime.resettimer;runtime.modtimer;runtime.wakeNetPoller;runtime.netpollBreak;runtime.write;runtime.write1 7
runtime.mstart;runtime.mstart1;runtime.sysmon;runtime.usleep 3
```

#### Using `go tool pprof`

`pprof` itself has an output mode called `-raw` that will show you the contents of a pprof file. However, it should be noted that this is not as `-raw` as it gets, checkout `protoc` below for that:

```
$ go tool pprof -raw examples/cpu/pprof.samples.cpu.001.pb.gz

PeriodType: cpu nanoseconds
Period: 10000000
Time: 2021-01-08 17:10:32.116825 +0100 CET
Duration: 3.13
Samples:
samples/count cpu/nanoseconds
         19  190000000: 1 2 3
          5   50000000: 4 5 2 3
          1   10000000: 6 7 8 9 10 11 12 13 14
          1   10000000: 15 16 17 11 18 14
          2   20000000: 6 7 8 9 10 11 18 14
          7   70000000: 19 20 21 22 23 24 14
          3   30000000: 25 26 27 28
Locations
     1: 0x1372f7f M=1 main.computeSum /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/cpu/main.go:39 s=0
     2: 0x13730f2 M=1 main.run.func2 /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/cpu/main.go:31 s=0
     3: 0x1372cf8 M=1 golang.org/x/sync/errgroup.(*Group).Go.func1 /Users/felix.geisendoerfer/go/pkg/mod/golang.org/x/sync@v0.0.0-20201207232520-09787c993a3a/errgroup/errgroup.go:57 s=0
     ...
Mappings
1: 0x0/0x0/0x0   [FN]
```

The output above is truncated, [examples/cpu/pprof.samples.cpu.001.pprof.txt](./examples/cpu/pprof.samples.cpu.001.pprof.txt) has the full version.

#### Using `protoc`

For those interested in seeing data closer to the raw binary storage, we need the `protoc` protocol buffer compiler. On macOS you can use `brew install protobuf` to install it, for other platform take a look at the [README's install section](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation).

Now let's take a look at the same CPU profile from above:

```
$ gzcat examples/cpu/pprof.samples.cpu.001.pb.gz | protoc --decode perftools.profiles.Profile ./profile.proto

sample_type {
  type: 1
  unit: 2
}
sample_type {
  type: 3
  unit: 4
}
sample {
  location_id: 1
  location_id: 2
  location_id: 3
  value: 19
  value: 190000000
}
sample {
  location_id: 4
  location_id: 5
  location_id: 2
  location_id: 3
  value: 5
  value: 50000000
}
...
mapping {
  id: 1
  has_functions: true
}
location {
  id: 1
  mapping_id: 1
  address: 20393855
  line {
    function_id: 1
    line: 39
  }
}
location {
  id: 2
  mapping_id: 1
  address: 20394226
  line {
    function_id: 2
    line: 31
  }
}
...
function {
  id: 1
  name: 5
  system_name: 5
  filename: 6
}
function {
  id: 2
  name: 7
  system_name: 7
  filename: 6
}
...
string_table: ""
string_table: "samples"
string_table: "count"
string_table: "cpu"
string_table: "nanoseconds"
string_table: "main.computeSum"
string_table: "/Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/cpu/main.go"
...
time_nanos: 1610122232116825000
duration_nanos: 3135113726
period_type {
  type: 3
  unit: 4
}
period: 10000000
```

The output above is truncated also, [pprof.samples.cpu.001.protoc.txt](./examples/cpu/pprof.samples.cpu.001.protoc.txt) has the full version.

## History

The [original pprof tool](https://github.com/gperftools/gperftools/blob/master/src/pprof) was a perl script developed internally at Google. Based on the copyright header, development might go back to 1998. It was first released in 2005 as part of [gperftools](https://github.com/google/tcmalloc/blob/master/docs/gperftools.md), and [added](https://github.com/golang/go/commit/c72fb37425f6ee6297371e0053d6d1f958d49a41) to the Go project in 2010.

In 2014 the Go project [replaced](https://github.com/golang/go/commit/8b5221a57b41a19abcb4e3dde20014af494048c2) the perl based version of the pprof tool with a Go implementation by [Raul Silvera](https://www.linkedin.com/in/raul-silvera-a0521b55/) that was already used inside of Google at this point. This implementation was re-released as a [standalone project](https://github.com/google/pprof) in 2016. Since then the Go project has been vendoring a copy of the upstream project, [updating](https://github.com/golang/go/commits/master/src/cmd/vendor/github.com/google/pprof) it on a regular basis.

Go 1.9 (2017-08-24) added support for pprof labels. It also started including symbol information into pprof files by default, which allows viewing profiles without having access to the binary.

## Todo

- Write more about using `go tool pprof` itself.
- Explain why pprof can be given a path to the binary the profile belongs to.
- Get into more details about line numbers / addresses.
- Talk about mappings and when a Go binary might have more than one

## Disclaimers

I'm [felixge](https://github.com/felixge) and work at [Datadog](https://www.datadoghq.com/) on [Continuous Profiling](https://www.datadoghq.com/product/code-profiling/) for Go. You should check it out. We're also [hiring](https://www.datadoghq.com/jobs-engineering/#all&all_locations) : ).

The information on this page is believed to be correct, but no warranty is provided. Feedback is welcome!
