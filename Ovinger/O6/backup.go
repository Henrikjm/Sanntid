package main

import (
	"fmt"
	."net"
	"strings"
	"time"
	"os/exec"
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


//BACKUP PROGRAM
func main(){
	//Initializing variables and Listening channels
	var(
		update string
		count string
	)
	alivePort := "26030"
	aliveConn := MakeListenerConn(alivePort)
	countPort := "26032"
	countConn := MakeListenerConn(countPort)
	incoming := make(chan string)
	countChan := make(chan string)
	go ListenToNetworkTimeLimited(aliveConn, incoming, 500) //receiving alivemsg
	go ListenToNetwork(countConn, countChan) // receiving count

	//Updating and listening
	count = func(incoming chan string, countChan chan string, count string) string{
		for{
			select{
			case update = <- incoming:
				//fmt.Println("Backup: Recieved ",update)
				if update == "connection is dead"{//Må oppdage død counter
					fmt.Println("You have entered the IF statement")
					return count 
				}
			case count = <- countChan://Må oppdatere count
				fmt.Println("Recieved: ",count)
			default:
			}
		}
	}(incoming, countChan, count)
	countConn.Close()
	aliveConn.Close()
	
	fmt.Println("Creating new Main")
	cmd := exec.Command("mate-terminal", "-x", "go", "run", "main.go")
	cmd.Run()

	fmt.Println("Count == ", count)
	for i := 0; i<500; i++{
		go SendToNetwork(countPort,count)
		time.Sleep(time.Millisecond*1)
	}
	
}