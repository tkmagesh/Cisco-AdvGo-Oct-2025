package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var no atomic.Int64

func main() {
	wg := &sync.WaitGroup{}
	for range 200 {
		wg.Add(1)
		go increment(wg)
	}
	wg.Wait()
	fmt.Println("no :", no.Load())
}

func increment(wg *sync.WaitGroup) {
	defer wg.Done()
	no.Add(1)
}
