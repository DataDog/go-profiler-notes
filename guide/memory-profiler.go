// +build ignore

package main

import (
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	file, _ := os.Create("./mem.pprof")
	defer pprof.Lookup("allocs").WriteTo(file, 0)

	go allocSmall()
	go allocBig()

	time.Sleep(1 * time.Second)
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
