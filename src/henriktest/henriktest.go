package main

//import "net"
import "fmt"
import "time"
//import "encoding/json"
import "network"
import ."types"


func main(){

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


}