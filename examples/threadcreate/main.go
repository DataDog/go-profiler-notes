package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	runtime.GOMAXPROCS(1)
	n := 16
	ch := make(chan int)
	fmt.Printf("start work\n")
	for i := 0; i < n; i++ {
		go func() {
			buf := make([]byte, 1024)
			fmt.Printf("syscall\n")
			n, err := syscall.Read(1, buf)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%d\n", n)
		}()
		go dowork(ch)
	}
	runtime.GOMAXPROCS(8)
	time.Sleep(50 * time.Millisecond)
	runtime.GC()
	if err := writeProfile("threadcreate"); err != nil {
		return err
	}
	fmt.Printf("profile\n")
	for i := 0; i < n; i++ {
		fmt.Printf("%d\n", <-ch)
	}
	fmt.Printf("done\n")

	return nil
}

func dowork(done chan int) {
	s := 0
	for i := 0; i < 1000000000; i++ {
		if i%2 == 0 {
			s++
		}
	}
	done <- s
}

func writeProfile(name string) error {
	f, err := os.Create(name + ".pb.gz")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := pprof.Lookup(name).WriteTo(f, 0); err != nil {
		return err
	}
	return nil
}
