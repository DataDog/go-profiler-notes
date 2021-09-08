// +build ignore

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var (
		sum int32
		wg  = &sync.WaitGroup{}
	)

	wg.Add(2)
	go add(&sum, 23, wg)
	go add(&sum, 42, wg)
	wg.Wait()

	fmt.Println(sum)
}

func add(dst *int32, delta int32, wg *sync.WaitGroup) {
	atomic.AddInt32(dst, delta)
	wg.Done()
}
