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
	ORDERPORT string = "42035"
	ALIVEPORT string = "42012"
	COSTPORT string = "42013"
	ELEVATORPORT string = "42064"
	ORDERCONFIRMATIONPORT string = "42075"
	COSTCONFIRMATIONPORT string = "42076"
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
		time.Sleep(time.Millisecond*5)
		conn.Write([]byte(msg))
	}
}


func RecieveAliveUDP(aliveChan chan string){
	data := make([]byte, 1024)
	conn := MakeListenerConn(ALIVEPORT, "RecieveAliveUDP")
	for {		
		_, addr, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")
		ip := strings.Trim(strings.SplitAfter(addr.String(), ":")[0], ":")
		aliveChan <- ip
	}

}


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
					changedElevatorChan <- Change{"new",incomingIP}
				}
			case <-updateForCostChan:
				updateForCostChan<-aliveMap
			case <-requestAliveChan:
				requestAliveChan<-aliveMap
			case <-updateForConfirmationChan:
				updateForConfirmationChan <- aliveMap
			default:
				time.Sleep(time.Millisecond * 1)
				for ip, value := range aliveMap {
					if time.Now().Sub(value) > 1000000000 {
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
		time.Sleep(time.Millisecond * 1)
		<-localIpChan
		localIpChan<-ip
	}
}


func RecieveOrderFromUDP(newOrderFromUDPChan chan Order, recieveCostChan chan map[string]Cost, updateForCostChan chan map[string]time.Time) { //Må beregne cost og sende ut
	conn := MakeListenerConn(ORDERPORT, "RecieveOrderFromUDP")
	sender := MakeSenderConn(ORDERCONFIRMATIONPORT)
	data := make([]byte, 1024)
	costConn := MakeListenerConn(COSTPORT, "RecieveCost")

	for {
		time.Sleep(time.Millisecond * 1)
		n, _, err := conn.ReadFromUDP(data) 	

		CheckError(err, "ERROR ReadFromUDP")
		sender.Write([]byte("OrderRecieved"))
		var newOrder Order
		json.Unmarshal(data[:n], &newOrder)
		
		fmt.Println("Recieved order over UDP:      ", newOrder)
			go func (newOrder Order, recieveCostChan chan map[string]Cost, updateForCostChan chan map[string]time.Time, conn *net.UDPConn){ //Må hente cost "indiciduelt..."
				newOrderFromUDPChan <- newOrder //videresend ordre til costevaluering
				if RecieveCost(newOrder , recieveCostChan, updateForCostChan, costConn){ //Vent til all cost er mottat og så send dette til kømodul
					fmt.Println("All cost evaluations recieved for order:", newOrder)
				}
			}(newOrder, recieveCostChan, updateForCostChan, costConn)
	}
}


func SendOrderToUDP(orderToNetworkChan chan Order, deadOrderToUDPChan chan Order, costChan chan map[string]Cost, updateForConfirmationChan chan map[string]time.Time){//IKKE FERDIG
	conn := MakeSenderConn(ORDERPORT)
	var order Order
	for{
		select{
			case order = <- orderToNetworkChan: 

			case order = <- deadOrderToUDPChan:
		}
		
		orderB,_ := json.Marshal(order)
		conn.Write([]byte(orderB))	
		time.Sleep(time.Millisecond * 25)
		}
}

func SendCost(sendCostChan chan Cost) {
	for{
		time.Sleep(time.Millisecond * 1)
		sender := MakeSenderConn(COSTPORT)
		cost := <- sendCostChan
		
		costB,_ := json.Marshal(cost)
		fmt.Println("Sending cost evaluation to UDP. Cost is: ", cost)
		go func(sender *net.UDPConn, costB []byte){
			fmt.Println("Started sending cost")
			for i := 0; i < 10; i++ {
				time.Sleep(time.Millisecond * 10)
				sender.Write(costB)
			}
		}(sender, costB)
	}
}


