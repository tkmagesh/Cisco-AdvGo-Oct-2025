package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ch := getNos()
	for data := range ch {
		time.Sleep(500 * time.Millisecond)
		fmt.Println(data)
	}
	fmt.Println("Done!")
}

func getNos() <-chan int {
	ch := make(chan int)
	count := rand.Intn(20)
	fmt.Printf("[getNos] processing %d numbers\n", count)
	go func() {
		for i := range count {
			ch <- (i + 1) * 10
		}
		close(ch)
	}()
	return ch
}
