This document was last updated for `go1.16` but probably still applies to older/newer versions for the most parts.

# Stack Traces in Go

Stack traces play a critical role in Go profiling. So let's try to understand them to see how they might impact the overhead and accuracy of our profiles. 

## Introduction

All Go profilers work by collecting samples of stack traces and putting them into [pprof profiles](./pprof.md). Ignoring some details, a pprof profile is just a frequency table of stack traces like shown below:

| stack trace  | count |
| ------------ | ----- |
| main;foo     | 5     |
| main;foo;bar | 3     |
| main;foobar  | 4     |

Let's zoom in on the first stack trace in the table above: `main;foo`. A Go developer will usually be more familiar with seeing a stack trace like this as rendered by `panic()` or  [`runtime.Stack()`](https://golang.org/pkg/runtime/#Stack) as shown below:

```
goroutine 1 [running]:
main.foo(...)
	/path/to/go-profiler-notes/examples/stack-trace/main.go:9
main.main()
	/path/to/go-profiler-notes/examples/stack-trace/main.go:5 +0x3a
```

This text format has been [described elsewhere](https://www.ardanlabs.com/blog/2015/01/stack-traces-in-go.html) so we won't discuss the details of it here. Instead we'll dive deeper into the source of this data.

## Goroutine Stack

As the name implies, stack traces originate from "the stack". Even so the details vary, most programming languages have a concept of a stack and use it to store things like local variables, arguments, return values and return addresses. Generating a stack trace usually involves navigating the stack in a process known as [Unwinding](#unwinding) that will be described in more detail later on.

Platforms like `x86-64` define a [stack layout](https://eli.thegreenplace.net/2011/09/06/stack-frame-layout-on-x86-64) and calling convention for C and encourage other programming languages to adopt it for interoperability. Go doesn't follow these conventions, and instead uses its own idiosyncratic [calling convention](https://dr-knz.net/go-calling-convention-x86-64.html). Future versions of Go (1.17?) will adopt a more traditional [register-based](https://go.googlesource.com/proposal/+/refs/changes/78/248178/1/design/40724-register-calling.md) convention that will improve performance . However compatibility with platform conventions is not planned as it would negatively impact goroutine scalability.

Even today, Go's stack layout is slightly different on different platforms. To keep things manageable, we'll assume that we're on `x86-64` for the remainder of this note.

### Stack Layout

Now let's take a closer look at the stack. Every goroutine has its own stack that is at least [2 KiB](https://sourcegraph.com/search?q=repo:golang/go+repo:%5Egithub%5C.com/golang/go%24+_StackMin+%3D&patternType=literal) and grows from a high memory address towards lower memory addresses. This can be a bit confusing and is mostly a historical convention from a time when memory was so limited that one had to worry about the stack colliding with other memory regions used by the program.

![](./goroutine-stack.png)

There is a lot going on in the picture above, but for now let's focus on the things highlighted in red. To get a stack trace, the first thing we need is the current program counter (pc) which identifies the function that is currently being executed. This is found in a CPU register called `rip` (instruction pointer register) that points to another region of memory that holds the executable machine code of our program. If you're not familiar with registers, you can think of them as special CPU variables that are incredibly fast to access.

The next step is to find the program counters of all the callers of the current function, i.e. all the `return address (pc)` values that are also highlighted in red. There are various techniques for doing this, which are described in the [Unwinding](#unwinding) section. The end result is a list of program counters that represent a stack trace. In fact, it's exactly the same list you can get from [`runtime.Callers()`](https://golang.org/pkg/runtime/#Callers) within your program. Last but not least, these `pc` values are translated into human readable file/line/function names as described in the [Symbolization](#symbolization) section below.

### Real Example

Looking at pretty pictures can be a good way to get a high level understanding of the stack, but it has its limits. Sometimes you need to look at the raw bits & bytes in order to get a full understanding. If you're not interested in that, feel free to skip ahead to the next section.

To take a look at the stack, we'll use [delve](https://github.com/go-delve/delve) which is a wonderful debugger for Go. In order to inspect the stack, I wrote a script called [stackannotate.star](./delve/stackannotate.star) that can used to print the annotated stack for a simple [example program](./examples/stackannotate/main.go):

```
$ dlv debug ./examples/stackannotate/main.go 
Type 'help' for list of commands.
(dlv) source delve/stackannotate.star
(dlv) continue examples/stackannotate/main.go:19
Breakpoint 1 set at 0x1067d94 for main.bar() ./examples/stackannotate/main.go:19
> main.bar() ./examples/stackannotate/main.go:19 (hits goroutine(1):1 total:1) (PC: 0x1067d94)
    14:	}
    15:	
    16:	func bar(a int, b int) int {
    17:		s := 3
    18:		for i := 0; i < 100; i++ {
=>  19:			s += a * b
    20:		}
    21:		return s
    22:	}
(dlv) stackannotate
regs    addr        offset  value               explanation                     
        c00004c7e8       0                   0  ?                               
        c00004c7e0      -8                   0  ?                               
        c00004c7e8     -16                   0  ?                               
        c00004c7e0     -24                   0  ?                               
        c00004c7d8     -32             1064ac1  return addr to runtime.goexit   
        c00004c7d0     -40                   0  frame pointer for runtime.main  
        c00004c7c8     -48             1082a28  ?                               
        c00004c7c0     -56          c00004c7ae  ?                               
        c00004c7b8     -64          c000000180  var g *runtime.g                
        c00004c7b0     -72                   0  ?                               
        c00004c7a8     -80     100000000000000  var needUnlock bool             
        c00004c7a0     -88                   0  ?                               
        c00004c798     -96          c00001c060  ?                               
        c00004c790    -104                   0  ?                               
        c00004c788    -112          c00001c060  ?                               
        c00004c780    -120             1035683  return addr to runtime.main     
        c00004c778    -128          c00004c7d0  frame pointer for main.main     
        c00004c770    -136          c00001c0b8  ?                               
        c00004c768    -144                   0  var i int                       
        c00004c760    -152                   0  var n int                       
        c00004c758    -160                   0  arg ~r1 int                     
        c00004c750    -168                   1  arg a int                       
        c00004c748    -176             1067c8c  return addr to main.main        
        c00004c740    -184          c00004c778  frame pointer for main.foo      
        c00004c738    -192          c00004c778  ?                               
        c00004c730    -200                   0  arg ~r2 int                     
        c00004c728    -208                   2  arg b int                       
        c00004c720    -216                   1  arg a int                       
        c00004c718    -224             1067d3d  return addr to main.foo         
bp -->  c00004c710    -232          c00004c740  frame pointer for main.bar      
        c00004c708    -240                   0  var i int                       
sp -->  c00004c700    -248                   3  var s int
```

The script isn't perfect and there are some addresses on the stack that it's unable to automatically annotate for now (contributions welcome!). But generally speaking, you should be able to use it to check your understanding against the abstract stack drawing that was presented earlier.

If you want to try it out yourself, perhaps modify the example program to spawn `main.foo()` as a goroutine and observe how that impacts the stack.

## Unwinding

Unwinding (or stack walking) is the process of collecting all the return addresses (see red elements in [Stack Layout](#stack-layout)) from the stack. Together with the current instruction pointer register (`rip`) they form a list of program counter (pc) values that represent a stack trace.

Conceptually unwinding can be seen as a separate step from [Symbolization](#symbolization), however in practice implementations often intertwine the two.

The Go runtime, including the builtin profilers, exclusively use  [.gopclntab](#gopclntab) for unwinding. However, we'll start with describing [Frame Pointer](#frame-pointer) unwinding first, because it's a much easier technique to understand.

### Frame Pointer

Frame pointer unwinding is the simple process of following the base pointer register (`rbp`) to the first frame pointer on the stack which points to the next frame pointer and so on. In other words, it is following the orange lines in the [Stack Layout](#stack-layout) graphic. For each visited frame pointer, the return address (pc) sits 8 bytes above the frame pointer address.

Go only started to include frame pointers on the stack by default since [Go 1.7](https://golang.org/doc/go1.7), and only does so for [64bit binaries](https://sourcegraph.com/search?q=framepointer_enabled+repo:%5Egithub%5C.com/golang/go%24+&patternType=literal). Even so that might cover the platforms most of us care about, it forces the Go runtime and builtin profilers to use a more complicated unwinding approach outlined in [.gopclntab](#gopclntab) that works on all supported platforms.

However, this doesn't mean the frame pointers are useless. 3rd party debuggers and profilers (e.g. [perf](http://www.brendangregg.com/perf.html)) do not understand Go's idiosyncratic unwinding tables, so they rely on frame pointers (and/or [DWARF](#dwarf)) for getting stack traces. 





https://github.com/felixge/gounwind/blob/main/gounwind.go

> Always compile with frame pointers. Omitting frame pointers is an evil compiler optimization that breaks debuggers, and sadly, is often the default.
> – [Brendan Gregg](http://www.brendangregg.com/perf.html)



> An unwinder that is several hundred lines long is simply not even *remotely* interesting to me.
> – [Linus Torvalds](https://lkml.org/lkml/2012/2/10/356)

> ... so as long as the Go compiler emits valid DWARF symbols, then anything is possible. Erm, by "anything" I mean, specificially, "unwinding". Many things are still not possible, even with DWARF (and more frequently than you might think, many things are not possible precisely because DWARF exists-- true love, for instance)
> – [David Sanchez](https://www.linkedin.com/in/david-sanchez-34289130/) (my colleague, who knows more about DWARF than any sane person probably should)



https://github.com/felixge/gounwind





For debuggers and profiling tools (not just in Go), frame pointer unwinding is the most universal, simplest, fastest and most reliable ways to get stack traces utilized 



 Some compilers offer options such as `-fomit-frame-pointers` to d

In order to support 3rd party profiling tools such as [perf](http://www.brendangregg.com/perf.html), Go started to include frame pointers in every stack frame of  since . 





 [.gopclntab](##gopclntab) for 



In fact, [until recently](https://lwn.net/Articles/727553/), the Linux Kernel rejected anything other than frame pointers for stack unwinding despite causing up to 5-10% slowdown for some workloads. Go reported the overhead as 2% in the Go 1.7 release notes.

### .gopclntab

To be written ...

### DWARF

To be written ...

### ORC, LBR, etc.

## Symbolization

To be written ...

## Overhead

To be written ...

- Frame Pointers: 2% for Go, Linux 5-10%

## Accuracy

To be written ...

### Function Inlining

### Frame Pointer Race Condition

To be written ...

### Goroutine Stack Truncation

To be written ...

### cgo

To be written ...

### pprof Labels

To be written ...

## History

In order to support 3rd profilers such as [perf](http://www.brendangregg.com/perf.html) the  [Go 1.7](https://golang.org/doc/go1.7) (2016-08-15) release started to enable frame pointers by default for [64bit binaries](https://sourcegraph.com/search?q=framepointer_enabled+repo:%5Egithub%5C.com/golang/go%24+&patternType=literal).



## Disclaimers

I'm [felixge](https://github.com/felixge) and work at [Datadog](https://www.datadoghq.com/) on [Continuous Profiling](https://www.datadoghq.com/product/code-profiling/) for Go. You should check it out. We're also [hiring](https://www.datadoghq.com/jobs-engineering/#all&all_locations) : ).

The information on this page is believed to be correct, but no warranty is provided. Feedback is welcome!