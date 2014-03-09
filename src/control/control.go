package driver
 
 import (
	"fmt"
	//"time"
	."types"
)
var motorChannel chan MoveDir
var readFloorChannel chan int

func ControlHandler(localOrdersChan chan Elevator, receiveQueueUpdateChan chan Elevator, updateQueueChan chan Elevator){
	
	//variables
	var(
		elevator Elevator
		reachedFloor int
	)

	motorChannel = make(chan MoveDir)
	receiveQueueUpdateChan = make(chan Elevator)
	updateQueueChan = make(chan Elevator)

	//Function calls
	IoInit()
	go ClearAllLights()
	go MotorControl()
	

	//initialize
	//testvariables
	OrderQueue := make([]Order,MAX_ORDERS)
	OrderQueue[0] = Order{1,ORDER_INTERNAL}
	OrderQueue[1] = Order{2,ORDER_INTERNAL}
	OrderQueue[2] = Order{3,ORDER_INTERNAL}
	OrderQueue[3] = Order{4,ORDER_INTERNAL}
	elevator.OrderQueue = OrderQueue
	elevator.LastFloor = 0

	//to find a starting point
	reachedFloor = ReadFloor()
	if reachedFloor == 0{
		motorChannel <- MOVE_DOWN
		for{
			if ReadFloor() > 0{
				motorChannel <- MOVE_STOP
				break
			}
		}
	}
	// moves on to objective
	elevator.LastFloor = reachedFloor
	fmt.Println(elevator) //checkpoiny
	SetNewDirection(&elevator)
	fmt.Println(elevator) //checkpoint
	motorChannel <- elevator.Direction
	fmt.Println(elevator.Direction)
	UpdateQueueChan <- elevator

	// State-Machine
	for {
		elevator = <-receiveUpdatesQueueChan
		fmt.Println(elevator)
		reachedFloor = ReadFloor()
		switch{
		case reachedFloor != elevator.LastFloor && reachedFloor > 0:
			SetFloorIndicatorLight(reachedFloor)
			elevator.LastFloor = reachedFloor	
			if elevator.OrderQueue[0].Floor == reachedFloor{
				motorChannel <- MOVE_STOP
				elevator.Direction = MOVE_STOP
				ReachedFloorClearOrders(&elevator)
				UpdateQueueChan <- elevator
			}
		case elevator.Direction == MOVE_STOP && elevator.OrderQueue[0].Floor > 0:
			SetNewDirection(&elevator)
			UpdateQueueChan <- elevator
			motorChannel <- elevator.Direction
		}
	}
}
