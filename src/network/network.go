package network

import (
	"fmt"
	"net"
	"time"
	"strings"
	"encoding/json"
<<<<<<< HEAD
	
=======
	".types"
>>>>>>> 02bce73356fd2c81974153a248e4bd35384fa1ea
)

const N_ELEVATORS int = 2

const (
	ORDERPORT string = "44001"
	ALIVEPORT string = "44002"
	COSTPORT string = "44003"
)


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

func RecieveAliveUDP(aliveChan chan string){
	data := make([]byte, 1024)
	conn := MakeListenerConn(ALIVEPORT)
	for {		
		_, addr, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")

		if (string(data) == "ImAlive!")){
			aliveChan <- addr.String()//add/update alive map
		}
	}

}

func UpdateAliveUDP(aliveChan chan string, updateFromAliveChan chan Change, requestAliveChan chan map[string]time.Time) {
	go ImAliveUDP()
	go RecieveAliveUDP(aliveChan)

	aliveMap := make(map[string]time.Time)
	var lengthOfMap int = 0

	for {	
		select{
			case incomingIP := <-aliveChan:
				if val,ok := aliveMap[incomingIP]; ok {
					aliveMap[incomingIP] = time.Now()
				}else{
					aliveMap[incomingIP]=time.Now()
					updateFromAliveChan <- Change{"new",incomingIP}
				}
			case <-updateForCostChan:
				updateForCostChan<-aliveMap
			case <-requestAliveChan:
				requestAliveChan<-aliveMap
			default:
				for ip, value := range aliveMap {//Iterate through alive-map and delete timed-out machines
					if time.Now().Sub(value) > 500000000 {
						delete(aliveMap, i)
						updateFromAliveChan <- Change{"dead", ip}
					}
				}
		}
	}
}


func SendUDP(port string, msg string) {
	sendAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:"+port)
	CheckError(err, "ERROR while resolving UDP addr")
	conn, err := net.DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")
	conn.Write([]byte(msg))
	conn.Close()
}

func GetLocalIp() string {
	conn, err := net.Dial("udp4", "google.com:80")
	CheckError(err, "ERROR: LocalIp: dialing to google.com:80")
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

func LocalIpSender(localIpChan chan string){
	for{
		<-localIpChan
		localIpChan<-GetLocalIp()
	}
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

func RecieveOrder(orderChannel chan Order) {
	conn := MakeListenerConn(ORDERPORT)
	data := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")

		var newOrder Order
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
	conn := MakeSenderConn(ORDERPORT)

	for  /*Should consider adding a limitation to # of tries*/{
	orderB,_ := json.Marshal(order)
	conn.Write([]byte(order))

	if RecieveCost(){
		break
		}
	}

	conn.Close()
}

func SendCost(cost Cost) {

}

func RecieveCost(order Order, chan costChan map[string]cost.Cost) bool{
	conn := MakeListenerConn(COSTPORT)

	//Oppdaterer liste over heiser som er tilkoblet
	dummy := make(map[string]time.Time)

	updateForCostChan <- dummy
	aliveMap <- updateForCostChan

	costMap := make(map[string]cost.Cost)

	//Les UDP i ett sekund eller til alle har levert cost rapport
	costInstance := make(Cost)
	data := make([]byte, 1024)
	t0 := time.Now()

	for {
		
		_,err := net.ReadFromUDP(data)
		json.Unmarshal(costInstance, &data)

		if costInstance.Order == order{
			costMap[costInstance.IP] = costInstance
		}

		if len(costMap) == len(aliveMap){
			
			costChan <- costMap
			
			return true
		}
		if time.Now().Sub(t0) > 500000000{
			return false
		}
	}
	//Lag et costMap med det som leses på UDP - Hvor mange ganger skal det leses over UDP? Lese kontinuerlig over en viss periode?
	//Det skal være like langt som IP-listen
	//Dersom vi har en cost for alle IP kan vi gå videre
	

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

func MakeSenderConn() *net.UDPConn{
	sendAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:"+ORDERPORT)
	CheckError(err, "ERROR while resolving UDP addr")
	conn, err := net.DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")
	return conn

}

func NetworkHandler(localIpChan chan string, updateFromAliveChan chan map[string]time.Time){
	aliveChan := make(chan string)
	requestAliveChan := make(chan map[string]time.Time)

	
	go localIpSender(localIpChan)
	go updateAliveUDP(aliveChan, updateFromAliveChan, requestAliveChan)
}

