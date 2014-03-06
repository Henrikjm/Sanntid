package main

import(
	"fmt"
)

type(
	MoveDir int
	OrderDir int
)

const(
	N_BUTTONS int = 3
	N_FLOORS  int = 4
	MAX_ORDERS int = 10
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
	Ip string
	//subject to change (will trigger select)
	OrderQueue []Order
	Direction MoveDir
	LastFloor int
}

type Order struct{
	Floor int
	Orientation OrderDir
}

type Cost struct{
	Cost int
	Order Order
	Ip string
}

type Change struct{
	Type string
	Ip string
}


func GetElevatorCost(elevator Elevator, order Order) int {
	score := 0
	//correct direction plays 4 points difference
	if elevator.Direction == MOVE_STOP{
		score += 0
	}else if ((elevator.Direction == MOVE_UP) && (order.Floor > elevator.LastFloor)) || ((elevator.Direction == MOVE_DOWN) && (order.Floor < elevator.LastFloor)){ //hvis bestilling er i riktig retning
		score -= 4
	}else{
		score += 4
	}
	// each order in queue before this order plays 1 point (NOTE: the internal and both the external orders play part consequently)
	score += GetNumberOfStopsBeforeOrder(elevator, order)
	return score
}

func GetNumberOfStopsBeforeOrder(elevator Elevator, order Order)int{
	placement := GetInsertOrderPlacement(elevator, order)
	stops := placement
	//fmt.Println("GetNumberOfStopsBeforeOrder: placement == ", placement)
	for j:= 0; j < placement; j++{ //Removing common objective orders from score
		//fmt.Println(elevator.OrderQueue[j].Floor == elevator.OrderQueue[j+1].Floor)
		//fmt.Println(elevator.OrderQueue[j].Orientation == ORDER_INTERNAL || elevator.OrderQueue[j+1].Orientation == ORDER_INTERNAL)
		if (elevator.OrderQueue[j].Floor == elevator.OrderQueue[j+1].Floor) && (elevator.OrderQueue[j].Orientation == ORDER_INTERNAL || elevator.OrderQueue[j+1].Orientation == ORDER_INTERNAL){
			j += 1
			stops -= 1
		}
	}
	return stops
}

func GetInsertOrderPlacement(elevator Elevator, order Order) int{
	priOrder := GetInsertOrderPriority(elevator, order)
	//fmt.Println("GetInserOrderPriority : order == ", order, "priOrder == ", priOrder)
	for i := 0; i < MAX_ORDERS; i++{
		if elevator.OrderQueue[i] == order{
			fmt.Println("ERROR in InsertOrder: identical order in queue")
			break
		}else if (GetInsertOrderPriority(elevator, elevator.OrderQueue[i]) >= priOrder) && (elevator.OrderQueue[i].Floor >= order.Floor) {
			return i
		}
	}
	return -1
}

func GetInsertOrderPriority(elevator Elevator, order Order) int{
		if order.Floor == 0{
			fmt.Println("WARNING GetInsertOrderPriority: order.Floor == 0")
			return 5
		}else if elevator.Direction == MOVE_UP{
			if order.Floor > elevator.LastFloor{
				if order.Orientation == ORDER_UP || order.Orientation == ORDER_INTERNAL{
					return 1
				}else if order.Orientation == ORDER_DOWN{
					return 2
				}
			}else if order.Floor <= elevator.LastFloor{
				if order.Orientation == ORDER_DOWN || order.Orientation == ORDER_INTERNAL{
					return 3
				}else if order.Orientation == ORDER_UP{
					return 4
				}
			}
		}else if elevator.Direction == MOVE_DOWN{
			if order.Floor < elevator.LastFloor{
				if order.Orientation == ORDER_DOWN || order.Orientation == ORDER_INTERNAL{
					return 1
				}else if order.Orientation == ORDER_UP{
					return 2
				}
			}else if order.Floor >= elevator.LastFloor{
				if order.Orientation == ORDER_UP || order.Orientation == ORDER_INTERNAL{
					return 3
				}else if order.Orientation == ORDER_DOWN{
					return 4
				}
			}
		}
		return -1
}

func InsertOrder(elevator Elevator, order Order){
	if order.Floor == 0{
		fmt.Println("ERROR in InsertOrder: order.Floor == 0")
		return
	}
	placement := GetInsertOrderPlacement(elevator, order)
	if placement == -1{
		fmt.Println("WARNING in InsertOrder: order existing, insertion cancelled")
	}
	var temp, insert Order
	insert = order
	for i := placement; i <MAX_ORDERS; i++{
		temp = elevator.OrderQueue[i]
		elevator.OrderQueue[i] = insert
		insert = temp
	}
}

