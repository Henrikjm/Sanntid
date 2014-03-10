package channelReceiver

import(
	"fmt"
)

func Receiver2Handle(channel *chan string){
	//var dummy string
	dummy = <- channel
	fmt.Println(dummy)

}