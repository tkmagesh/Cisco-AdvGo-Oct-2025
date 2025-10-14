package main

import (
	"fmt"
	"sync"
	"time"
)

// share memory by communicating (replacing "communicate by sharing memory")

func main() {
	wg := &sync.WaitGroup{}
	stopCh := time.After(5 * time.Second)
	wg.Add(1)
	go genNos(wg, stopCh)
	wg.Wait()
	fmt.Println("Done!")
}

func genNos(wg *sync.WaitGroup, stopCh <-chan time.Time) {
	defer wg.Done()
	tickerCh := time.Tick(500 * time.Millisecond)
LOOP:
	for i := 0; ; i++ {
		select {
		case <-stopCh:
			fmt.Println("[genNos] stop signal received...!")
			break LOOP
		case <-tickerCh:
			fmt.Println("Time :", i)
		}
	}
}
