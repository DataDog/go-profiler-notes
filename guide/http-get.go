//go:build ignore

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/trace"
	"time"
)

func main() {
	traceF, _ := os.Create("trace.out")
	trace.Start(traceF)
	defer trace.Stop()

	start := time.Now()
	res, err := http.Get("https://example.org/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", res.StatusCode)
	log.Printf("main() took: %s\n", time.Since(start))
}
