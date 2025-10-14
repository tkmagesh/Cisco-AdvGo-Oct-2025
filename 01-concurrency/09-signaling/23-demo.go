package main

import (
	"fmt"
	"sync"
	"time"
)

// communicate by sharing memory
var stop bool = false

func main() {
	wg := &sync.WaitGroup{}
	go func() {
		fmt.Println("Hit ENTER to stop...!")
		fmt.Scanln()
		stop = true
	}()
	wg.Add(1)
	go genNos(wg)
	wg.Wait()
	fmt.Println("Done!")
}

func genNos(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; ; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("No :", i)
		if stop {
			fmt.Println("[genNos] stop signal received...!")
			break
		}
	}
}
