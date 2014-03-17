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

func ReachedFloorClearOrders(elevator *Elevator){//, floor int){
	jump := 1
	if (elevator.OrderQueue[0].Floor == elevator.OrderQueue[1].Floor) && 
	(elevator.OrderQueue[0].Orientation == ORDER_INTERNAL || elevator.OrderQueue[1].Orientation == ORDER_INTERNAL){
		jump = 2
	}
	for i:= 0; i < jump; i++{
		clearOrderLightChannel <- elevator.OrderQueue[i]
	}
	for i:= 0; i < MAX_ORDERS-jump; i++{
		elevator.OrderQueue[i] = elevator.OrderQueue[i+jump]
	}
	for i:= 0; i < jump; i++{
		elevator.OrderQueue[MAX_ORDERS-1-i] = Order{}
		//clearOrderLightChannel <- elevator.OrderQueue[i]
	}
}

func TimedUpdate(timedUpdateChan chan string, interval int){
	for{
		timedUpdateChan <- "Update"
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
}




func ControlHandler(localOrderChan chan Order, updateDriverChan chan Elevator, updateQueueChan chan Elevator){
	fmt.Println("ControlHandler started.")
	//variables
	var(
		elevator Elevator
		reachedFloor int
		state string
		dummyElev Elevator
		waitTime time.Time
	)

	queueUpdateInterval := 50

	//channels
	motorChannel = make(chan MoveDir)
	setOrderLightChannel = make(chan []Order)
	clearOrderLightChannel = make(chan Order)
	stopButtonChannel = make(chan bool)
	timedUpdateQueueChan := make(chan string)

	//Function calls
	IoInit()
	go ClearAllLights()
	go MotorControl()
	go GetOrderButton(localOrderChan)
	go SetOrderLights()
	go ClearOrderLight()
	go TimedUpdate(timedUpdateQueueChan, queueUpdateInterval)
	go SetOrderLights()
	
	

	//testvariables
	//OrderQueue := []Order{Order{1, ORDER_INTERNAL}, Order{1, ORDER_UP}, Order{2, ORDER_UP}, Order{2, ORDER_INTERNAL}, Order{3, ORDER_UP}, Order{3, ORDER_INTERNAL}, Order{4, ORDER_INTERNAL}, Order{4, ORDER_DOWN}, Order{3, ORDER_DOWN},Order{2,ORDER_DOWN}}	
	
	dummyElev = <-updateDriverChan
	elevator.OrderQueue = dummyElev.OrderQueue
	state = "start"
	
	
	for{
		select{

		case <- timedUpdateQueueChan:
			
			//go func(elevator Elevator){ //litt dirty men går jo veldig bra
			updateQueueChan <- elevator
			//}(elevator)
		
		case dummyElev = <- updateDriverChan:
			
			elevator.OrderQueue = dummyElev.OrderQueue
			setOrderLightChannel <- elevator.OrderQueue
		
		default: //STATE MACHINE!
			
				switch state{
			case "start":
				
				elevator.Direction = MOVE_STOP
				elevator.LastFloor = 0
				reachedFloor = ReadFloor()
				if reachedFloor == 0{
					motorChannel <- MOVE_UP
				}
				for{
					time.Sleep(time.Millisecond*1)
					fmt.Println("Searching for floor.")
					reachedFloor = ReadFloor()
					if  reachedFloor > 0{
					motorChannel <- MOVE_STOP
					break
					}
				}
				elevator.LastFloor = reachedFloor
				setOrderLightChannel <- elevator.OrderQueue
				state = "idle"
				fmt.Println(state)
			case "moving":
				
				reachedFloor = ReadFloor()
				/*if ReadBit(STOP) {
					motorChannel <- MOVE_STOP 
					state = "stop"
					SetStopLight()
					time.Sleep(time.Millisecond*500)
				}else */if reachedFloor > 0{ 
					elevator.LastFloor = reachedFloor
					state = "floor"
				}
			case "floor":
				
				if elevator.LastFloor == elevator.OrderQueue[0].Floor{
					motorChannel <- MOVE_STOP
					elevator.Direction = MOVE_STOP
					state = "arrived"
					ReachedFloorClearOrders(&elevator)
					waitTime = time.Now().Add(2*time.Second)
				}else{
					state = "moving"
				}
				SetFloorIndicatorLight(reachedFloor)
			case "arrived":
				
		        if time.Now().After(waitTime){
					state = "idle"
		        }
			case "idle":
				
				clearOrderLightChannel <- Order{elevator.LastFloor,ORDER_INTERNAL}
				clearOrderLightChannel <- Order{elevator.LastFloor,ORDER_UP}
				clearOrderLightChannel <- Order{elevator.LastFloor,ORDER_DOWN}
				if elevator.OrderQueue[0].Floor > 0{
					if elevator.OrderQueue[0].Floor == elevator.LastFloor{
						ReachedFloorClearOrders(&elevator) //changing variable n
					}else{
						SetNewDirection(&elevator)
						fmt.Println("RAW TO MOTORCHANNEL: ", elevator.Direction)
						motorChannel <- elevator.Direction
						state = "moving"
						fmt.Println(state)
					}	
				}
			/*case "stop":
				if ReadBit(STOP){
					ClearStopLight()
					state = "moving"
					motorChannel <- elevator.Direction
					time.Sleep(time.Second*1)
				}*/
			}
		}
	}
}
