package network

import (
	"fmt"
	"net"
	"time"
	"strings"
	"encoding/json"
	."types"
)

const (
	ORDERPORT string = "44001"
	ALIVEPORT string = "44002"
	COSTPORT string = "44003"
	ELEVATORPORT string = "44004"
	ORDERCONFIRMATIONPORT string = "44005"
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

		if (string(data) == "ImAlive!"){
			aliveChan <- addr.String()//add/update alive map
		}
	}

}

func UpdateAliveUDP(aliveChan chan string, updateFromAliveChan chan Change, requestAliveChan chan map[string]time.Time, updateForConfirmationChan chan map[string]time.Time, updateForCostChan chan map[string]time.Time) {
	go ImAliveUDP()
	go RecieveAliveUDP(aliveChan)

	aliveMap := make(map[string]time.Time)


	for {	
		select{
			case incomingIP := <-aliveChan:
				if _,ok := aliveMap[incomingIP]; ok {
					aliveMap[incomingIP] = time.Now()
				}else{
					aliveMap[incomingIP]=time.Now()
					updateFromAliveChan <- Change{"new",incomingIP}
				}
			case <-updateForCostChan:
				updateForCostChan<-aliveMap
			case <-requestAliveChan:
				requestAliveChan<-aliveMap
			case <-updateForConfirmationChan:
				updateForConfirmationChan <- aliveMap
			default:
				for ip, value := range aliveMap {//Iterate through alive-map and delete timed-out machines
					if time.Now().Sub(value) > 500000000 {
						delete(aliveMap, ip)
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

func RecieveOrderFromUDP(newOrdersChan chan Order, recieveCostChan chan map[string]Cost, updateForCostChan chan map[string]time.Time) { //Må beregne cost og sende ut
	conn := MakeListenerConn(ORDERPORT)
	data := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFromUDP(data) 	//motta ordre
		CheckError(err, "ERROR ReadFromUDP")
		
		var newOrder Order
		json.Unmarshal(data, &newOrder)
		
		for{
			newOrdersChan <- newOrder //videresend ordre til costevaluering
			if RecieveCost(newOrder , recieveCostChan, updateForCostChan){ //Vent til all cost er mottat og så send dette til kømodul
				break
			}
		}
	}
	
	
	//bekreft mottat ordre
}

func SendOrderToUDP(orderChannel chan Order, costChan chan map[string]Cost, updateForConfirmationChan chan map[string]time.Time){//IKKE FERDIG
	conn := MakeSenderConn(ORDERPORT)
	orderConfirmationChan := make(chan bool)
	for{
	order := <- orderChannel
	go RecieveOrderConfirmation(orderConfirmationChan, order, updateForConfirmationChan)
	for  /*Should consider adding a limitation to # of tries*/{
	orderB,_ := json.Marshal(order)
	conn.Write([]byte(orderB))

	if <-orderConfirmationChan{
		break
		}
	}
	}
}

func SendCost() {//IKKE FERDIG
	//Recieve order
	//Request cost evaluation
	//Send cost to UDP
}

func RecieveCost(order Order, recieveCostChan chan map[string]Cost, updateForCostChan chan map[string]time.Time) bool{//IKKE FERDIG
	conn := MakeListenerConn(COSTPORT)

	//Oppdaterer liste over heiser som er tilkoblet
	aliveMap := make(map[string]time.Time)


	updateForCostChan <- aliveMap
	aliveMap = <- updateForCostChan

	costMap := make(map[string]Cost)

	//Les UDP i 500ms eller til alle har levert cost rapport
	var costInstance Cost
	data := make([]byte, 1024)
	t0 := time.Now()

	for {
		
		conn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Millisecond))
		_,_,err := conn.ReadFromUDP(data)

		if err.Error() != ("read udp4 0.0.0.0:"+ COSTPORT +": i/o timeout"){
			CheckError(err, "ERROR!! while recieving cost")
		}

		json.Unmarshal(data, &costInstance	)

		if costInstance.Order == order{
			costMap[costInstance.Ip] = costInstance
		}

		if len(costMap) == len(aliveMap){
			
			recieveCostChan <- costMap
			
			return true
		}
		if time.Now().Sub(t0) > 500000000{
			break
		}
	}
	//Lag et costMap med det som leses på UDP - Hvor mange ganger skal det leses over UDP? Lese kontinuerlig over en viss periode?
	//Det skal være like langt som IP-listen
	//Dersom vi har en cost for alle IP kan vi gå videre
	conn.Close()
	return false
}

func MakeListenerConn(port string) *net.UDPConn{
	udpAddr, err := net.ResolveUDPAddr("udp4", ":"+port)
	CheckError(err, "ERROR while resolving UDPaddr for listen")
	fmt.Println("Establishing ListenToNetwork")
	conn, err := net.ListenUDP("udp4", udpAddr)
	fmt.Println("Listening on port ", udpAddr.String())
	CheckError(err, "Error while establishing listening connection")
	return conn
}

func MakeSenderConn(port string) *net.UDPConn{
	sendAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:"+port)
	CheckError(err, "ERROR while resolving UDP addr for sending")
	conn, err := net.DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")
	return conn

}

func RecieveOrderConfirmation(orderConfirmationChan chan bool, order Order, updateForConfirmationChan chan map[string]time.Time) bool{ //IKKE TESTET
	//INITIALISERINGER
	conn := MakeListenerConn(ORDERCONFIRMATIONPORT)
	data := make([]byte, 1024)
	confirmationMap := make(map[string]time.Time)
	aliveMap := make(map[string]time.Time)
	
	//HENTER INN ANTALL MASKINER
	updateForConfirmationChan<-aliveMap
	aliveMap = <- updateForConfirmationChan
	t0 := time.Now() //REFERANSETID
	
	//VENTER I 500MS PÅ SVAR FRA ALLE MASKINER
	for{
		_,addr,err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR!! RecieveOrderConfirmation")
		confirmationMap[addr.String()] = time.Now()

		if len(confirmationMap) == len(confirmationMap) {
			return true
		}
		if time.Now().Sub(t0) > 500000000{
			return false
		}
	}
}

func SendElevator(SOMEELEVATORCHAN chan Elevator){
	conn := MakeSenderConn(ELEVATORPORT)
	//data := make([]byte, 1024)

	//SEND REQUEST
	//SOMEELEVATORCHAN <- REQUEST
	//MOTTA SVAR
	elevator := <- SOMEELEVATORCHAN
	elevatorB,_ := json.Marshal(elevator)
	conn.Write([]byte(elevatorB))
	//TA INN ELEVATOR FRA QUEUE





}//IKKE FERDIG

func RecieveElevator(){
	//conn := MakeListenerConn(ELEVATORPORT)
	//data := make([]byte, 1024)


}//IKKE FERDIG

func NetworkHandler(localIpChan chan string, updateFromAliveChan chan Change){//IKKE FERDIG
	aliveChan := make(chan string)
	requestAliveChan := make(chan map[string]time.Time)
	updateForConfirmationChan := make(chan map[string]time.Time)
	updateForCostChan := make(chan map[string]time.Time)


	go LocalIpSender(localIpChan)
	go UpdateAliveUDP(aliveChan, updateFromAliveChan, requestAliveChan, updateForConfirmationChan, updateForCostChan)
	 
}

