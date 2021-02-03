package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if os.Getenv("WORKER") == "" {
		return leader()
	} else {
		return worker()
	}
}

func leader() error {
	var (
		blockprofilerates = flagIntSlice("blockprofilerates", []int{0, 1, 10, 100, 1000, 10000, 100000, 1000000}, "The runtime.SetBlockProfileRate() values to benchmark.")
		workloads         = flagStringSlice("workloads", []string{"mutex", "chan"}, "The workloads to benchmark.")
		ops               = flag.Int("ops", 1000, "The number of operations to perform for each workload.")
		goroutines        = flagIntSlice("goroutines", []int{runtime.NumCPU()}, "The number of goroutine values to use for each workloads.")
		runs              = flag.Int("runs", 3, "The number of times to repeat the same benchmark to understand variance.")
		depths            = flagIntSlice("depths", []int{2, 4, 8, 16, 32}, "The different frame depths values to use for each workload.")
	)
	flag.Parse()

	cw := csv.NewWriter(os.Stdout)
	cw.Write(Headers())
	cw.Flush()

	for _, workload := range *workloads {
		for _, goroutine := range *goroutines {
			for _, blockprofilerate := range *blockprofilerates {
				for _, depth := range *depths {
					for run := 1; run <= *runs; run++ {
						cmd := exec.Command(os.Args[0],
							"-run", fmt.Sprintf("%d", run),
							"-blockprofilerate", fmt.Sprintf("%d", blockprofilerate),
							"-ops", fmt.Sprintf("%d", *ops),
							"-goroutines", fmt.Sprintf("%d", goroutine),
							"-depth", fmt.Sprintf("%d", depth),
							"-workload", workload,
						)

						buf := &bytes.Buffer{}
						cmd.Stdout = buf
						cmd.Stderr = os.Stderr
						cmd.Env = append(cmd.Env, "WORKER=yeah")

						if err := cmd.Run(); err != nil {
							return err
						}

						buf.WriteTo(os.Stdout)
					}
				}
			}
		}
	}

	return nil
}

func worker() error {
	var (
		run              = flag.Int("run", 1, "The number of run. Has no impact on the benchmark, but gets included in the csv output line.")
		blockprofilerate = flag.Int("blockprofilerate", 1, "The block profile rate to use.")
		workload         = flag.String("workload", "mutex", "The workload to simulate.")
		out              = flag.String("blockprofile", "", "Path to a file for writing the block profile.")
		depth            = flag.Int("depth", 16, "The stack depth at which to perform blocking events.")
		ops              = flag.Int("ops", 100000, "The number of operations to perform.")
		goroutines       = flag.Int("goroutines", runtime.NumCPU(), "The number of goroutines to utilize.")
	)
	flag.Parse()

	if *blockprofilerate > 0 {
		runtime.SetBlockProfileRate(*blockprofilerate)
	}

	start := time.Now()
	switch *workload {
	case "mutex":
		if err := mutexWorkload(*goroutines, *ops, *depth); err != nil {
			return err
		}
	case "chan":
		if err := chanWorkload(*goroutines, *ops, *depth); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown workload: %q", *workload)
	}
	duration := time.Since(start)

	if *blockprofilerate > 0 && *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
			return err
		}
	}

	cw := csv.NewWriter(os.Stdout)
	record, err := (&Record{
		Workload:         *workload,
		Ops:              *ops,
		Goroutines:       *goroutines,
		Depth:            *depth,
		Blockprofilerate: *blockprofilerate,
		Run:              *run,
		Duration:         duration,
	}).MarshalRecord()
	if err != nil {
		return err
	}
	cw.Write(record)
	cw.Flush()
	return cw.Error()
}

func atStackDepth(depth int, fn func()) {
	pcs := make([]uintptr, depth*10)
	n := runtime.Callers(1, pcs)
	if n > depth {
		panic("depth exceeded")
	} else if n < depth {
		atStackDepth(depth, fn)
		return
	}

	fn()
}

func flagIntSlice(name string, value []int, usage string) *[]int {
	val := &intSlice{vals: value}
	flag.Var(val, name, usage)
	return &val.vals
}

type intSlice struct {
	vals []int
}

func (i *intSlice) Set(val string) error {
	var vals []int
	for _, val := range strings.Split(val, ",") {
		num, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		vals = append(vals, num)
	}
	i.vals = vals
	return nil
}

func (i *intSlice) String() string {
	return fmt.Sprintf("%v", i.vals)
}

func flagStringSlice(name string, value []string, usage string) *[]string {
	val := &strSlice{vals: value}
	flag.Var(val, name, usage)
	return &val.vals
}

type strSlice struct {
	vals []string
}

func (s *strSlice) Set(val string) error {
	s.vals = strings.Split(val, ",")
	return nil
}

func (s *strSlice) String() string {
	return fmt.Sprintf("%v", s.vals)
}
