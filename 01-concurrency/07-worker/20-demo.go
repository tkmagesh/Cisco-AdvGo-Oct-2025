/*
Create a handful (customizable) number of goroutines that will share the load of processing the numbers (instead of one goroutine per number)
*/
package main

import (
	"fmt"
	"sync"
)

func main() {
	var start, end int = 2, 1000
	primesCh := generatePrimes(start, end, 10)
	for primeNo := range primesCh {
		fmt.Println("Prime No :", primeNo)
	}
}

func dataProducer(start, end int) <-chan int {
	dataCh := make(chan int)

	// feed data to the data channel
	go func() {
		for no := start; no <= end; no++ {
			dataCh <- no
		}
		close(dataCh)
	}()
	return dataCh
}

func worker(workerId int, wg *sync.WaitGroup, dataCh <-chan int, primeCh chan<- int) {
	defer wg.Done()
	fmt.Printf("[worker] worker : %d starting....\n", workerId)
	for no := range dataCh {
		if isPrime(no) {
			fmt.Printf("[worker] worker id : %d and prime no : %d\n", workerId, no)
			primeCh <- no
		}
	}
	fmt.Printf("[worker] worker : %d completed\n", workerId)
}

func generatePrimes(start, end int, workerCount int) <-chan int {
	primesCh := make(chan int)
	dataCh := dataProducer(start, end)
	wg := &sync.WaitGroup{}
	for id := range workerCount {
		wg.Add(1)
		go worker(id+1, wg, dataCh, primesCh)
	}
	go func() {
		wg.Wait() // wait for all the workers to complete
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
