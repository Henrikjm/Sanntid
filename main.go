package main 

import (
	"fmt"
	"net"
	"time"
)

var connection *net.UDPConn

func server() {
	fmt.Println("Starting udp server")
	udpAddr, err := net.ResolveUDPAddr("udp4", ":20018")
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	
	connection = conn
	
	buffer := make([]byte, 1024)
	for {
		len, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		
		fmt.Println("Got", len, "from", addr, "bytes:", string(buffer))
	}
}

func main() {
	go server()
	
	sendAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:20018")
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	
	time.Sleep(1*time.Second)
	
	if connection == nil {
		fmt.Println("Error blah")
		return
	}
	
	for {
		connection.WriteToUDP([]byte("Hello :D"), sendAddr)
		time.Sleep(1*time.Second)
	}
}

