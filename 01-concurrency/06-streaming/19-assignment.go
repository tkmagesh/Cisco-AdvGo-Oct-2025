/*
Modify the solution to adopt "share memory by communicating" (using channels)
*/
package main

import (
	"fmt"
	"sync"
)

func main() {
	var start, end int = 2, 1000
	primesCh := generatePrimes(start, end)
	for primeNo := range primesCh {
		fmt.Println("Prime No :", primeNo)
	}
}

func generatePrimes(start, end int) <-chan int {
	primesCh := make(chan int)
	go func() {
		wg := sync.WaitGroup{}
		for no := start; no <= end; no++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if isPrime(no) {
					primesCh <- no
				}
			}()
		}
		wg.Wait()
		close(primesCh)
	}()
	return primesCh
}

func isPrime(no int) bool {
	for i := 2; i <= (no / 2); i++ {
		if no%i == 0 {
			return false
		}
	}
	return true
}
