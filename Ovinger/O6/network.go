package network

import (
	"fmt"
	"net"
	"time"
)

func CheckError(err error, errorMsg string) {
	if err != nil {
		fmt.Println("!!Error type: " + errorMsg,"!!")
	}
}

func ImAliveUDP(port string) {
	fmt.Println("Establishing IAmAlive...")
	sendAddr, err := ResolveUDPAddr("udp4", "129.241.187.255:"+port)
	CheckError(err, "ERROR while resolving UDP addr")
	conn, err := DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")
	msg := GetLocalIp() + " :IsAlive"
	for {
		time.Sleep(time.Millisecond * 100)
		conn.Write([]byte(msg))
	}
}

func SendToNetworkUDP(port string, msg string) {
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

func ListenToNetworkUDP(conn *UDPConn, incoming chan string) {
	data := make([]byte, 1024)
	//ownAddr := *GetLocalIp();
	for {
		_, _, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")
		//if addr.String() == ownAddr{ //OBS
		if err != nil{
			data = []byte("connection is dead")
		}
		fmt.Println("ListenToNetwork: Channeling data :" + string(data))
		incoming <- string(data)
	}
	conn.Close()
}



func MakeListenerConn(port string) *UDPConn{
	udpAddr, err := ResolveUDPAddr("udp4", ":"+port)
	CheckError(err, "ERROR while resolving UDPaddr for ListenToNetwork")
	fmt.Println("Establishing ListenToNetwork")
	conn, err := ListenUDP("udp4", udpAddr)
	fmt.Println("Listening on port ", udpAddr.String())
	CheckError(err, "Error while establishing listening connection")
	return conn
}