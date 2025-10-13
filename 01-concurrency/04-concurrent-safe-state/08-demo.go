package main

import (
	"fmt"
	"sync"
)

var no int
var mutex sync.Mutex

func main() {
	wg := &sync.WaitGroup{}
	for range 200 {
		wg.Add(1)
		go increment(wg)
	}
	wg.Wait()
	fmt.Println("no :", no)
}

func increment(wg *sync.WaitGroup) {
	defer wg.Done()
	mutex.Lock()
	{
		no += 1
	}
	mutex.Unlock()
}
