package main

import (
	"time"
)

func main() {
	start := time.Now()
	for time.Since(start) < time.Second {
	}
	//time.Sleep(time.Second)
}
