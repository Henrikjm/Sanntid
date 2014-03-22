package main

//import "net"

//import "strings"
//import "strconv"
//import "driver"
//import "time"
//import "network"
//import ."types"

import (
   // "bufio"
    "fmt"
    //"io"
    "encoding/json"
    "io/ioutil"
    //"os"
)

func check(e error) {
    if e != nil {
        panic(e)
	   
    }
}

func checkForInternalOrders() []int{
	dat, err := ioutil.ReadFile("internalOrderBackupFile")
	
	var readOrders []int

	if err != nil {
		internalOrders := []int{0,0,0,0}
		d1,_ := json.Marshal(internalOrders)
	    err := ioutil.WriteFile("internalOrderBackupFile", d1, 0644)
	    fmt.Println("No internal orders stored. Writing to new file...")
	    check(err)

	    dat, err := ioutil.ReadFile("internalOrderBackupFile")
	  
		err = json.Unmarshal(dat, &readOrders)
	    fmt.Println("Made new file with empty orders: ", readOrders)
	}else{
		
		err = json.Unmarshal(dat, &readOrders)
		fmt.Println("Found file for internal orders. The orders were: ", readOrders)
	}
	if len(readOrders) != 4 {
		fmt.Println("CORRUPTED READ!!! Parsing empty order list to elevator.")
		return []int{0,0,0,0}
	}

	return readOrders
}

func main(){

	internal := checkForInternalOrders()
	fmt.Println("The check gave: ", internal)

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

		
		
		
	