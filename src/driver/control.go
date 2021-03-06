package driver
 
 import (
	"fmt"
	"time"
	."types"
	)

func SetNewDirection(elevator Elevator) Elevator{
	fmt.Println("Setting new direction.")
	if elevator.OrderQueue[0].Floor == 0{
		elevator.Direction = MOVE_STOP
	}else if elevator.LastFloor < elevator.OrderQueue[0].Floor{
		elevator.Direction = MOVE_UP
	}else if elevator.LastFloor > elevator.OrderQueue[0].Floor{
		elevator.Direction = MOVE_DOWN
	}else{
		elevator.Direction = MOVE_STOP
	}

	return elevator
}

func ReachedFloorClearOrders(elevator Elevator,  changeInElevator chan bool) Elevator{
	var temp1 []Order
	var temp2 []Order
	var result []Order

	
	currentFloor := ReadFloor()


	if len(elevator.OrderQueue) < 1{
		for len(elevator.OrderQueue)<MAX_ORDERS {
			elevator.OrderQueue = append(elevator.OrderQueue, Order{0,0})
		}
	}


	if currentFloor == elevator.OrderQueue[0].Floor{
		//Removes all orders identical to the one handled now.
		for _, orderInstance := range elevator.OrderQueue{
				if orderInstance != elevator.OrderQueue[0] && orderInstance.Floor != 0{
					temp1 = append(temp1, orderInstance)
				}
		}
		//Removes all orders which are internal orders on the same floor.
		currentOrder := Order{currentFloor, ORDER_INTERNAL}
		for _, orderInstance := range temp1{
				if orderInstance != currentOrder && orderInstance.Floor != 0{
					temp2 = append(temp2, orderInstance)
				}
		}
		//Removes all existing orders on the same floor if the next order is on the same floor.
		if elevator.OrderQueue[0].Floor == elevator.OrderQueue[1].Floor{
			for _, orderInstance := range temp2{
				if orderInstance != elevator.OrderQueue[1] && orderInstance.Floor != 0{
					result = append(result, orderInstance)
				}
			}
		}else{
			result = temp2
		}
	}
	for {
    	if len(result) < MAX_ORDERS{
    		result = append(result,Order{0,0})
    	}else{break}
    	if len(result) > MAX_ORDERS{
    		fmt.Println("ERROR IN QUEUE LENGTH")
    	}
    }
    elevator.OrderQueue = result
	changeInElevator <- true
	return elevator
}

func UpdateFloor(updateFloorChan chan int){
	for{
		select{
			case <- updateFloorChan:
				updateFloorChan <- ReadFloor()
		}
	}
}

func UpdateLights(timedLightUpdate chan []Elevator) {

	var lightArray [2][4]int
	var elevatorsForLight []Elevator


	for{
		
		elevatorsForLight = <- timedLightUpdate
		for i := 0; i < len(lightArray); i++ {
			for j := 0; j < len(lightArray[0]); j++ {
				lightArray[i][j] = 0
			}
		}

		for i := 0; i < len(elevatorsForLight); i++ {
			for j := 0; j < len(elevatorsForLight[i].OrderQueue); j++ {
				order := elevatorsForLight[i].OrderQueue[j]
				if (order.Orientation != ORDER_INTERNAL) && (order.Orientation != 0 && order.Floor != 0){
					if order.Orientation == ORDER_DOWN{
						lightArray[0][order.Floor-1] = 1 
					}else if order.Orientation == ORDER_UP{
						lightArray[1][order.Floor-1] = 1
					}
				}

			}
		}
		
		SetGlobalExternalLights(lightArray)
	}
}

func SetGlobalExternalLights(lightArray [2][4]int) {
	if lightArray[0][1] == 1{SetBit(LIGHT_DOWN2)}else{ClearBit(LIGHT_DOWN2)}
	if lightArray[0][2] == 1{SetBit(LIGHT_DOWN3)}else{ClearBit(LIGHT_DOWN3)}
	if lightArray[0][3] == 1{SetBit(LIGHT_DOWN4)}else{ClearBit(LIGHT_DOWN4)}
	if lightArray[1][0] == 1{SetBit(LIGHT_UP1)}else{ClearBit(LIGHT_UP1)}
	if lightArray[1][1] == 1{SetBit(LIGHT_UP2)}else{ClearBit(LIGHT_UP2)}
	if lightArray[1][2] == 1{SetBit(LIGHT_UP3)}else{ClearBit(LIGHT_UP3)}
}

