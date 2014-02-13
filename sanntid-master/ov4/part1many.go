package main

import (
	"time"
	"fmt"
)

func sillyRoutine(c chan int) {
	i := <-c
	for {
		fmt.Printf("I'm number: %d\n", i)
		<-time.After(2*time.Second)
	}
}

const numGoroutines = 1e4

func main() {
	myChan := make(chan int, numGoroutines)
	blocker := make(chan int)
	for i := 1; i <= numGoroutines; i++ {
		myChan <- i
		go sillyRoutine(myChan)
	}
	<-blocker
}
