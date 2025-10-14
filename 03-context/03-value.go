package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	/*
		rootCtx := context.Background()
		valCtx := context.WithValue(rootCtx, "main-key", "main-value")
		timeoutCtx, cancel := context.WithTimeout(valCtx, 10*time.Second)
	*/
	ctx := context.Background()
	ctx = context.WithValue(ctx, "main-key", "main-value")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	go func() {
		fmt.Println("Hit ENTER to stop (before timeout occurs)...!")
		fmt.Scanln()
		cancel() // send the cancellation signal
	}()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	// go genNos(wg, timeoutCtx)
	go genNos(wg, ctx)
	wg.Wait()
	fmt.Println("Done!")
}

func genNos(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	fmt.Printf("[genNos] value of (main-key) from context : %v\n", ctx.Value("main-key"))
	ticker := time.NewTicker(500 * time.Millisecond)
LOOP:
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.Canceled:
				fmt.Println("[genNos] cancellation [programmatic] signal received!")
			case context.DeadlineExceeded:
				fmt.Println("[genNos] cancellation [timeout] signal received!")
			}
			break LOOP
		case <-ticker.C:
			fmt.Printf("[genNos] no : %d\n", i)
		}
	}
	ticker.Stop()
}
