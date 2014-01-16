// Go 1.2
// go run helloworld_go.go

package main

import (
	. "fmt"     // Using '.' to avoid prefixing functions with their package names
	. "runtime" //   This is probably not a good idea for large projects...
	//. "time"
)

var chan1, chan2, chan4 chan int

var k int

func thread_1(quit chan int) {
	for x := 0; x < 999999; x++ {
		chan1 <- 1
	}
	quit <- 1
}
func thread_2(quit chan int) {
	for x := 0; x < 1000000; x++ {
		chan2 <- 1
	}
	quit <- 1
}

func server() {
	var i = 0
	for {
		select {
		case <-chan1:
			i++
		case <-chan2:
			i--
		case chan4 <- i:
		}
	}
}

func main() {
	GOMAXPROCS(NumCPU()) // I guess this is a hint to what GOMAXPROCS does...
	routineQuit := make(chan int)
	chan1 = make(chan int)
	chan2 = make(chan int)
	chan4 = make(chan int)
	Println("bap")

	go thread_1(routineQuit) // This spawns adder() as a goroutine
	go thread_2(routineQuit)
	go server()
	Println("bop")
	<-routineQuit
	<-routineQuit
	// No way to wait for the completion of a goroutine (without additional syncronization)
	// We'll come back to using channels in Exercise 2. For now: Sleep

	Println("Done_2:", <-chan4)
}
