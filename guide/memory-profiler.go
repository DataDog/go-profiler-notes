//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func main() {
	file, _ := os.Create("./memory-profiler.pprof")
	defer pprof.Lookup("allocs").WriteTo(file, 0)
	defer runtime.GC()

	go allocSmall()
	go allocBig()

	time.Sleep(1 * time.Second)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%#v\n", m)
}

//go:noinline
func allocSmall() {
	for i := 0; ; i++ {
		_ = alloc(32)
	}
}

//go:noinline
func allocBig() {
	for i := 0; ; i++ {
		_ = alloc(256)
	}
}

//go:noinline
func alloc(size int) []byte {
	return make([]byte, size)
}
