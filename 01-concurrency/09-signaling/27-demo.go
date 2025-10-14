package main

import (
	"fmt"
	"time"
)

// share memory by communicating (replacing "communicate by sharing memory")

func main() {
	stopCh := time.After(5 * time.Second)
	doneCh := genNos(stopCh)
	<-doneCh
	fmt.Println("Done!")
}

func genNos(stopCh <-chan time.Time) <-chan struct{} {
	doneCh := make(chan struct{}) // to notify that the goroutine is completed
	go func() {
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
		close(doneCh)
	}()
	return doneCh
}
