package main

import "fmt"
import "time"

func server(listenChan <-chan int, sendChan chan<- int) {
	sharedVar := 0
	for {
		select {
		case val := <-listenChan:
			sharedVar = sharedVar + val
		case sendChan<-sharedVar:
		}
	}
}

func printer(recvChan <-chan int) {
	for {
		i := <-recvChan
		fmt.Printf("Shared var: %d\n", i)
		<-time.After(time.Second)
	}
}

func client(value int, sendChan chan<- int) {
	for {
		<-time.After(1*time.Second)
		sendChan <- value
	}
}

func main() {
	c := make(chan int)
	s := make(chan int)
	blocker := make(chan int)
	go server(c, s)
	go printer(s)
	for i:=0; i <= 5; i++ {
		go client(i, c)
		<-time.After(time.Second)
	}
	<-blocker
}
