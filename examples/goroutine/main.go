package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/pprof"
	"time"
)

type Profile struct {
	Name    string
	WriteTo func(w io.Writer) error
}

var profiles = []Profile{
	{
		Name: "runtime.stack.txt",
		WriteTo: func(w io.Writer) error {
			buf := make([]byte, 1024*1024)
			n := runtime.Stack(buf, true)
			buf = buf[:n]
			_, err := w.Write(buf)
			return err
		},
	},
	{
		Name: "runtime.goroutineprofile.json",
		WriteTo: func(w io.Writer) error {
			p := make([]runtime.StackRecord, 1000)
			n, ok := runtime.GoroutineProfile(p)
			if !ok {
				return errors.New("runtime.GoroutineProfile: not ok")
			}
			p = p[0:n]
			e := json.NewEncoder(w)
			e.SetIndent("", "  ")
			return e.Encode(p)
		},
	},
	{
		Name: "pprof.lookup.goroutine.debug0.pb.gz",
		WriteTo: func(w io.Writer) error {
			profile := pprof.Lookup("goroutine")
			return profile.WriteTo(w, 0)
		},
	},
	{
		Name: "pprof.lookup.goroutine.debug1.txt",
		WriteTo: func(w io.Writer) error {
			profile := pprof.Lookup("goroutine")
			return profile.WriteTo(w, 1)
		},
	},
	{
		Name: "pprof.lookup.goroutine.debug2.txt",
		WriteTo: func(w io.Writer) error {
			profile := pprof.Lookup("goroutine")
			return profile.WriteTo(w, 2)
		},
	},
	{
		Name: "net.http.pprof.goroutine.debug0.pb.gz",
		WriteTo: func(w io.Writer) error {
			return writeHttpProfile(w, 0)
		},
	},
	{
		Name: "net.http.pprof.goroutine.debug1.txt",
		WriteTo: func(w io.Writer) error {
			return writeHttpProfile(w, 1)
		},
	},
	{
		Name: "net.http.pprof.goroutine.debug2.txt",
		WriteTo: func(w io.Writer) error {
			return writeHttpProfile(w, 2)
		},
	},
}

func writeHttpProfile(w io.Writer, debug int) error {
	url := fmt.Sprintf("http://%s/debug/pprof/goroutine?debug=%d", listenAddr, debug)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(w, res.Body)
	return err
}

func writeProfiles(n int) error {
	for _, profile := range profiles {
		buf := &bytes.Buffer{}
		filename := fmt.Sprintf("%d.%s", n, profile.Name)
		if err := profile.WriteTo(buf); err != nil {
			return err
		} else if err := ioutil.WriteFile(filename, buf.Bytes(), 0666); err != nil {
			return err
		}
	}
	return nil
}

var listenAddr = "127.0.0.1:8080"

func main() {
	flag.Parse()

	errCh := make(chan error, 1)
	go func() {
		fmt.Printf("Listening for pprof requests on %s\n", listenAddr)
		errCh <- http.ListenAndServe(listenAddr, nil)
	}()

	labels := pprof.Labels("test_label", "test_value")
	ctx := pprof.WithLabels(context.Background(), labels)
	pprof.SetGoroutineLabels(ctx)

	go shortSleepLoop()
	go sleepLoop(time.Hour)
	go chanReceiveForever()
	go indirectShortSleepLoop()

	sleep := time.Second
	fmt.Printf("Sleeping for %s followed by gc\n", sleep)

	time.Sleep(time.Second)
	runtime.GC()

	fmt.Printf("Dump 1\n")
	if err := writeProfiles(1); err != nil {
		panic(err)
	}

	sleep = 1*time.Minute + 10*time.Second
	fmt.Printf("Sleeping for %s followed by gc\n", sleep)
	time.Sleep(sleep)
	runtime.GC()

	fmt.Printf("Dump 2\n")
	if err := writeProfiles(2); err != nil {
		panic(err)
	}

	fmt.Printf("Waiting forever\n")

	//fmt.Printf("waiting indefinitely so you can press ctrl+\\ to compare the output\n")
	//runtime.GC()
	<-errCh
}

func shortSleepLoop() {
	for {
		time.Sleep(time.Second)
	}
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

func indirectShortSleepLoop() {
	indirectShortSleepLoop2()
}

func indirectShortSleepLoop2() {
	go shortSleepLoop()
}
