package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/trace"
)

func main() {
	trace.Start(os.Stderr)
	defer trace.Stop()

	res, err := http.Get("https://example.org/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", res.StatusCode)
}
