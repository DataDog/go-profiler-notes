⬅ [Block Profiling in Go](./block.md)

⚠️This document describes a sampling bias issue I discovered while researching the block profiler for Go. I have since [landed a fix for it](https://go-review.googlesource.com/c/go/+/299991) that should appear in Go 1.17.

# Block Profiler Sampling Bias

**tl;dr:** Setting your sampling `rate` too high will bias your results towards infrequent long events over frequent short events.

As described in the [Usage](#usage) section, the block profiler will sample as follows:

- Events with `duration >= rate` will be sampled 100%
- Events with `duration < rate` have a `duration / rate` chance of getting sampled.

The [implementation](https://github.com/golang/go/blob/go1.15.7/src/runtime/mprof.go#L408) for this looks like that:

```go
func blocksampled(cycles int64) bool {
	rate := int64(atomic.Load64(&blockprofilerate))
	if rate <= 0 || (rate > cycles && int64(fastrand())%rate > cycles) {
		return false
	}
	return true
}
```

This means that if you set your profiling `rate` low enough, you'll get very accurate results. However, if your `rate` is higher than the `duration` of some of the events you are sampling, the sampling process will exhibit a bias favoring infrequent events of higher `duration` over frequent events with lower `duration` even so they may contribute to the same amount of overall block duration in your program.

## Simple Example

Let's say your `blockprofilerate` is `100ns` and your application produces the following events:

- `A`: `1` event with a duration of `100ns`.
- `B`: `10` events with a duration of `10ns` each.

Given this scenario, the `blockprofiler` is guaranteed to catch and accurately report event `A` as `100ns` in the profile. For event `B`  the most likely outcome is that the profiler will capture only a single event (10% of 10 events) and report `B` as `10ns` in the profile. So you might find yourself in a situation where you think event `A` is causing 10x more blocking than event `B`, which is not true.

## Simulation & Proposal for Improvement

For an even better intuition about this, consider the [simulated example](./sim/block_sampling.ipynb) below. Here we have a histogram of all durations collected from 3 types of blocking events. As you can see, they all have different mean durations (`1000ns`, `2000ns`, `3000ns`) and they are occurring at different frequencies, with `count(a) > count(b) > count(c)`. What's more difficult to see, is that the cumulative durations of these events are the same, i.e. `sum(a) = sum(b) = sum(c)`, but you can trust me on that : ).

<img src="/Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/sim/block_sampling_hist.png" style="zoom: 80%;" />

So given that your application might produce events like this, how will they show up in your block profile as you try out different `blockprofilerate` values? As you can see below, all is well and fine until a `blockprofilerate` of `1000ns`. Each event shows up with the same total duration in the profile (the red and green dots are hidden below the blue ones). However starting at `1000ns` you see that event `a` starts to fade from our profile and at `2000ns` you'd already think that events `b` and `c` are causing twice as much blocking time as event `a`.

<img src="/Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/sim/block_sampling_biased.png" style="zoom:80%;" />

So what can we do? Do we always need to live in fear of bias when working with block profiles? No! If the [Overhead](#overhead) for your workload allows it, the simplest solution is to use a low enough `blockprofilerate` in order to capture most blocking events.

But perhaps there is an even better way. I'm thinking we could correct for the current bias by keeping the same logic of sampling `duration / rate` fraction of events when `duration < rate`. However, when this happens we could simply multiply the sampled duration by `rate/duration` like this:

```
duration = duration * (rate/duration)
# note: the expression above can be simplified to just `duration = rate`
```

Doing so could be done with a [trivial patch](https://github.com/felixge/go/compare/master...debias-blockprofile-rate) to the go runtime and the picture below shows the results from simulating it.

<img src="/Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/sim/block_sampling_debiased.png" alt="" style="zoom: 80%;" />

## Disclaimers

I'm [felixge](https://github.com/felixge) and work at [Datadog](https://www.datadoghq.com/) on [Continuous Profiling](https://www.datadoghq.com/product/code-profiling/) for Go. You should check it out. We're also [hiring](https://www.datadoghq.com/jobs-engineering/#all&all_locations) : ).

The information on this page is believed to be correct, but no warranty is provided. Feedback is welcome!