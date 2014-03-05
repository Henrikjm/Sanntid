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
	msg := "ImAlive!"
	for {
		time.Sleep(time.Millisecond * 300)
		conn.Write([]byte(msg))
	}
}

func RecieveAliveUDP(alivePort string, aliveChan *chan map[string]time.Time){
	data := make([]byte, 1024)
	ownAddr := *GetLocalIp();
	conn := network.MakeListenerConn(alivePort)
	for {		
		_, addr_, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")

		if (string(data) == "ImAlive!") && (addr.String() != ownAddr){
			*aliveChan <- addr.String()//add/update alive map
		}
	}

}

func UpdateAliveUDP(aliveChan *chan string, updateChan *chan map[string]time.Time) {
	for {
		select{
			case incomingIP := <-*aliveChan:
				aliveMap[incomingIP] = time.Now()
			case <-updateChan:
				*updateChan<-aliveMap
			default:
				for i, value := range aliveMap {//Iterate through alive-map and delete timed-out machines
					if time.Now().Sub(value) > 500000000 {
						delete(aliveMap, i)
					}
				}
				if lengthOfMap != len(aliveMap) {
					lengthOfMap = len(aliveMap)
					*updateChan <- aliveMap
				}
		}
	}
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

