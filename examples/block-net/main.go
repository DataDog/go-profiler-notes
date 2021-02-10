package main

import (
	"fmt"
	"net"
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
	runtime.SetBlockProfileRate(1)

	listening := make(chan struct{})
	g := &errgroup.Group{}
	g.Go(func() error {
		return server(listening)
	})

	<-listening
	time.Sleep(time.Second)

	if err := client(); err != nil {
		return err
	}
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

func server(listening chan struct{}) error {
	l, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return err
	}

	listening <- struct{}{}

	conn, err := l.Accept()
	if err != nil {
		return err
	}
	fmt.Printf("[server] accepted: %v\n", conn)
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}
	buf = buf[:n]
	fmt.Printf("[server] read: %s\n", buf)
	return nil
}

func client() error {
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		return err
	}
	defer conn.Close()

	time.Sleep(time.Second)

	fmt.Printf("[client] connected: %v\n", conn)
	_, err = conn.Write([]byte("hello world"))
	return err
}
