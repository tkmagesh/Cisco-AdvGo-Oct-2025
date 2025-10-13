/*
The prime check for each number is executed sequentially
Modify the below to execute them concurrently
*/
package main

import (
	"fmt"
)

func main() {
	var start, end int = 2, 1000
	for no := start; no <= end; no++ {
		PrintIfPrime(no)
	}
}

func PrintIfPrime(no int) {
	for i := 2; i <= (no / 2); i++ {
		if no%i == 0 {
			return
		}
	}
	fmt.Printf("Prime No : %d\n", no)
}
