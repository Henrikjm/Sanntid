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
		println(string(data))
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
	//part1	funker
	conn, err := Dial("tcp", "129.241.187.161:33546")
	checkError(err, "dialupError")
	writeMsg(conn, "\nFakka yuu devil machine!!!\x00")
	readMsg(conn)

	//part2 under arbeid
	Dialup := "\nConnect to: 129.241.187.148:33546\x00"
	writeMsg(conn, Dialup)
	readMsg(conn)

	
}