func InitElev(localElevator Elevator, updateDirectionLastFloorChan chan Elevator)(Elevator, string){
	fmt.Println("Initiating elevator...")

	reachedFloor := ReadFloor()
	if reachedFloor == 0{
		motorChannel <- MOVE_UP
		fmt.Println("Searching for floor.")
	}
	
	for{
		time.Sleep(time.Millisecond*100)
		reachedFloor = ReadFloor()
		if  reachedFloor > 0{
			motorChannel <- MOVE_STOP
			localElevator.Direction = MOVE_STOP
			localElevator.LastFloor = reachedFloor
			go func(elevator Elevator){
						updateDirectionLastFloorChan <- elevator
					}(localElevator)
			break
		}
	}
	fmt.Println("Local elevator initiated in floor: ", localElevator.LastFloor)
	return localElevator, "idle"
}

func TimedUpdate(timedUpdateChan chan string){
	for{
		timedUpdateChan <- "Update"
		time.Sleep(time.Millisecond * 150)
	}
}

func ControlHandler(
	localOrderChan chan Order,
	updateDriverChan chan Elevator,
	receiveDriverUpdateChan chan Elevator,
	updateFloorChan chan int,
	timedLightUpdate chan []Elevator,
	localUpdateDriverChan chan Elevator,
	updateFromDriverChan chan Elevator,
	readyForUpdateChan chan bool,
	updateDirectionLastFloorChan chan Elevator){

	fmt.Println("ControlHandler started.")
	var(
		reachedFloor int
		waitTime time.Time
	)

	motorChannel = make(chan MoveDir)
	setInternalOrderLightChannel = make(chan []Order)
	stopButtonChannel = make(chan bool)
	changeInElevator := make(chan bool, 1)
	timedUpdateChan := make(chan string)
	
	IoInit()

	
	go ClearAllLights()
	go MotorControl()
	go GetOrderButton(localOrderChan)
	go SetInternalOrderLights()
	go TimedUpdate(timedUpdateChan)
	go UpdateFloor(updateFloorChan)
	go UpdateLights(timedLightUpdate)

	
	localElevator := <- updateDriverChan

	localElevator, state := InitElev(localElevator, updateDirectionLastFloorChan)
	reachedFloor = localElevator.LastFloor
	oldstate := state
	for{
		time.Sleep(time.Millisecond * 1)

		select{
		case <-changeInElevator:
			updateFromDriverChan <- localElevator
    		localElevator = <- updateFromDriverChan
    		setInternalOrderLightChannel <- localElevator.OrderQueue

    	case localElevator = <- localUpdateDriverChan:
    		setInternalOrderLightChannel <- localElevator.OrderQueue
		
		case <-timedUpdateChan:
			readyForUpdateChan <- true
			localElevator = <- updateDriverChan
			setInternalOrderLightChannel <- localElevator.OrderQueue

		default: 
			if oldstate != state{
				fmt.Println(oldstate,"-->",state)
				oldstate = state
			}
			switch state{
			case "moving":
				reachedFloor = ReadFloor()
				if reachedFloor > 0{ 
					state = "floor"
					
				}
			case "floor":
				
				if localElevator.LastFloor != reachedFloor {
					localElevator.LastFloor = reachedFloor
					go func(elevator Elevator){
						updateDirectionLastFloorChan <- elevator
					}(localElevator)
				}
							
				if reachedFloor == localElevator.OrderQueue[0].Floor{
					motorChannel <- MOVE_STOP
					ClearInternalOrderLight(Order{ReadFloor(),ORDER_INTERNAL})
					state = "arrived"
					fmt.Println("Erasing ", localElevator.OrderQueue[0] , "from floor ", reachedFloor, "STATE:" , state)
					localElevator = ReachedFloorClearOrders(localElevator, changeInElevator)
					waitTime = time.Now().Add(2*time.Second)
				}else if localElevator.Direction != MOVE_STOP{
					state = "moving"
				}else{
					state = "idle"

				}
				SetFloorIndicatorLight(reachedFloor)

			case "arrived":
				SetDoorOpenLight()
		        if time.Now().After(waitTime){
					state = "idle"
					ClearDoorOpenLight()
		        }
			case "idle":
				

				ClearInternalOrderLight(Order{ReadFloor(),ORDER_INTERNAL})
				if len(localElevator.OrderQueue)>1{
					if localElevator.OrderQueue[0].Floor > 0{
						if localElevator.OrderQueue[0].Floor == reachedFloor{
							localElevator = ReachedFloorClearOrders(localElevator, changeInElevator) 			
						}else{
							localElevator = SetNewDirection(localElevator)
							motorChannel <- localElevator.Direction
							go func(elevator Elevator){
								updateDirectionLastFloorChan <- elevator
							}(localElevator)
							state = "moving"
						}	
					}else if(ReadFloor() == 0){
						fmt.Println("ERROR!!! The elevator has stopped between floors.")
					}else {
						localElevator.Direction = MOVE_STOP
						go func(elevator Elevator){
						updateDirectionLastFloorChan <- elevator
						}(localElevator)						
					}

				}else{
					fmt.Println("Corrupted localElevator!!!")
					localElevator = ReachedFloorClearOrders(localElevator,  changeInElevator)
				}
			}
		}
	}

}
