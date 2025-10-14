package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	rootCtx := context.Background()
	cancelCtx, cancel := context.WithCancel(rootCtx)
	go func() {
		fmt.Println("Hit ENTER to stop...!")
		fmt.Scanln()
		cancel() // send the cancellation signal
	}()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go genNos(wg, cancelCtx)
	wg.Wait()
	fmt.Println("Done!")
}

func genNos(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	ticker := time.NewTicker(500 * time.Millisecond)
LOOP:
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("[genNos] cancellation signal received!")
			break LOOP
		case <-ticker.C:
			fmt.Printf("[genNos] no : %d\n", i)
		}
	}
	ticker.Stop()
}
