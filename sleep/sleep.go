package main

import (
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			sleep()
			wg.Done()
		}()
	}
	wg.Wait()
}

func sleep() {
	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)
	}
}
