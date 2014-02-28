package main

import (
	"fmt"
	. "net"
	"strings"
	"time"
	"os/exec"
	"strconv"
)

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

func ListenToNetworkTimeLimited(port string, incoming chan string, timeLimit int) {
	udpAddr, err := ResolveUDPAddr("udp4", ":"+port)
	CheckError(err, "ERROR while resolving UDPaddr for ListenToNetwork")
	fmt.Println("Establishing ListenToNetwork")

	conn, err := ListenUDP("udp4", udpAddr)
	fmt.Println("Listening on port ", udpAddr.String())
	CheckError(err, "Error while establishing listening connection")

	conn.SetReadDeadline( time.Now().Add(time.Duration(timeLimit) * time.Millisecond) )
	
	data := make([]byte, 1024)
	//ownAddr := *GetLocalIp();
	for {
		_, _, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")
		//if addr.String() == ownAddr{ //OBS
		if err != nil{
			data = []byte("connection is dead")
			fmt.Println("Channeling data " + string(data))
			incoming <- string(data)
			break
		}
		fmt.Println("Channeling data " + string(data))
		incoming <- string(data)
		//}

	}
	conn.Close()
}





//MAIN PROGRAM
func main() {

	incoming := make(chan string)
	alivePort := "26030"
	sendPort := "26032"
	i := 0
	go ImAlive(alivePort)
	go ListenToNetworkTimeLimited(sendPort, incoming, 200)


	//Checks to see if there is old data in a backup
	continueCount := <-incoming
	fmt.Println(continueCount)
	if continueCount != "connection is dead"{
		i, _ = strconv.Atoi(continueCount)
	}else{
		i = 0
	}

	//Makes a new backup 
	fmt.Println("Creating backup")
	cmd := exec.Command("mate-terminal", "-x", "go", "run", "backup.go")
	cmd.Run()

	//Counts and updates over UDP
	for{
		SendToNetwork(sendPort, string(i))
		i += 1
		time.Sleep(time.Second*1)
		fmt.Println(i)
	}
}
