package main

import (
	"fmt"
	"sync"
)

func chanWorkload(goroutines, ops, depth, bufsize int) error {
	if goroutines%2 != 0 {
		return fmt.Errorf("bad goroutines: %d: must be a multiple of 2", goroutines)
	}

	wg := &sync.WaitGroup{}
	for j := 0; j < goroutines/2; j++ {
		ch := make(chan struct{}, bufsize)
		wg.Add(1)
		go atStackDepth(depth, func() {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				ch <- struct{}{}
			}
		})
		wg.Add(1)
		go atStackDepth(depth, func() {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				<-ch
			}
		})
	}
	wg.Wait()
	return nil
}
