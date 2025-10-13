/*
The prime check for each number is executed sequentially
Modify the below to execute them concurrently
*/
package main

import "fmt"

func main() {
	var start, end int = 2, 100
LOOP:
	for no := start; no <= end; no++ {
		for i := 2; i <= (no / 2); i++ {
			if no%i == 0 {
				continue LOOP
			}
		}
		fmt.Printf("Prime No : %d\n", no)
	}
}
