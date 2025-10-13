package main

import (
	"fmt"
	"sync"
)

// custom type encapsulating the state & concurrent safe manipulation of the state
type SafeCounter struct {
	sync.Mutex
	no int
}

func (sf *SafeCounter) Add(delta int) {
	sf.Lock()
	{
		sf.no += 1
	}
	sf.Unlock()
}

var sf SafeCounter

func main() {
	wg := &sync.WaitGroup{}
	for range 200 {
		wg.Add(1)
		go increment(wg)
	}
	wg.Wait()
	fmt.Println("no :", sf.no)
}

func increment(wg *sync.WaitGroup) {
	defer wg.Done()
	sf.Add(1)
}
