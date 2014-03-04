package network

import (
	"fmt"
	"net"
	"time"
	"strings"
)

const N_ELEVATORS int = 2

type AliveArray [][]string

func CheckError(err error, errorMsg string) {
	if err != nil {
		fmt.Println("!!Error type: " + errorMsg,"!!")
	}
}

func ImAliveUDP(port string) {
	fmt.Println("Establishing IAmAlive...")
	sendAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:"+port)
	CheckError(err, "ERROR while resolving UDP addr")
	conn, err := net.DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")
	msg := *GetLocalIp()
	for {
		time.Sleep(time.Millisecond * 300)
		conn.Write([]byte(msg))
	}
}

func RecieveAliveUDP(port string, aliveChan *chan string){
	data := make([]byte, 1024)
	
	for {
		_,addr,err := conn.ReadFromUDP(data)
		CheckError(err,"ERROR ReadFromUDP")
		aliveChan <- data
		
	}
}

func UpdateAliveUDP(aliveChan *chan string, newRequestChan *chan AliveArray){
	//Oppdaterer alivearray ihht. nye meldinger
	var workingArray AliveArray

	for{
		select{
			case newAlive := <- *aliveChan:
				//sjekk om newAlive er i array -> oppdater time, eller legg til ny med time
			case <- newRequestChan:
				newRequestChan <- 
			default:
				//sjekk etter døde maskiner
		}
	}
	//Fjerner døde maskiner etter x tid
	//Sender en array med alle koblingene dersom det blir forespurt
}

func SendToNetworkUDP(port string, msg string) {
	sendAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:"+port)
	CheckError(err, "ERROR while resolving UDP addr")
	conn, err := net.DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")
	conn.Write([]byte(msg))
	conn.Close()
}

func GetLocalIp() *string {
	conn, err := net.Dial("udp4", "google.com:80")
	CheckError(err, "ERROR: LocalIp: dialing to google.com:80")
	return &strings.Split(conn.LocalAddr().String(), ":")[0]
}

func ListenToNetworkUDP(conn *net.UDPConn, incoming *chan string) {
	data := make([]byte, 1024)
	//ownAddr := *GetLocalIp();
	for {
		_, /*addr*/_, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")
		//if addr.String() == ownAddr{ //OBS
			if err != nil{
				data = []byte("connection is dead")
			}
			//fmt.Println("ListenToNetwork: Channeling data :" + string(data))
			*incoming <- string(data)
		//}
	}
	conn.Close()
}



func MakeListenerConn(port string) *net.UDPConn{
	udpAddr, err := net.ResolveUDPAddr("udp4", ":"+port)
	CheckError(err, "ERROR while resolving UDPaddr for ListenToNetwork")
	fmt.Println("Establishing ListenToNetwork")
	conn, err := net.ListenUDP("udp4", udpAddr)
	fmt.Println("Listening on port ", udpAddr.String())
	CheckError(err, "Error while establishing listening connection")
	return conn
}

func NetworkHandler(){
	
}

