package channelReceiver

import "fmt"

func Receiver1Handle(channel *chan string){
	
	dummy := "poop"
	channel <- dummy
	fmt.Println("skrevet")
	
}