//go:build ignore

package main

import (
	"net/http"
	_ "net/http/pprof"
)

func main() {
	http.ListenAndServe("localhost:1234", nil)
}
