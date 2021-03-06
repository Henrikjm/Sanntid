package main

import(
	."fmt"
	."net"
	"bufio"
	"os"
	"strings"
)

func ListenToNetwork(chanCon chan *UDPConn, setting string, port string) string{
	Println("Start UDP server")
	udpAddr, err := ResolveUDPAddr("udp4", ":" + port) //resolving
	CheckError("ERROR: Resolving error")
	conn, err := ListenUDP("udp", udpAddr) //initiating listening
	CheckError("ERROR: Listening error")
   	chanCon <- conn
	for{
		data := make([]byte,1024)
		_, addr, err := conn.ReadFromUDP(data) //kan bruke addr til å sjekke hvor melding kommer fra f.eks if addr not [egen i.p]
		CheckError("ERROR: while reading")
		return string(data)
	}	
}

func SendToNetwork(chanCon chan *UDPConn, port string, string msg){
    sendAddr, err := ResolveUDPAddr("udp4","129.241.187.255:" + port) //Spesifiserer adresse
	CheckError("ERROR while resolving UDP addr")
	connection := <- chanCon
	if msg == "writeFromConsole" {
		reader := bufio.NewReader(os.Stdin)
		for{
			text, _ := reader.ReadString('\n')
			testmsg := []byte(strings.TrimSpace(text))
			if connection ==  nil{
				println("!!ERROR, connection = nil")
			}
			if testmsg == "exit"{	
				return
			}
			connection.WriteToUDP(testmsg, sendAddr)
		}
	}
		if connection ==  nil{
			println("!!ERROR, connection = nil")
	}
	connection.WriteToUDP(testmsg, sendAddr)
}

func CheckError(err error, errorMsg string) {
	if err != nil {
		Println("!!Error type: " + errorMsg)
	}
}

