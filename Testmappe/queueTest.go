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
	elevIp string
	//subject to change (will trigger select)
	orderQueue []Order
	direction MoveDir
	lastFloor int
}

type Order struct{
	floor int
	orientation OrderDir
}


func GetElevatorEligibilityScore(elevator Elevator, order Order) int {
	score := 0
	//correct direction plays 4 points difference
	if elevator.direction == MOVE_STOP{
		score += 0
	}else if ((elevator.direction == MOVE_UP) && (order.floor > elevator.lastFloor)) || ((elevator.direction == MOVE_DOWN) && (order.floor < elevator.lastFloor)){ //hvis bestilling er i riktig retning
		score =+ 4
	}else{
		score =- 4
	}
	
	// each order in queue before this order plays 1 point (NOTE: the internal and both the external orders play part consequently)
	placement := 0
	priOrder := GetInsertOrderPri(elevator, order)
	for i := 0; i < MAX_ORDERS; i++{
		if elevator.orderQueue[i] == order{
			fmt.Println("!!Error in GetElevatorEligibilityScore: identical order in orderQueue")
			score -= 100
		}else if GetInsertOrderPri(elevator, elevator.orderQueue[i]) > priOrder{
			placement = i
			for j:= 0; j < i; j++{ //Removing common objective orders from score
				fmt.Println("forloop placement = ", placement)
				if (elevator.orderQueue[j].floor == elevator.orderQueue[j+1].floor) && (elevator.orderQueue[j].orientation == ORDER_INTERNAL || elevator.orderQueue[j+1].orientation == ORDER_INTERNAL){
					placement -= 1
					j += 1
				}
			}
			break
		}
	}
	fmt.Println("placement = ", placement)
	score += placement
	return score
}


func GetInsertOrderPri(elevator Elevator, order Order) int{
		if order.floor == 0{
			return 5
		}else if elevator.direction == MOVE_UP{
			if order.floor > elevator.lastFloor{
				if order.orientation == ORDER_UP || order.orientation == ORDER_INTERNAL{
					return 1
				}else if order.orientation == ORDER_DOWN{
					return 2
				}
			}else if order.floor <= elevator.lastFloor{
				if order.orientation == ORDER_DOWN || order.orientation == ORDER_INTERNAL{
					return 3
				}else if order.orientation == ORDER_UP{
					return 4
				}
			}

		}else if order.orientation == ORDER_DOWN{
			if order.floor < elevator.lastFloor{
				if order.orientation == ORDER_DOWN || order.orientation == ORDER_INTERNAL{
					return 1
				}else if order.orientation == ORDER_UP{
					return 2
				}
			}else if order.floor >= elevator.lastFloor{
				if order.orientation == ORDER_UP || order.orientation == ORDER_INTERNAL{
					return 3
				}else if order.orientation == ORDER_DOWN{
					return 4
				}
			}
		}
		return -1
}

func InsertOrder(elevator Elevator, order Order){
	if order.floor == 0{
		fmt.Println("!!Error in InsertOrder: order.floor == 0")
		return
	}

	var placement, priOrder int
	priOrder = GetInsertOrderPri(elevator, order)

	for i := 0; i < MAX_ORDERS; i++{
		if elevator.orderQueue[i] == order{
			fmt.Println("!!Error in InsertOrder: identical order in queue")
		}else if GetInsertOrderPri(elevator, elevator.orderQueue[i]) < priOrder{
			placement = i
			break
		}
	}
	
	var temp, insert Order
	insert = order

	for i := placement; i <MAX_ORDERS; i++{
		temp = elevator.orderQueue[i]
		elevator.orderQueue[i] = insert
		insert = temp
	}
}




func main(){
	exOrder := make([]Order,MAX_ORDERS)
	exOrder[0] = Order{3, ORDER_UP}
	exOrder[1] = Order{3, ORDER_INTERNAL}
	exOrder[2] = Order{4, ORDER_DOWN}

	nuOrder := Order{2,ORDER_UP}
	elevator := Elevator{"elevator1:", exOrder, MOVE_UP, 1}

	fmt.Println(elevator, " s.t Order = ", nuOrder)
	InsertOrder(elevator, nuOrder)
	fmt.Println(elevator)

	EligScore := GetElevatorEligibilityScore(elevator, nuOrder)
	fmt.Println("eligble score of elevator1 :", EligScore)

	EligScore = GetElevatorEligibilityScore(elevator, Order{4,ORDER_INTERNAL})
	fmt.Println("eligble score of elevator1 :", EligScore)

}