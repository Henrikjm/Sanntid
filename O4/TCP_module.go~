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
	Println("Connection succesfull: ", adress)
	return conn
	
	
}

func main(){
	 TCP_hookup("129.241.187.162:33546")
	
}

