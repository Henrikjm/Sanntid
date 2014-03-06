package main

import "net"
import "fmt"
import "time"

func main(){
	udpAddr, err := net.ResolveUDPAddr("udp4", ":"+"20202")
	//CheckError(err, "ERROR while resolving UDPaddr for ListenToNetwork")
	fmt.Println("Establishing ListenToNetwork")
	conn, err := net.ListenUDP("udp4", udpAddr)
	fmt.Println("Listening on port ", udpAddr.String())
	//CheckError(err, "Error while establishing listening connection")

	conn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Millisecond))
	data := make([]byte, 1024)
	_,_,err = conn.ReadFromUDP(data)
	fmt.Println("Error says: ", err)
	kek := err.Error()
	port := "20202"
	fmt.Println( kek == "read udp4 0.0.0.0:"+port+": i/o timeout")
}