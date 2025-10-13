/*
Modify as below:
do not print the prime number in the "PrintIfPrime()"
instead, collect all the prime numbers on a slice
and print them in the "main()"
*/
package main

import (
	"fmt"
	"sync"
)

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
