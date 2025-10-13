/*
Modify the solution to adopt "share memory by communicating" (using channels)
*/
package main

import (
	"fmt"
	"sync"
)

// communicate by sharing memory
var primes []int
var mutex sync.Mutex

func main() {
	var start, end int = 2, 1000
	var wg sync.WaitGroup
	for no := start; no <= end; no++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			checkPrime(no)
		}()
	}
	wg.Wait()
	for _, primeNo := range primes {
		fmt.Println("Prime No :", primeNo)
	}
}

func checkPrime(no int) {
	for i := 2; i <= (no / 2); i++ {
		if no%i == 0 {
			return
		}
	}
	mutex.Lock()
	{
		primes = append(primes, no)
	}
	mutex.Unlock()
}
