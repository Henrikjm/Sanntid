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
	ORDERPORT string = "44021"
	ALIVEPORT string = "44042"
	COSTPORT string = "44013"
	ELEVATORPORT string = "44064"
	ORDERCONFIRMATIONPORT string = "44075"
)


//OK
func CheckError(err error, errorMsg string) {
	if err != nil {
		fmt.Println("!!Error type: " + errorMsg,"!!")
	}
}

//OK
func ImAliveUDP() {
	fmt.Println("Establishing IAmAlive...")
	sendAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:"+ALIVEPORT)
	CheckError(err, "ERROR while resolving UDP addr")
	conn, err := net.DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")
	msg := "ImAlive!"
	for {
		time.Sleep(time.Millisecond*200)
		conn.Write([]byte(msg))
	}
}

//OK
func RecieveAliveUDP(aliveChan chan string){
	data := make([]byte, 1024)
	conn := MakeListenerConn(ALIVEPORT)
	for {		
		_, addr, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")
		ip := strings.Trim(strings.SplitAfter(addr.String(), ":")[0], ":") //Fjerner PORT og semikolon

		aliveChan <- ip //add/update alive map
		
	}

}

//OK
func UpdateAliveUDP(aliveChan chan string, changedElevatorChan chan Change, requestAliveChan chan map[string]time.Time, updateForConfirmationChan chan map[string]time.Time, updateForCostChan chan map[string]time.Time) {
	fmt.Println("UpdateAliveUDP Started")
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
					fmt.Println("|||||||||||||||||||||||||Forwarding:", incomingIP,"|||||||||||||||||||||||||||||||")
					changedElevatorChan <- Change{"new",incomingIP}
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
						changedElevatorChan <- Change{"dead", ip}
					}
				}
				
		}
	}
}


func GetLocalIp() string {
	conn, err := net.Dial("udp4", "google.com:80")
	CheckError(err, "ERROR: LocalIp: dialing to google.com:80")
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

func LocalIpSender(localIpChan chan string){
	ip := GetLocalIp()
	for{
		<-localIpChan
		localIpChan<-ip
	}
}


func RecieveOrderFromUDP(newOrderFromUDPChan chan Order, recieveCostChan chan map[string]Cost, updateForCostChan chan map[string]time.Time) { //Må beregne cost og sende ut
	conn := MakeListenerConn(ORDERPORT)
	sender := MakeSenderConn(ORDERCONFIRMATIONPORT)
	data := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFromUDP(data) 	//motta ordre
		CheckError(err, "ERROR ReadFromUDP")
		sender.Write([]byte("OrderRecieved"))
		var newOrder Order
		json.Unmarshal(data, &newOrder)
		
		fmt.Println("New order recieved. Waiting for cost evaluations.")
		for{
			newOrderFromUDPChan <- newOrder //videresend ordre til costevaluering
			if RecieveCost(newOrder , recieveCostChan, updateForCostChan){ //Vent til all cost er mottat og så send dette til kømodul
				fmt.Println("All cost evaluations recieved. Breaking loop.")
				break
			}
		}
	}
	
}


func SendOrderToUDP(orderChan chan Order, deadOrderToUDPChan chan Order, costChan chan map[string]Cost, updateForConfirmationChan chan map[string]time.Time){//IKKE FERDIG
	conn := MakeSenderConn(ORDERPORT)
	orderConfirmationChan := make(chan bool, 1)
	status := 1
	var order Order
	for{
		select{
			case order = <- orderChan: //Venter på ordre

			case order = <- deadOrderToUDPChan:
		}
		fmt.Println("Sending order. Waiting for RecieveOrderConfirmation.")
		go RecieveOrderConfirmation(order, orderConfirmationChan, updateForConfirmationChan)
		for  i := 0; i < 50; i++{
		orderB,_ := json.Marshal(order) //sender ut ordren til den har fått bekreftet at alle har mottatt
		conn.Write([]byte(orderB))

		if <- orderConfirmationChan{ //Sjekker bekreftelse, orderConfimartionChan er buffered med 1
			fmt.Println("OrderConfirmation received. Breaking loop.")
			status = 0
			break
			}

		}
		if status == 1{
			fmt.Println("ERROR!! SendOrder failed.")
		}
	}
}

func SendCost(sendCostChan chan Cost) {
	sender := MakeSenderConn(COSTPORT)
	for{
	fmt.Println("!!!!!!SendCost waiting...!!!!!")
	cost := <- sendCostChan
	fmt.Println("!!!!!SendCost recieved work...!!!!")
	costB,_ := json.Marshal(cost)
	fmt.Println("Sending cost evaluation to UDP.")
	sender.Write(costB)
	}
}

