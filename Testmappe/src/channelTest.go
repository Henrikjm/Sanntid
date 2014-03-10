package main

import(
	//"fmt"
	"channelReceiver"
)
/*
var channel chan string
var exitChan chan int
*/

func main(){
	channel := make(chan string)
	exitChan := make(chan int)
	go channelReceiver.Receiver1Handle(channel)
	go channelReceiver.Receiver2Handle(channel)
	

	<-exitChan
	//var exit string
	//fmt.Scanln(&exit)
}