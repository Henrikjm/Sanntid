// Test of TCP-module

package main

import(
	. "net"
	. "fmt"
	//. "time"
)



func TCP_hookup( adress string) Conn {
	conn, err := Dial("tcp", adress)
	if err != nil {
		//Handle connection error
		Println("Connection Error")
	}
	//Println("Connection succesfull: ", adress)
	return conn
}

/*
func TCP_listener( port string ) Conn {
	listener, err := Listen("tcp", port )
	if err != nil{
		//Handle listener error
		Println("Listener error")
	}
	channel, err := listener.Accept()
		if err != nil{
		//Handle listener accept error
		Println("Listener Accept Error")
	}
*/ 

	func main(){
	 TCP_hookup("129.241.187.149:33501")
}

