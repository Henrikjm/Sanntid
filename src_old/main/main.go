package main

import(
	"driver"
	"network"
	"queue"
	."types"
	"fmt"
)


func main() {
/*
//NOTES::::::::::::
1. Concidering changing the channel updating the DRIVER to be only a Order type, not entire elevator
concerened that queue might overwrite the direction and lastfloor variables, not ideal.



//::::::::::::::::::::
*/
fmt.Println(MOVE_UP, MOVE_DOWN, MOVE_STOP)
exitChan := make(chan string)
//---NETWORK - QUEUE
//------- Update
receiveElevatorChan := make(chan Elevator)
updateNetworkChan := make(chan Elevator)
//-------- Orders
newOrderFromUDPChan := make(chan Order)
deadOrderToUDPChan := make(chan Order)
orderToNetworkChan := make(chan Order)
//-------- Costs
sendCostChan := make(chan Cost, 2)
receivedCostsChan := make(chan []Cost)
//-------- Change
changedElevatorChan := make(chan Change)
//-------- Get
localIpChan := make(chan string)

//---DRIVER - QUEUE
// ------- I/O
localOrderChan := make(chan Order,3)
// ------- Update
receiveDriverUpdateChan := make(chan Elevator,1)
updateDriverChan := make(chan Elevator)



	

recieveCostChan := make(chan map[string]Cost)
costChan := make(chan map[string]Cost)

 	
 	




go driver.ControlHandler(localOrderChan, updateDriverChan, receiveDriverUpdateChan)
go queue.QueueHandler(receiveElevatorChan, updateNetworkChan, newOrderFromUDPChan, deadOrderToUDPChan, sendCostChan, receivedCostsChan, 
	changedElevatorChan, localIpChan , localOrderChan, updateDriverChan, receiveDriverUpdateChan, orderToNetworkChan)
go network.NetworkHandler(localIpChan, changedElevatorChan, sendCostChan, newOrderFromUDPChan, recieveCostChan, orderToNetworkChan, deadOrderToUDPChan,
 costChan, updateNetworkChan, receiveElevatorChan)

<-exitChan
}

//CHANNEL OVERWIEV
//Network-queue channel interface:
//------- Update
// receiveElevatorChan - for receiving updates on the elevators status
// updateNetworkChan - for sending updates on local elevator status
//-------- Orders
// sendLocalOrderChan - every non INTERNAL order must be relayed to
// newOrderChan - First instance of a new order, gives an order for calculation of cost
// deadOrderChan - sends orders from dead elevator to network module (to be used as new orders)
//-------- Costs
// sendCostChan - For sending cost after receiving newOrder, will be made a map in network and sent to all machines
// receivedCostsChan - for receiving costs, identefy whether to apply change localy (if cost.ip is local)
//-------- Changes
// changedElevatorChan - dead or new elevator, will be handled by first elevator in list of elevators
//-------- Get
// localIpChan - sends request for local ip and waits to receive it

//Driver-queue channel interface:
// ------- I/O
// localOrdersChan - for channeling orders received on internal buttons
// ------- Update
// receiveDriverUpdateChan - for updating the local elevator
// UpdateDriverChan - channel for sending elevator to driver, for setting lights (and more?)

//Local channels: queue
// ------- Update
//TimedUpdateChanNetwork - channel that activates withing gorutined function with sleep every eks 100 milisec
//TimedUpdateChanDriver - for activating a driver update of elevator
