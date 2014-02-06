// Test of TCP-module

package main


import(
	. "net"
	. "fmt"
	//. "time"
)


func TCP_hookup( adress string) Conn {
	conn, err := Dial("tcp", adress)
	checkError(err, "Dial Error")
	return conn
}


func TCP_listener( port string ) Conn {
	port = ":" + port
	listener, err := Listen("tcp", port)
	conn, err := listener.Accept()
	checkError(err, "Listen Error")
	return conn
}


func checkError(err error, errorMsg string) {
	if err != nil {
		Println("!!Error type: " + errorMsg)
	}
}


func writeMsg(conn Conn, msg string) {
	data := []byte(msg)
	_, err = conn.Write(data)
	checkError(err, "writeError")
}


func main(){
	 TCP_hookup("129.241.187.149:33501")
}

