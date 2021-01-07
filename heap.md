# Go's Heap Profiler

Totally a work in progress ... !!! Please come back later : ).

## How does it work?

Look at the go source code to understand how the data is captured, what role runtime.MemProfileRate plays, etc.

## How does it fail?

Think through sampling rate issues, alloc size classes, etc.

## Data Format

Figure out how the data ends up in the pprof file.

## GC Control

```
# turn of gc
GOGC=off go run <code>
# print gc events to stdout
GODEBUG=gctrace=1 go run <code>
```

## Questions

- What are the [docs](https://golang.org/pkg/runtime/pprof/#Profile) talking about here? How do I actually use this?

  > Pprof's -inuse_space, -inuse_objects, -alloc_space, and -alloc_objects flags select which to display, defaulting to -inuse_space (live objects, scaled by size).

  A: Those flags are deprecated. Easiest way to select this stuff is via the pprof web ui's sample drop down.

- The [docs](https://golang.org/pkg/runtime/pprof/#Profile) say I should get some kind of data, even if there is no GC. I can reproduce that, but the data seems to not change?

  > If there has been no garbage collection at all, the heap profile reports all known allocations. This exception helps mainly in programs running without garbage collection enabled, usually for debugging purposes.