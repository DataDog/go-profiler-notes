//go:build ignore

package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var buf [][]byte
var ballast []byte

func main() {
	file, _ := os.Create("mem.pprof")
	space := 5 * 1024 * 1024 // 5 MiB
	allocs := 1000

	ballast = make([]byte, space*5)
	runtime.GC()

	log.Println("start alloc")
	for i := 0; i < allocs; i++ {
		buf = append(buf, allocBytes(space/allocs))
	}
	runtime.GC() // without this the allocs above are invisible
	log.Println("end alloc")

	pprof.WriteHeapProfile(file)
}

func allocBytes(n int) []byte {
	return make([]byte, n)
}
