package main

import (
	"fmt"
	"sync"
)

func mutexWorkload(goroutines, ops, depth int) error {
	if goroutines%2 != 0 {
		return fmt.Errorf("bad goroutines: %d: must be a multiple of 2", goroutines)
	}

	wg := &sync.WaitGroup{}
	for j := 0; j < goroutines/2; j++ {
		m := &sync.Mutex{}
		wg.Add(1)
		go atStackDepth(depth, func() {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				m.Lock()
				m.Unlock()
			}
		})
		wg.Add(1)
		go atStackDepth(depth, func() {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				m.Lock()
				m.Unlock()
			}
		})
	}
	wg.Wait()
	return nil
}
