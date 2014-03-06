package main

import(
	"fmt"
	"driver"
	"network"
	"queue"
)


func main() {
	
var(
//Queue channels
	/* very much subject of change
	recieveElevatorChan chan Elevator
	sendElevatorChan chan Elevator
	newOrderChan chan Order
	deadOrderChan chan Order
	sendCostChan chan Cost
	receivedCostsChan chan []Cost 
	changedElevatorChan chan Change
	localIpChan chan string
	updateDriverChan chan Elevator
	lreceiveDriveUpdateChan chan Orderr
*/
)

}

//CHANNEL OVERWIEV
//Network-queue channel interface:
//------- Update
// receiveElevatorChan - for receiving updates on the elevators status
// updateNetworkChan - for sending updates on local elevator status
//-------- Orders
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
