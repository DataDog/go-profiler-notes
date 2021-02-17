package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	rand.Seed(time.Now().UnixNano())

	workCh := make(chan int)
	go worker(workCh)
	go producer(workCh)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	return nil
}

func producer(workCh chan<- int) {
	for msg := 0; ; msg++ {
		workCh <- msg
		fmt.Printf("producer sent: %d\n", msg)
		if rand.Int63n(10) == 0 {
			takeNap()
		}
	}
}

func worker(workCh <-chan int) {
	for {
		msg := <-workCh
		fmt.Printf("worker received: %d\n", msg)
		if rand.Int63n(10) == 0 {
			takeNap()
		}
	}
}

func takeNap() {
	var forever chan struct{}
	<-forever
}
