package main

import (
	"fmt"
	"sync"
	"time"
)

// share memory by communicating (replacing "communicate by sharing memory")

func main() {
	wg := &sync.WaitGroup{}
	stopCh := make(chan struct{})
	go func() {
		fmt.Println("Hit ENTER to stop...!")
		fmt.Scanln()
		// send the signal

		// stopCh <- struct{}{}
		close(stopCh)
	}()
	wg.Add(1)
	go genNos(wg, stopCh)
	wg.Wait()
	fmt.Println("Done!")
}

func genNos(wg *sync.WaitGroup, stopCh chan struct{}) {
	defer wg.Done()

LOOP:
	for i := 0; ; i++ {
		select {
		case <-stopCh:
			fmt.Println("[genNos] stop signal received...!")
			break LOOP
		default:
			time.Sleep(500 * time.Millisecond)
			fmt.Println("No :", i)
		}
	}
}
