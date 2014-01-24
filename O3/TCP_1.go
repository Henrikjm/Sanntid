
package main


import (
	. "fmt"
	. "net"
	"time"

)

var err error

func readMsg(conn Conn) {
	
		data := make([]byte, 1024)
		_, err := conn.Read(data)
		checkError(err, "readError")
		Println(string(data))
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

func readSimon(conn Conn) {
		data := make([]byte, 1024)
		_, err := conn.Read(data)
		checkError(err, "readError")
		sub := string(data)
		Println(sub)
		simon := "Simon says "
		Println(sub[len(simon):])
		if sub[0:len(simon)] == simon{
			sub = sub[len(simon):]
			Println(sub)			
			writeMsg(conn, sub)
		return
		}
		Println("no match")
		writeMsg(conn, "")

}
		
	


func main() {
	//part1	
	conn, err := Dial("tcp", "129.241.187.161:33546")
	checkError(err, "dialupError")
	writeMsg(conn, "\nPart 1 check\n \x00")
	readMsg(conn)

	//part2
	listener, err := Listen("tcp", ":33504") //må endre port iblant (når panic:)
	checkError(err, "ListenError")
	writeMsg(conn, "Connect to: 129.241.187.151:33504\x00") //må endres hvis man sitter på ny PC
	channel, err := listener.Accept()
	checkError(err, "AcceptError")
	readMsg(channel)
	writeMsg(channel, "\nPart2 check\x00")
	readMsg(channel)
	
	//part3
	writeMsg(channel, "Play Simon says\x00")
	for{
		time.Sleep(1 * time.Second)	
		readSimon(channel)
	}
}
