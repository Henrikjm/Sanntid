
package main

import(
	"fmt"
	"time"
	"network"
)

func main(){
	

	var msg string
	port := "09001"
	var incoming chan string
	conn := network.MakeListenerConn(port)
	go network.ListenToNetworkUDP(conn, &incoming)
	go func(port string){
		for{
			network.SendToNetworkUDP(port, "helu")
			time.Sleep(time.Millisecond*500)
		}
	}(port)

	
	go func(incoming chan string){
		for{
			select{
			case msg = <-incoming:
				fmt.Println(msg)
			default:
			}
		}
	}(incoming)
	

	
	var exit string
	fmt.Scanln(&exit)

}