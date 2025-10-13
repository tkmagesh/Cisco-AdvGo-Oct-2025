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

func main() {
	var start, end int = 2, 1000
	var wg sync.WaitGroup
	for no := start; no <= end; no++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			PrintIfPrime(no)
		}()
	}
	wg.Wait()
}

func PrintIfPrime(no int) {
	for i := 2; i <= (no / 2); i++ {
		if no%i == 0 {
			return
		}
	}
	fmt.Printf("Prime No : %d\n", no)
}
