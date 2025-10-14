package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	rootCtx := context.Background()
	valCtx := context.WithValue(rootCtx, "main-key", "main-value")
	cancelCtx, cancel := context.WithCancel(valCtx)
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

	// accessing the value from the context
	fmt.Printf("[genNos] value of (main-key) from context : %v\n", ctx.Value("main-key"))

	subValCtx := context.WithValue(ctx, "main-key", "new-main-value")

	cWg := &sync.WaitGroup{}
	cWg.Add(1)
	oddCtx, cancel := context.WithTimeout(subValCtx, 10*time.Second)
	defer cancel()
	go genOddNos(cWg, oddCtx)

	cWg.Add(1)
	evenCtx, cancel := context.WithTimeout(subValCtx, 15*time.Second)
	defer cancel()
	go genEvenNos(cWg, evenCtx)

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
	cWg.Wait()
}

func genOddNos(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	// accessing the value from the context
	fmt.Printf("[genOddNos] value of (main-key) from context : %v\n", ctx.Value("main-key"))
	ticker := time.NewTicker(300 * time.Millisecond)
LOOP:
	for i := 1; ; i += 2 {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.Canceled:
				fmt.Println("[genOddNos] cancellation [programmatic] signal received!")
			case context.DeadlineExceeded:
				fmt.Println("[genOddNos] cancellation [timeout] signal received!")
			}
			break LOOP
		case <-ticker.C:
			fmt.Printf("[genOddNos] no : %d\n", i)
		}
	}
	ticker.Stop()
}

func genEvenNos(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	// accessing the value from the context
	fmt.Printf("[genEvenNos] value of (main-key) from context : %v\n", ctx.Value("main-key"))
	ticker := time.NewTicker(700 * time.Millisecond)
LOOP:
	for i := 0; ; i += 2 {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.Canceled:
				fmt.Println("[genEvenNos] cancellation [programmatic] signal received!")
			case context.DeadlineExceeded:
				fmt.Println("[genEvenNos] cancellation [timeout] signal received!")
			}
			break LOOP
		case <-ticker.C:
			fmt.Printf("[genEvenNos] no : %d\n", i)
		}
	}
	ticker.Stop()
}
