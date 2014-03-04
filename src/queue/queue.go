package main

import(
	"fmt"
	//"encoding/json"
)

type(
	MoveDir int
	OrderDir int
)


const(
	N_BUTTONS int = 3
	N_FLOORS  int = 4
	MAX_ORDERS int = 11
	N_ELEVATORS int = 2

	ORDER_UP OrderDir = iota //matched with FLOOR for actuall order
	ORDER_DOWN
	ORDER_INTERNAL

	MOVE_UP MoveDir = iota //defines av elevators direction 
	MOVE_DOWN
	MOVE_STOP
)



type Elevator struct{
	//constant
	elevIp string
	//subject to change (will trigger select)
	workQueue []int
	direction MoveDir
	lastFloor int
}

type Order{
	floor int
	orientation OrderDir
}

func GetElevNumber(elevators []Elevator)(int){
	ownIp := network.GetLocalIp()
	for i := 0; i < N_ELEVATORS; i++{
		if elevators[i].elevIp == ownIp{
			return i
		}
	}
}

//Compare order returns a value based on place in priority (& thus place in workQueue)
//ASSUMES: that workQueue[0] is the current optimal priority
func CompareOrder(workQueue []int, lastFloor int, order Order)(int){
	if order.floor == 0 {
		Println("!!Error in CompareOrder ordered floor == 0")
		return nil
	}
	for i := 0; i < N_FLOORS; i++{
		
		if workQueue[i] == 0 {
			return i

		}else if order.direction == MOVE_UP{
			if order.floor > lastFloor && workQueue[i] > order.floor{
				return i
			}else if order.floor < lastFloor && workQueue[i] < order.floor {
				return i
			}

		}else if order.direction == MOVE_DOWN{
			if order.floor < lastFloor && workQueue[i] < order.floor{
				return i
			}else if order.floor > last.floor && workQueue[i] > order.floor{
				return i
			}
		}
	}



}

func InsertOrder(elevators []Elevator)[]Elevator{
}

//Update Localy
func HandleInternalOrder(elevators []Elevator, order Order){}
	//if (order.floor > elevators[elevNumber].lastFloor) && (elevators[elevNumber].direction == MOVE_UP)


//Update Globaly
func HandleEksternalOrder(elevators []Elevator, order Order ){}

//Update Globaly (will all do this?)
func HandleDeadElev(elevators []Elevator, deadIp string){}



func main(){
	emptyWorkQueue := make([]int, N_FLOORS)
	var elevators [N_ELEVATORS]Elevator
	var internalOrderChan, eksternalOrderChan chan Order
	var internalOrder, eksternalOrder Order
	var elevStatusChan chan Elevator
	var dummyElevator Elevator
	var aliveChan chan string
	var deadIp, queuePort string
	var elevNumber int

	go driver.MonitorIO(internalOrderChan)
	go driver.MonitorIO(eksternalOrderChan)
	go driver.MonitorElev(elevStatusChan) //feeds elevStatusChan new signals from elevator (like floors((not defined)
	go network.MonitorAlive(aliveChan) //feeds aliveChan the dead IP (udefinert)

	elevNumber = GetElevNumber(elevators)

	for{
		select{
		case internalOrder = <- internalOrderChan:
			network.SendToNetworkUDP(queuePort, "Handling new internal order") //the common picture must be updated regardless
			HandleInternalOrder(internalOrder)
		case eksternalOrder = <- eksternalOrderChan:
			network.SendToNetworkUDP(queuePort, "Handling new eksternal order") //the common picture must be updated regardless
			HandleEksternalOrder(eksternalOrder)
		case dummyElevator = <- elevStatusChan: //making sure we know what our elevator is doing
			elevators[elevNumber].direction = dummyElevator.direction
			elevators[elevNumber].lastFloor = dummyElevator.lastFloor
		case deadIp = <- aliveChan:
			HandleDeadElev(elevators, deadIp)
		default: //for å ikke okkupere programtelleren
		}
	}


	elevators[0] = Elevator{"129.241.187.147", MOVE_STOP, emptyWorkQueue}
	fmt.Println(elevators[0])

	//elevators[1] = Elevator{"129.241.187.147", MOVE_STOP, [11]int}

	

}