func GetLocalElevatorIndex(elevators []Elevator, localIpChan chan string)int{
	localIpChan <- "LocalIp"
	localIp := <- localIpChan
	for i := 0; i < N_ELEVATORS; i++{
		if elevators[i].Ip == localIp{
			fmt.Println("LocalElevatorIndex of ", elevators[i].Ip, " = ", i)
			return i
		}
	}
	fmt.Println("ERROR local Ip not found i elevators")
	return -1
}

func HandleDeadElev(elevators []Elevator, ip string, newOrdersChan chan Order){
	var i int
	var deadElevQueue []Order
	for i = 0 ; i < N_ELEVATORS; i++{
		if elevators[i].Ip == ip{
			deadElevQueue = elevators[i].OrderQueue
			break
		}
	}
	for i = 0; i < len(deadElevQueue); i++{
		newOrdersChan <- deadElevQueue[i]
	}
}



func queueHandler(elevatorsChan chan []Elevator, newOrdersChan chan Order, localOrdersChan chan Order, receivedCostsChan chan []Cost, elevatorsChangeChan chan Change, localIpChan chan string, sendCostChan chan Cost){
	
	//Making situation picture
	elevators := <- elevatorsChan
	localElevatorIndex := GetLocalElevatorIndex(elevators, localIpChan)

	//Variables
	var newOrder, localOrder Order
	var localCost, receivedCost Cost
	var receivedCosts []Cost

//Network channel interface:
// elevatorsChan - for sending and receiving updates on the elevators status
// newOrdersChan - First instance of a new order, gives an order for calculation of cost
// sendCostChan - For sending cost after receiving newOrder, will be made a map in network and sent to all machines
// receivedCostsChan - for receiving costs, identefy whether to apply change localy (if cost.ip is local)

//Driver channel interface:
// localOrdersChan - for channeling orders received on internal buttons


	//Listening and handling
	for{
		select{
		case elevators = <-elevatorsChan: //recieves update for local queue
		case localOrder = <- localOrdersChan: //recieves local orders from driver, imedeatly insert localy and send update
			InsertOrder(elevators[localElevatorIndex], localOrder)
			elevatorsChan <- elevators
		case newOrder = <-newOrdersChan: //receives new order and replies with sending local Cost
			localCost = Cost{GetElevatorCost(elevators[localElevatorIndex], newOrder), newOrder, elevators[localElevatorIndex].Ip}
			sendCostChan <- localCost
		case receivedCosts = <- receivedCostsChan: //receives a map of costs and ip's
			best := Cost{}
			best.Cost = 20
			for _, receivedCost =  range receivedCosts{
				if receivedCost.Cost < best.Cost{
					best = receivedCost
				}
			}
			if best.Ip == elevators[localElevatorIndex].Ip{
				InsertOrder(elevators[localElevatorIndex], best.Order)
				elevatorsChan <- elevators
			}
		}
	}
}

func main() {
	
	fmt.Println("debugging")

	//Initiating test elevators / orders
	exOrder1 := make([]Order,MAX_ORDERS)
	exOrder1[0] = Order{3, ORDER_UP}
	exOrder1[1] = Order{3, ORDER_INTERNAL}
	exOrder1[2] = Order{4, ORDER_DOWN}
	nuOrder1 := Order{2,ORDER_UP}
	nuOrder2 := Order{4,ORDER_INTERNAL}
	nuOrder3 := Order{2.ORDER_DOWN}
	nuOrder4 := Order{1,ORDER_UP}
	elevator1:= Elevator{"elevator1:", exOrder1, MOVE_UP, 1}

	//Testing Inserting
	fmt.Println(elevator1, " s.t Order = ", nuOrder1)
	InsertOrder(elevator1, nuOrder1)
	fmt.Println(elevator1)
	fmt.Println("s.t Order = ", nuOrder3)
	InsertOrder(elevator1, nuOrder3)
	fmt.Println("s-t Order = ", nuOrder4)
	InsertOrder(elevator1, nuOrder4)
	fmt.Println(elevator1)

	//Testing GetNumberOfStops
	fmt.Println("num stops for ", nuOrder2, "in elevator1 = ", GetNumberOfStopsBeforeOrder(elevator1, nuOrder2))

	//Testing scoreSys
	fmt.Println("Cost for elevator 1 and ",nuOrder2," = ", GetElevatorCost(elevator1, nuOrder2))


}