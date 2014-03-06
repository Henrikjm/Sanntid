package network

import (
	"fmt"
	"net"
	"time"
	"strings"
	"encoding/json"
	"queue"
)

const N_ELEVATORS int = 2

const (
	ORDERPORT string = "44001"
	ALIVEPORT string = "44002"
	COSTPORT string = "44003"
)

type AliveArray [][]string

type TestVariable struct{
	Kek int 
	Lol  string 
}

func CheckError(err error, errorMsg string) {
	if err != nil {
		fmt.Println("!!Error type: " + errorMsg,"!!")
	}
}

func ImAliveUDP() {
	fmt.Println("Establishing IAmAlive...")
	sendAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:"+ALIVEPORT)
	CheckError(err, "ERROR while resolving UDP addr")
	conn, err := net.DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")
	msg := "ImAlive!"
	for {
		time.Sleep(time.Millisecond * 300)
		conn.Write([]byte(msg))
	}
}

func RecieveAliveUDP(aliveChan *chan string){
	data := make([]byte, 1024)
	ownAddr := *GetLocalIp();
	conn := MakeListenerConn(ALIVEPORT)
	for {		
		_, addr, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")

		if (string(data) == "ImAlive!") && (addr.String() != ownAddr){
			*aliveChan <- addr.String()//add/update alive map
		}
	}

}

func UpdateAliveUDP(aliveChan *chan string, updateChan *chan map[string]time.Time) {
	aliveMap := make(map[string]time.Time)
	var lengthOfMap int = 0

	for {	
		select{
			case incomingIP := <-*aliveChan:
				aliveMap[incomingIP] = time.Now()
			case <-*updateForCostChan:
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

func AliveOrDead(){

}

func SendUDP(port string, msg string) {
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

func ListenToOrderUDP(conn *net.UDPConn, incoming *chan string) {
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

func RecieveOrder(orderChannel chan queue.Order) {
	conn := MakeListenerConn(ORDERPORT)
	data := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")

		var newOrder queue.Order
		json.Unmarshal(newOrder, &data)
		go func(order Order, orderChannel chan Order){
			orderChannel <- newOrder
		}(newOrder)

	}
	//motta ordre
	//videresend ordre
	//bekreft mottat ordre
}

func SendOrder(order Order){
	sendAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:"+ORDERPORT)
	CheckError(err, "ERROR while resolving UDP addr")
	conn, err := net.DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")

	for{
	orderB,_ := json.Marshal(order)
	conn.Write([]byte(order))

	if RecieveCost(){
		break
	}
	//Motta cost fra alle
	//Hvis ikke cost mottas fra alle, oppdater koblingsliste og send på nytt
	}

	
	conn.Close()
}

func SendCost() {
	
}

func RecieveCost() bool{
	conn := MakeListenerConn(COSTPORT)

	dummy := make(map[string]time.Time)
	updateForCostChan <- dummy
	aliveMap <- updateForCostChan

	costMap := make(map[string]int)

	//Les UDP i ett sekund eller til alle har levert cost rapport
	costInstance := make(queue.Cost)
	data := make([]byte, 1024)
	for {
		
		_,err := net.ReadFromUDP(data)
		json.Unmarshal(costInstance, &data)



	}
	//Lag et costMap med det som leses på UDP - Hvor mange ganger skal det leses over UDP? Lese kontinuerlig over en viss periode?
	//Det skal være like langt som IP-listen
	//Dersom vi har en cost for alle IP kan vi gå videre
	

	conn.close()
}

func confirmOrder() {
	
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

