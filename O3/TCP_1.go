package main

import (
	. "fmt"
	. "net"
)

var err error

func readMsg(conn Conn) {
	for {
		data := make([]byte, 1024)
		_, err = conn.Read(data)
		checkError(err, "readError")
		msg := string(data[:])
		println(msg)
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
	//Dialup := "Connect to: 129.241.187.161:33546"
	conn, err := Dial("tcp", "129.241.187.161:33546")
	checkError(err, "dialupError")

	readMsg(conn)

	conn.Close()

}