func RecieveCost(order Order, recieveCostChan chan map[string]Cost, updateForCostChan chan map[string]time.Time, conn *net.UDPConn) bool{ //Lag et map sortert på ordre og lytt til
	
	aliveMap := make(map[string]time.Time)
	costMap := make(map[string]Cost)
	data := make([]byte, 1024)
	var costInstance Cost


	updateForCostChan <- aliveMap
	aliveMap = <- updateForCostChan

	t0 := time.Now()
	fmt.Println("Listening for cost updates.")
	for {
		time.Sleep(time.Millisecond * 1)
		conn.SetReadDeadline(time.Now().Add(time.Duration(100) * time.Millisecond))
		n,_,err := conn.ReadFromUDP(data)

		if (err != nil) && (err.Error() != ("read udp4 0.0.0.0:"+ COSTPORT +": i/o timeout")) {
				CheckError(err, "ERROR!! while recieving cost")
		}

		json.Unmarshal(data[:n], &costInstance)
		if err == nil && order == costInstance.Order{
			if _,ok := costMap[costInstance.Ip]; !ok{
				fmt.Println("Order is: ", order, "costInstance.Order is:", costInstance.Order)
			}
			if _,ok := costMap[costInstance.Ip]; !ok && costInstance.Order == order{
				fmt.Println("Saving cost from: ", costInstance.Ip)
				costMap[costInstance.Ip] = costInstance
			}
			if len(costMap) == len(aliveMap){
				fmt.Println("Cost recived from everyone. Forwarding results.")
				recieveCostChan <- costMap
				return true
			}
		}


		if (time.Now().Sub(t0) > 1000000000){
			fmt.Println("Cost listen timed out.")
			break
		}
		
	}
	
	return false
}

func MakeListenerConn(port string, process string) *net.UDPConn{
	udpAddr, err := net.ResolveUDPAddr("udp4", ":"+port)
	CheckError(err, "ERROR while resolving UDPaddr for listen")
	conn, err := net.ListenUDP("udp4", udpAddr)
	CheckError(err, "Error while establishing listening connection for "+process)
	return conn
}

func MakeSenderConn(port string) *net.UDPConn{
	sendAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:"+port)
	CheckError(err, "ERROR while resolving UDP addr for sending")
	conn, err := net.DialUDP("udp4", nil, sendAddr)
	CheckError(err, "ERROR while dialing")
	return conn

}

func RecieveOrderConfirmation(order Order, orderConfirmationChan chan bool, updateForConfirmationChan chan map[string]time.Time){ //IKKE TESTET
	conn := MakeListenerConn(ORDERCONFIRMATIONPORT, "RecieveOrderConfirmation")
	data := make([]byte, 1024)
	confirmationMap := make(map[string]time.Time)
	aliveMap := make(map[string]time.Time)
	
	
	updateForConfirmationChan<-aliveMap
	aliveMap = <- updateForConfirmationChan
	t0 := time.Now() 
	
	
	for{
		time.Sleep(time.Millisecond * 1)
		_,addr,err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR!! RecieveOrderConfirmation")
		confirmationMap[addr.String()] = time.Now()

		if len(confirmationMap) == len(aliveMap) {
			conn.Close()
			orderConfirmationChan <- true
			break
		}
		if time.Now().Sub(t0) > 500000000{
			conn.Close()
			orderConfirmationChan <- false
			break
		}
	}
}

func SendElevator(updateNetworkChan chan Elevator){ 
	conn := MakeSenderConn(ELEVATORPORT)
	for{
		time.Sleep(time.Millisecond * 100)
		dataToSend := <- updateNetworkChan
		elevatorB,_ := json.Marshal(dataToSend)
		conn.Write([]byte(elevatorB))

	}
}
func RecieveElevator(receiveElevatorChan chan Elevator){ 
	conn := MakeListenerConn(ELEVATORPORT, "RecieveElevator")
	data := make([]byte, 1024)

	var newElevator Elevator
	
	for{
		time.Sleep(time.Millisecond * 1)
		n,_,err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR!! RecieveElevator")

		json.Unmarshal(data[:n], &newElevator)
		receiveElevatorChan <- newElevator
	}	
}

func NetworkHandler(localIpChan chan string, changedElevatorChan chan Change, sendCostChan chan Cost, 
	newOrderFromUDPChan chan Order, recieveCostChan chan map[string]Cost, orderToNetworkChan chan Order, 
	deadOrderToUDPChan chan Order, costChan chan map[string]Cost, updateNetworkChan chan Elevator,
	 receiveElevatorChan chan Elevator){


	fmt.Println("NetworkHandler Started...")
	exitChan :=make(chan string)
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

	
	<- exitChan
	 
}