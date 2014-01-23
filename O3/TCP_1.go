
package main


import (
	. "fmt"
	. "net"

)

var err error

func readMsg(conn Conn) {
	
		data := make([]byte, 1024)
		_, err := conn.Read(data)
		checkError(err, "readError")
		Println(string(data))
		
		
	
}

func readSimon(conn Conn) {
	
		data := make([]byte, 1024)
		_, err := conn.Read(data)
		checkError(err, "readError")
		sub := string(data)
		simon := "Simon says "
		println(sub)

		equal := true

		for i :=0; i<len(simon) ; i++{
			if simon[i] != sub[i]{
				equal = false
			}		
		}	
		
		println(sub)		
	
		if equal == true{
			_, err = conn.Write(data)			
			//writeMsg(conn, sub)
		}		
	}
		
	


func writeMsg(conn Conn, msg string) {
	data := []byte(msg)
	_, err = conn.Write(data)
	checkError(err, "writeError")
}

func checkError(err error, errorMsg string) {
	if err != nil {
		Println("!!Error type: " + errorMsg)
	}
}



func main() {
	//part1	
	conn, err := Dial("tcp", "129.241.187.161:33546")
	checkError(err, "dialupError")
	writeMsg(conn, "\nPart 1\n \x00")
	readMsg(conn)

	//part2
	

	listener, err := Listen("tcp", "0:33500")
	checkError(err, "ListenError")	
	writeMsg(conn, "Connect to: 129.241.187.148:33500\x00")
	channel, err := listener.Accept()
	checkError(err, "AcceptError")
	readMsg(channel)
	writeMsg(channel, "Play Simon says\x00")


	for{
	
		//time.Sleep(3 * time.Second)
		//writeMsg(channel,"Hello from client!\x00")

		
		readSimon(channel)
		
		
	}

	
}
