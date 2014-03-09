package driver
 
 import (
	"fmt"
	"time"
	."types"
)

func SetNewDirection(elevator *Elevator){
	fmt.Println("setting new direction")
	if elevator.OrderQueue[0].Floor == 0{
		elevator.Direction = MOVE_STOP
	}else if elevator.LastFloor < elevator.OrderQueue[0].Floor{
		elevator.Direction = MOVE_UP
	}else if elevator.LastFloor > elevator.OrderQueue[0].Floor{
		elevator.Direction = MOVE_DOWN
	}else{
		fmt.Println("WARNING from SetNewDirection: first in orderQueue == LastFloor")
	}
}

func ReachedFloorClearOrders(elevator *Elevator){
	jump := 1
	if (elevator.OrderQueue[0].Floor == elevator.OrderQueue[1].Floor) && 
	(elevator.OrderQueue[0].Orientation == ORDER_INTERNAL || elevator.OrderQueue[1].Orientation == ORDER_INTERNAL){
		jump = 2
	}
	for i:= 0; i < MAX_ORDERS-jump; i++{
		elevator.OrderQueue[i] = elevator.OrderQueue[i+jump]
	}
	for i:= 1; i < jump + 1; i++{
		elevator.OrderQueue[MAX_ORDERS-i] = Order{}
	}
}

func InitializeElevator(elevator *Elevator){
	//to find a starting point
	elevator.Direction = MOVE_STOP
	elevator.LastFloor = 0
	reachedFloor := ReadFloor()
	if reachedFloor == 0{
		motorChannel <- MOVE_UP
		for{
			reachedFloor = ReadFloor()
			if  reachedFloor > 0{

				motorChannel <- MOVE_STOP
				break
			}
		}
	}
	// moves on to objective
	elevator.LastFloor = reachedFloor
	fmt.Println(elevator) //checkpoiny
}



func ControlHandler(localOrdersChan chan Order, receiveQueueUpdateChan chan Elevator, updateQueueChan chan Elevator){
	
	//variables
	var(
		elevator Elevator
		reachedFloor int
		state string
	)

	motorChannel = make(chan MoveDir)
	receiveQueueUpdateChan = make(chan Elevator)
	updateQueueChan = make(chan Elevator)

	//Function calls
	IoInit()
	ClearAllLights()
	go MotorControl()
	go GetOrderButton(localOrdersChan)
	

	//initialize
	//testvariables
	OrderQueue := []Order{Order{1, ORDER_INTERNAL}, Order{1, ORDER_UP}, Order{2, ORDER_UP}, Order{2, ORDER_INTERNAL}, Order{3, ORDER_UP}, Order{3, ORDER_INTERNAL}, Order{4, ORDER_INTERNAL}, Order{4, ORDER_DOWN}, Order{3, ORDER_DOWN},Order{2,ORDER_DOWN}}	
	elevator.OrderQueue = OrderQueue
	state = "start"
	
	
	InitializeElevator(&elevator)
	SetNewDirection(&elevator)
	fmt.Println(elevator) //checkpoint
	motorChannel <- elevator.Direction
	fmt.Println(elevator.Direction)
	if elevator.OrderQueue[0].Floor == reachedFloor{
		ReachedFloorClearOrders(&elevator)
		fmt.Println("DELETED? ", elevator)
		time.Sleep(time.Second*3)
	}

	
	// State-Machine
	for {
		fmt.Println(elevator)
		fmt.Println("elevator dir = ", elevator.Direction, " elevator.OrderQueue[0] = ", elevator.OrderQueue[0].Floor)
		//elevator = <- receiveQueueUpdateChan
		reachedFloor = ReadFloor()
		switch{
			case reachedFloor > 0 && reachedFloor != elevator.LastFloor:
				fmt.Println("case reachedFloor")
				SetFloorIndicatorLight(reachedFloor)
				elevator.LastFloor = reachedFloor
				if elevator.OrderQueue[0].Floor == reachedFloor{
					motorChannel <- MOVE_STOP
					elevator.Direction = MOVE_STOP
					ReachedFloorClearOrders(&elevator)
					time.Sleep(time.Second*3)
					//updateQueueChan <- elevator
				}

			case elevator.Direction == MOVE_STOP && elevator.OrderQueue[0].Floor > 0:
				fmt.Println("case arrived")
				SetNewDirection(&elevator)
				//updateQueueChan <- elevator
				motorChannel <- elevator.Direction
			}
			//select{
		//case <- updateQueueChan:
		//	updateQueueChan <- elevator
		//default:
		//}
	}
}
