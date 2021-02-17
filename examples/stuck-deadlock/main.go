package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	var (
		a = &sync.Mutex{}
		b = &sync.Mutex{}
	)
	go bob(a, b)
	go alice(a, b)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	return nil
}

func bob(a, b *sync.Mutex) {
	for {
		fmt.Println("bob is okay")
		a.Lock()
		b.Lock()
		a.Unlock()
		b.Unlock()
	}
}

func alice(a, b *sync.Mutex) {
	for {
		fmt.Println("alice is okay")
		b.Lock()
		a.Lock()
		b.Unlock()
		a.Unlock()
	}
}
