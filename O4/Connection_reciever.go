

package main


import (

	. "net"

)

var err error

func main() {
	listener, _ := Listen("tcp", ":33501")
	for{
		data := make([]byte, 1024)
		conn, _ := listener.Accept()
		_, _ = conn.Read(data)
		msg := string(data)
		if msg == "speil"{
			for{
				_, _ = conn.Read(data)
				msg = string(data)
				if string(msg) == "close"{ 
					break
				}
				conn.Write(data)
			}
		}
		
		
		conn.Close()
		
	}
}