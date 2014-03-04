package main

import (
	"fmt"
	. "net"
	"strings"
	"time"
	"os/exec"
)

func CheckError(err error, errorMsg string) {
	if err != nil {
		fmt.Println("!!Error type: " + errorMsg)
	}
}

func IAmAlive(port string) {
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
	conn.Write([]byte("I Am Alive!"))
	conn.Close()
}

func GetLocalIp() *string {
	conn, err := Dial("udp4", "google.com:80")
	CheckError(err, "ERROR: LocalIp: dialing to google.com:80")
	return &strings.Split(conn.LocalAddr().String(), ":")[0]
}

func ListenToNetwork(port string, incoming chan string, timeLimit int) {
	udpAddr, err := ResolveUDPAddr("udp4", ":"+port)
	CheckError(err, "ERROR while resolving UDPaddr for ListenToNetwork")
	fmt.Println("Establishing ListenToNetwork")

	conn, err := ListenUDP("udp4", udpAddr)
	fmt.Println("Listening on port ", udpAddr.String())
	CheckError(err, "Error while establishing listening connection")

	conn.SetReadDeadline(time.Now()+time.Milliseconds*timeLimit)
	
	data := make([]byte, 1024)
	//ownAddr := *GetLocalIp();
	for {
		_, _, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")
		//if addr.String() == ownAddr{ //OBS
		if err != nil{
			data = "main is dead"
		}
		fmt.Println("Channeling data " + string(data))
		incoming <- string(data)
		//}

	}
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
			data = "main is dead"
		}
		fmt.Println("Channeling data " + string(data))
		incoming <- string(data)
		//}

	}
}

func MainCounter(){
	incoming := make(chan string)
	countChan := make(chan int)
	alivePort := 22031
	 BackupCounter(incoming, countChan, alivePort)		//Må spawne backup
	cmd := exec.Command("mate-terminal", "-x", "go", "run", "telleProg.go")
	cmd.Run()
	//Må sende ImAlive
	//Må counte
	//Må sende CountUpdate
}

func BackupCounter(incoming chan string, countChan chan int, alivePort string){
	go ListenToNetwork(alivePort, 200)
	for{
		select{
		case update := <- incoming:
			if update == "main is dead"://Må oppdage død counter
				break
		case count := <- countChan://Må oppdatere count
		}
	}

	go MainCounter()
	
	//Må spawne Main
}

func main() {
	incoming := make(chan string)
	alivePort := "22222"

	addr := GetLocalIp()
	fmt.Println(*addr)
	go IAmAlive(alivePort)
	go ListenToNetwork(alivePort, incoming, 200)
	for {
		select {
		case data := <-incoming:
			fmt.Println("He said: ",data)
		}
	}

	fmt.Println("We're done here!")

}