func RecieveCost(order Order, recieveCostChan chan map[string]Cost, updateForCostChan chan map[string]time.Time) bool{
	conn := MakeListenerConn(COSTPORT)
	aliveMap := make(map[string]time.Time)
	costMap := make(map[string]Cost)
	data := make([]byte, 1024)
	var costInstance Cost


	//Henter informasjon om hvor mange maskiner som er tilkoblet
	updateForCostChan <- aliveMap
	aliveMap = <- updateForCostChan


	//Lytter til UDP i 500ms eller til alle har levert "kostrapport". Dersom alle leverer kost vil den videresende et Map med IP og kost
	t0 := time.Now()
	fmt.Println("Listening for cost updates.")
	for {

		conn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Millisecond))
		_,_,err := conn.ReadFromUDP(data)

		if (err != nil) && (err.Error() != ("read udp4 0.0.0.0:"+ COSTPORT +": i/o timeout")) {
				CheckError(err, "ERROR!! while recieving cost")
		}
		json.Unmarshal(data, &costInstance)
		if costInstance.Order == order{ //Legger til
			costMap[costInstance.Ip] = costInstance
		}
		if len(costMap) == len(aliveMap){ //Sjekker om vi har fått svar fra alle
			recieveCostChan <- costMap
			fmt.Println("Cost recived from everyone. Forwarding results.")
			return true
		}
		if time.Now().Sub(t0) > 500000000{ //Sjekker om vi har brukt >500ms
			fmt.Println("Cost listen timed out.")
			break
		}
	}
	conn.Close()
	return false
}

func MakeListenerConn(port string) *net.UDPConn{
	udpAddr, err := net.ResolveUDPAddr("udp4", ":"+port)
	CheckError(err, "ERROR while resolving UDPaddr for listen")
	//fmt.Println("Establishing ListenToNetwork")
	conn, err := net.ListenUDP("udp4", udpAddr)
	//fmt.Println("Listening on port ", udpAddr.String())
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

func RecieveOrderConfirmation(order Order, orderConfirmationChan chan bool, updateForConfirmationChan chan map[string]time.Time) bool{ //IKKE TESTET
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

		if len(confirmationMap) == len(aliveMap) {
			orderConfirmationChan <- true
		}
		if time.Now().Sub(t0) > 500000000{
			orderConfirmationChan <- false
		}
	}
}

func SendElevator(updateNetworkChan chan Elevator){ //Ikke testet. Designet for å være goroutine.
	conn := MakeSenderConn(ELEVATORPORT)

	var localElevator Elevator

	for{
	localElevator = <- updateNetworkChan

	elevatorB,_ := json.Marshal(localElevator)
	conn.Write([]byte(elevatorB))
	}
}
func RecieveElevator(receiveElevatorChan chan Elevator){ //Ikke testet. Designet for å være goroutine.
	conn := MakeListenerConn(ELEVATORPORT)
	data := make([]byte, 1024)

	var newElevator Elevator
	
	for{
		_,_,err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR!! RecieveElevator")

		json.Unmarshal(data, &newElevator)
		receiveElevatorChan <- newElevator
	}	
}

func NetworkHandler(localIpChan chan string, changedElevatorChan chan Change, sendCostChan chan Cost, 
	newOrderFromUDPChan chan Order, recieveCostChan chan map[string]Cost, orderToNetworkChan chan Order, 
	deadOrderToUDPChan chan Order, costChan chan map[string]Cost, updateNetworkChan chan Elevator,
	 receiveElevatorChan chan Elevator){


	fmt.Println("NetworkHandler Started...")

	aliveChan := make(chan string)
	requestAliveChan := make(chan map[string]time.Time)
	updateForConfirmationChan := make(chan map[string]time.Time)
	updateForCostChan := make(chan map[string]time.Time)

	go SendCost(sendCostChan)
	go LocalIpSender(localIpChan)
	go UpdateAliveUDP(aliveChan, changedElevatorChan, requestAliveChan, updateForConfirmationChan, updateForCostChan)
	go SendOrderToUDP(orderToNetworkChan, deadOrderToUDPChan, costChan, updateForConfirmationChan)
	go RecieveOrderFromUDP(newOrderFromUDPChan, recieveCostChan, updateForCostChan)
	
	go SendElevator(updateNetworkChan)
	go RecieveElevator(receiveElevatorChan)

	for{time.Sleep(time.Millisecond * 200)
	}
	 
}