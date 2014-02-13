package main

import (
	"time"
	"fmt"
)

func numbaOne(c chan int) {
	for{
		i := <-c
		c <- i+1
		time.Sleep(1*time.Second)
		fmt.Printf("Numba one: %d\n", i)
	}
}

func numbaNone(c chan int) {
	for{
		i := <-c
		c <- i+2
		<-time.After(1*time.Second)
		fmt.Printf("Numba none: %d\n", i)
	}
}

func main() {
	c := make(chan int)
	go numbaOne(c)
	go numbaNone(c)
	c <- 10
	blocker := make(chan int)

	<-blocker
}
