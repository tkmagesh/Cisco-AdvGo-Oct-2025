/*
The prime check for each number is executed sequentially
Modify the below to execute them concurrently
*/
package main

import (
	"fmt"
	"sync"
)

func main() {
	var start, end int = 2, 1000
	wg := &sync.WaitGroup{}
	for no := start; no <= end; no++ {
		wg.Add(1)
		go PrintIfPrime(no, wg)
	}
	wg.Wait()
}

func PrintIfPrime(no int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 2; i <= (no / 2); i++ {
		if no%i == 0 {
			return
		}
	}
	fmt.Printf("Prime No : %d\n", no)
}
