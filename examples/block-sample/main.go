package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	labels := pprof.Labels("test_label", "test_value")
	ctx := pprof.WithLabels(context.Background(), labels)
	pprof.SetGoroutineLabels(ctx)

	runtime.SetBlockProfileRate(int((40 * time.Microsecond).Nanoseconds()))
	done := make(chan struct{})
	g := errgroup.Group{}
	g.Go(func() error {
		return eventA(done)
	})
	g.Go(func() error {
		return eventB(done)
	})
	time.Sleep(time.Second)
	close(done)
	if err := g.Wait(); err != nil {
		return err
	}

	f, err := os.Create("block.pb.gz")
	if err != nil {
		return err
	}
	defer f.Close()
	if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
		return err
	}
	return nil
}

func eventA(done chan struct{}) error {
	return simulateBlockEvents(20*time.Microsecond, done)
}

func eventB(done chan struct{}) error {
	return simulateBlockEvents(40*time.Microsecond, done)
}

const tolerance = 1.1

func simulateBlockEvents(meanDuration time.Duration, done chan struct{}) error {
	var (
		prev   time.Time
		sum    time.Duration
		count  int
		ticker = time.NewTicker(meanDuration)
	)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			if !prev.IsZero() {
				sum += now.Sub(prev)
				count += 1
				if count > 1000 {
					actualMean := float64(sum) / float64(count)
					max := tolerance * float64(meanDuration)
					min := float64(meanDuration) / tolerance
					if actualMean <= min || actualMean >= max {
						return fmt.Errorf("low clock accuracy: got=%s want=%s", time.Duration(actualMean), meanDuration)
					}
				}
			}
			prev = now
		case <-done:
			return nil
		}
	}
}

/*

Bias in current go version:

	$ go version
	go version go1.15.6 darwin/amd64
	$ go run . && go tool pprof -raw block.pb.gz
	PeriodType: contentions count
	Period: 1
	Time: 2021-02-05 15:27:01.371414 +0100 CET
	Samples:
	contentions/count delay/nanoseconds
				23271  892063188: 1 2 3 4
				22612  438270491: 1 2 5 4
	Locations
			 1: 0x10453af M=1 runtime.selectgo /usr/local/Cellar/go/1.15.6/libexec/src/runtime/select.go:511 s=0
			 2: 0x10cf56b M=1 main.simulateBlockEvents /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:66 s=0
			 3: 0x10cf8b2 M=1 main.eventB /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:52 s=0
							 main.run.func2 /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:28 s=0
			 4: 0x10cefd8 M=1 golang.org/x/sync/errgroup.(*Group).Go.func1 /Users/felix.geisendoerfer/go/pkg/mod/golang.org/x/sync@v0.0.0-20201207232520-09787c993a3a/errgroup/errgroup.go:57 s=0
			 5: 0x10cf852 M=1 main.eventA /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:48 s=0
							 main.run.func1 /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:25 s=0
	Mappings
	1: 0x0/0x0/0x0   [FN]

After removing bias with this patch:
https://github.com/felixge/go/commit/7c3f70c378d3f9a331a2079d17cd4cc420c70190

	$ sgo version
	go version devel +7c3f70c378 2021-02-05 15:19:48 +0100 darwin/amd64
	$ sgo run . && go tool pprof -raw block.pb.gz
	PeriodType: contentions count
	Period: 1
	Time: 2021-02-05 15:18:56.401244 +0100 CET
	Samples:
	contentions/count delay/nanoseconds
				22833  931500628: 1 2 3 4
				22453  902100036: 1 2 5 4
						1      39999: 6 7 8 9 10
	Locations
			 1: 0x10471d9 M=1 runtime.selectgo /Users/felix.geisendoerfer/go/src/github.com/golang/go/src/runtime/select.go:492 s=0
			 2: 0x10d023e M=1 main.simulateBlockEvents /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:66 s=0
			 3: 0x10d0592 M=1 main.eventB /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:52 s=0
							 main.run.func2 /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:28 s=0
			 4: 0x10cfcd8 M=1 golang.org/x/sync/errgroup.(*Group).Go.func1 /Users/felix.geisendoerfer/go/pkg/mod/golang.org/x/sync@v0.0.0-20201207232520-09787c993a3a/errgroup/errgroup.go:57 s=0
			 5: 0x10d0532 M=1 main.eventA /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:48 s=0
							 main.run.func1 /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:25 s=0
			 6: 0x107b024 M=1 sync.(*WaitGroup).Wait /Users/felix.geisendoerfer/go/src/github.com/golang/go/src/sync
	/waitgroup.go:130 s=0
			 7: 0x10cfb30 M=1 golang.org/x/sync/errgroup.(*Group).Wait /Users/felix.geisendoerfer/go/pkg/mod/golang.org/x/sync@v0.0.0-20201207232520-09787c993a3a/errgroup/errgroup.go:40 s=0
			 8: 0x10cff7a M=1 main.run /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:32 s=0
			 9: 0x10cfda5 M=1 main.main /Users/felix.geisendoerfer/go/src/github.com/felixge/go-profiler-notes/examples/block-sample/main.go:14 s=0
			10: 0x10371b5 M=1 runtime.main /Users/felix.geisendoerfer/go/src/github.com/golang/go/src/runtime/proc.go:225 s=0
	Mappings
	1: 0x0/0x0/0x0   [FN]

*/
