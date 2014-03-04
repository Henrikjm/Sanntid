// Go 1.2
// go run helloworld_go.go

package main

import (
	. "fmt"     // Using '.' to avoid prefixing functions with their package names
	. "runtime" //   This is probably not a good idea for large projects...
	//. "time"
)

var i = 0

func thread_1(quit chan int) {
	for x := 0; x < 1000000; x++ {
		i++
		quit <- 1
	}
}
func thread_2(quit chan int) {
	for x := 0; x < 1000000; x++ {
		i--
		quit <- 1
	}
}

func main() {
	GOMAXPROCS(NumCPU()) // I guess this is a hint to what GOMAXPROCS does...
	routineQuit := make(chan int)
	go thread_1(routineQuit) // This spawns adder() as a goroutine
	go thread_2(routineQuit)
	<-routineQuit
	<-routineQuit
	// No way to wait for the completion of a goroutine (without additional syncronization)
	// We'll come back to using channels in Exercise 2. For now: Sleep

	Println("Done_2:", i)
}
