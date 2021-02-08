package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	sleep := 1*time.Minute + 10*time.Second

	go sleepLoop(sleep)
	go sleepLoop(time.Second)
	go chanReceiveForever()
	go indirectSleepLoop(time.Second)

	time.Sleep(time.Second)
	// g.wait
	runtime.GC()
	fmt.Printf("sleeping for %s before showing stack traces\n", sleep)
	time.Sleep(sleep)
	runtime.GC()

	fmt.Printf("%s has passed, dumping all stacks", sleep)
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, true)
	buf = buf[:n]
	fmt.Printf("%s\n", buf)

	fmt.Printf("waiting indefinitely so you can press ctrl+\\ to compare the output\n")
	runtime.GC()
	chanReceiveForever()
}

func sleepLoop(d time.Duration) {
	for {
		time.Sleep(d)
	}
}

func chanReceiveForever() {
	forever := make(chan struct{})
	<-forever
}

func indirectSleepLoop(d time.Duration) {
	indirectSleepLoop2(d)
}

func indirectSleepLoop2(d time.Duration) {
	go sleepLoop(d)
}
