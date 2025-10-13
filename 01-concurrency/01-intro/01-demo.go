package main

import (
	"fmt"
	"time"
)

func main() {
	go f1() //scheduling the exection of f1() through the scheduler (to be executed in future)
	f2()
	// when the main function execution ends, the application is shutdown (irrespective of any goroutines that are scheduled and waiting to be executed)

	// f1 will get the opportunity to executed with the current function (main) is blocked for any reason (cooperative multitasking)

	// block the execution of "this" function so that the scheduler will look for other goroutines scheduled and schedule them for execution
	// Poor man's synchronization techniques (DO NOT do this in your applications)
	time.Sleep(4 * time.Second)

}

func f1() {
	fmt.Println("f1 started")
	time.Sleep(3 * time.Second)
	fmt.Println("f1 completed")
}

func f2() {
	fmt.Println("f2 invoked")
}
