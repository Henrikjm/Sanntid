package main

import(
	"driver"
	"network"
	"queue"
	."types"
)


func main() {

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
recieveCostChan := make(chan map[string]Cost)
//-------- Change
changedElevatorChan := make(chan Change)
//-------- Get
localIpChan := make(chan string)

//---DRIVER - QUEUE
// ------- I/O
localOrderChan := make(chan Order)
// ------- Update
receiveDriverUpdateChan := make(chan Elevator,1)
updateDriverChan := make(chan Elevator)
updateFloorChan := make(chan int)
timedLightUpdate := make(chan []Elevator)
localUpdateDriverChan := make(chan Elevator)
updateFromDriverChan := make(chan Elevator)
readyForUpdateChan := make(chan bool)
costChan := make(chan map[string]Cost)
updateDirectionLastFloorChan := make(chan Elevator)

 	




go driver.ControlHandler(localOrderChan, updateDriverChan, receiveDriverUpdateChan, updateFloorChan, timedLightUpdate, localUpdateDriverChan, updateFromDriverChan, readyForUpdateChan, updateDirectionLastFloorChan)
go queue.QueueHandler(receiveElevatorChan, updateNetworkChan, newOrderFromUDPChan, deadOrderToUDPChan, sendCostChan, recieveCostChan, 
	changedElevatorChan, localIpChan , localOrderChan, updateDriverChan, receiveDriverUpdateChan, orderToNetworkChan, updateFloorChan, timedLightUpdate, localUpdateDriverChan, updateFromDriverChan, readyForUpdateChan, updateDirectionLastFloorChan)

go network.NetworkHandler(localIpChan, changedElevatorChan, sendCostChan, newOrderFromUDPChan, recieveCostChan, orderToNetworkChan, deadOrderToUDPChan,
 costChan, updateNetworkChan, receiveElevatorChan)

<-exitChan
}