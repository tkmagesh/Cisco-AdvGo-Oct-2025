package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		time.Sleep(5 * time.Second)
		ch1 <- 100
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- 200
	}()

	go func() {
		time.Sleep(3 * time.Second)
		d3 := <-ch3
		fmt.Println("ch3 :", d3)
	}()

	/*
		wg := sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("ch1 :", <-ch1)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("ch2 :", <-ch2)
		}()
		wg.Wait()
	*/

	// Using "select-case" for the above
	time.Sleep(2500 * time.Millisecond) // wait until the data sent to "ch2" OR one of the channels
	for range 3 {
		select { //blocked until any of the channel operations in the "case" statement is successful (unblocked), unless a "default" section is there
		case d1 := <-ch1:
			fmt.Println("ch1 :", d1)
		case d2 := <-ch2:
			fmt.Println("ch2 :", d2)
		case ch3 <- 300:
			fmt.Println("[select] data sent to ch3")
		default:
			fmt.Println("[select] no channel operations are successful!")
		}
	}
}
