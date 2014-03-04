package main

import (
	"fmt"
	//."net"
	//"strings"
	//"time"
	"strconv"
)
/*
func CheckError(err error, errorMsg string) {
	if err != nil {
		fmt.Println("!!Error type: " + errorMsg)
	}
}

func ImAlive(port string) {
	fmt.Println("Establishing IAmAlive...")
	sendAddr, err := ResolveUDPAddr("udp4", "129.241.187.255:"+port)
	CheckError(err, "ERROR while resolving UDP addr")
	conn, err := DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")
	for {
		time.Sleep(time.Second * 3)
		conn.Write([]byte("I Am Alive!"))
	}
}

func SendToNetwork(port string, msg string) {
	sendAddr, err := ResolveUDPAddr("udp4", "129.241.187.255:"+port)
	CheckError(err, "ERROR while resolving UDP addr")
	conn, err := DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")
	conn.Write([]byte(msg))
	conn.Close()
}

func GetLocalIp() *string {
	conn, err := Dial("udp4", "google.com:80")
	CheckError(err, "ERROR: LocalIp: dialing to google.com:80")
	return &strings.Split(conn.LocalAddr().String(), ":")[0]
}

func ListenToNetwork(port string, incoming chan string) {
	udpAddr, err := ResolveUDPAddr("udp4", ":"+port)
	CheckError(err, "ERROR while resolving UDPaddr for ListenToNetwork")
	fmt.Println("Establishing ListenToNetwork")

	conn, err := ListenUDP("udp4", udpAddr)
	fmt.Println("Listening on port ", udpAddr.String())
	CheckError(err, "Error while establishing listening connection")
	
	data := make([]byte, 1024)
	//ownAddr := *GetLocalIp();
	for {
		_, _, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")
		//if addr.String() == ownAddr{ //OBS
		if err != nil{
			data = []byte("connection is dead")
		}
		fmt.Println("Channeling data " + string(data))
		incoming <- string(data)
	}
	conn.Close()
}

func ListenToNetworkTimeLimited(port string, outgoing chan string, timeLimit int) {
	udpAddr, err := ResolveUDPAddr("udp4", ":"+port)
	CheckError(err, "ERROR while resolving UDPaddr for ListenToNetwork")
	fmt.Println("Establishing ListenToNetwork")

	conn, err := ListenUDP("udp4", udpAddr)
	fmt.Println("Listening on port ", udpAddr.String())
	CheckError(err, "Error while establishing listening connection")

	
	
	data := make([]byte, 1024)
	//ownAddr := *GetLocalIp();
	for {
		conn.SetReadDeadline( time.Now().Add(time.Duration(timeLimit) * time.Millisecond) )
		_, _, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")
		//if addr.String() == ownAddr{ //OBS
		if err != nil{
			fmt.Println("Channeling data: connection is dead")
			outgoing <- "connection is dead"
			fmt.Println("Backup: Breaking listen-loop")
			break
		}
		fmt.Println("Channeling data: " + string(data))
		outgoing <- string(data)
		//}

	}
	conn.Close()
}


func MakeListenerConn(port string) Conn{
	udpAddr, err := ResolveUDPAddr("udp4", ":"+port)
	CheckError(err, "ERROR while resolving UDPaddr for ListenToNetwork")
	fmt.Println("Establishing ListenToNetwork")
	conn, err := ListenUDP("udp4", udpAddr)
	fmt.Println("Listening on port ", udpAddr.String())
	CheckError(err, "Error while establishing listening connection")
	return conn
}


*/



func main() {
	variablelele := "8"
	

	i,_ := strconv.Atoi(variablelele)
	fmt.Println("variablelele = ", variablelele,"i = ", i)

}
