package main

//import "net"
import "fmt"
//import "strings"
//import "strconv"
//import "driver"
//import "time"
//import "encoding/json"
//import "network"
//import ."types"




func main(){

	var lightArray [2][3]int
	fmt.Println(lightArray)
}
/*
	aliveChan := make(chan string)
	updateFromAliveChan := make(chan Change)
	requestAliveChan := make(chan map[string]time.Time)
	updateForConfirmationChan := make(chan map[string]time.Time)
	updateForCostChan := make(chan map[string]time.Time)

	requestMap := make(map[string]time.Time)

	go network.UpdateAliveUDP(aliveChan , updateFromAliveChan , requestAliveChan , updateForConfirmationChan , updateForCostChan )

	
	t0 := time.Now()
	for {
		select{
		case changeVariable := <- updateFromAliveChan:
			fmt.Println(changeVariable)
		default:
		}
		if time.Now().Sub(t0) > 2000000000{
			requestAliveChan<-requestMap
			fmt.Println(<-requestAliveChan)

			updateForConfirmationChan<-requestMap
			fmt.Println(<-updateForConfirmationChan)

			updateForCostChan<-requestMap
			fmt.Println(<-updateForCostChan)

			t0 = time.Now()
		}

		

	}
*/
/*
	
	localIpChan := make(chan string, 2)
	updateFromAliveChan := make(chan Change)
	sendCostChan := make(chan Cost, 2)
	newOrderChan := make(chan Order)
	recieveCostChan := make(chan map[string]Cost)
	orderChannel := make(chan Order)
 	costChan := make(chan map[string]Cost)
 	updateNetworkChan := make(chan Elevator)
 	receiveElevatorChan := make(chan Elevator)

 	fmt.Println("Starting...")
	go network.NetworkHandler(localIpChan, updateFromAliveChan, sendCostChan , newOrderChan, recieveCostChan, orderChannel,
 costChan, updateNetworkChan, receiveElevatorChan)

	newOrder := Order{1, 1}
	var change Change
	for{
		select{
		case <- recieveCostChan:
			fmt.Println("Loop complete?")

		case change = <-updateFromAliveChan:
			fmt.Println(change)
			orderChannel<-newOrder
		case  <- newOrderChan:
			fmt.Println("Sending cost")
			var ip string
			localIpChan <- ip
			sendCostChan <- Cost{1, Order{1,1}, <-localIpChan}
			fmt.Println("Cost sendt")
		
*/

		
		
		
	