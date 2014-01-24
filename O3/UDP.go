package main

import( 
	. "net"
	. "fmt"
)

//var err error

func writeMsgUDP(conn Conn, msg string){
	data := []byte(msg)
	_,err := conn.Write(data)
	if err != nil {
		Println("!!Error type: writeError")
	}
}

func readMsgUDP(conn Conn){
	data := make([]byte, 1024)
	_,err := conn.Read(data)
	if err != nil {
		Println("!!Error type: readError")
	}
	Println(string(data))
}




/*
func reciever(){
	buffer = make([]byte, 1024)	
}
*/

func main(){
	
	conn, err := Dial("udp", "129.241.187.161:20018")
	if err != nil {
		Println("!!Error type: connectError")
	}
	writeMsgUDP(conn, "hallais")
	readMsgUDP(conn)
	


}
