package main

import (
	"fmt"
	."net"
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
		time.Sleep(time.Millisecond * 100)
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

func ListenToNetwork(conn *UDPConn, incoming chan string) {
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

func ListenToNetworkTimeLimited(conn *UDPConn, outgoing chan string, timeLimit int) {
	data := make([]byte, 1024)
	//ownAddr := *GetLocalIp();
	for {
		conn.SetReadDeadline( time.Now().Add(time.Duration(timeLimit) * time.Millisecond) )
		_, _, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")
		//if addr.String() == ownAddr{ //OBS
		if err != nil{
			fmt.Println("Channeling data: connection is dead")
			go func(outgoing chan string){
				outgoing <- "connection is dead"
			}(outgoing)
			fmt.Println("Backup: Breaking listen-loop")
			break
		}
		//fmt.Println("Channeling data: " + string(data))
		outgoing <- string(data)
		//}

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



//MAIN PROGRAM
func main() {

	incoming := make(chan string)
	alivePort := "26030"
	countPort := "26032"
	countConn := MakeListenerConn(countPort)
	var i int

	go ImAlive(alivePort)
	go ListenToNetworkTimeLimited(countConn, incoming, 500)


	//Checks to see if there is old data in a backup
	continueCount := <-incoming
	/*
	fmt.Println(continueCount)
	if continueCount != "connection is dead"{
		a, _ := strconv.Atoi(continueCount)
		i = i + a
		fmt.Println("string", continueCount, "int", i)
	}else{		
		i = 0
	}
	fmt.Println(i)
	*/
	if continueCount != "connection is dead"{
		continueCount = "0"
	}

	i,_ = strconv.Atoi(continueCount)
	fmt.Println("continueCount == ", continueCount, "i == ", i)

	//Makes a new backup 
	countConn.Close()
	fmt.Println("Creating backup")
	cmd := exec.Command("mate-terminal", "-x", "go", "run", "backup.go")
	cmd.Run()

	//Counts and updates over UDP
	go func(countPort string, i int){
		for{
			SendToNetwork(countPort, strconv.Itoa(i))
			i += 1
			time.Sleep(time.Second*1)
			fmt.Println(i)
		}
	}(countPort, i)
	var exit string
	fmt.Scanln(&exit)
}

