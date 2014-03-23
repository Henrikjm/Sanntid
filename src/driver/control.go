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
		//elevator.Direction = MOVE_STOP //fmt.Println("WARNING from SetNewDirection: first in orderQueue == LastFloor")
	}
	return elevator
}

func ReachedFloorClearOrders(elevator Elevator,  changeInElevator chan bool) Elevator{//, floor int){
	var temp1 []Order
	var temp2 []Order
	var result []Order
	//var oppositeDir OrderDir
	
	currentFloor := ReadFloor()


	if len(elevator.OrderQueue) < 1{
		for len(elevator.OrderQueue)<MAX_ORDERS {
			elevator.OrderQueue = append(elevator.OrderQueue, Order{0,0})
		}
	}


	if currentFloor == elevator.OrderQueue[0].Floor{
		//Fjerner alle ordre lik den som håndteres nå.
		for _, orderInstance := range elevator.OrderQueue{
				if orderInstance != elevator.OrderQueue[0] && orderInstance.Floor != 0{
					temp1 = append(temp1, orderInstance)
				}
		}
		//Fjerner alle instanser med interne ordre på samme etasje
		currentOrder := Order{currentFloor, ORDER_INTERNAL}
		for _, orderInstance := range temp1{
				if orderInstance != currentOrder && orderInstance.Floor != 0{
					temp2 = append(temp2, orderInstance)
				}
		}

		//Fjerner alle eksterne ordre på samme etasje dersom dette er den neste ordren.
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

	/*jump := 1
	if (elevator.OrderQueue[0].Floor == elevator.OrderQueue[1].Floor) && 
		(elevator.OrderQueue[0].Orientation == ORDER_INTERNAL || 
		elevator.OrderQueue[1].Orientation == ORDER_INTERNAL){
		jump = 2
	}
	for i:= 0; i < MAX_ORDERS-jump; i++{
		elevator.OrderQueue[i] = elevator.OrderQueue[i+jump]
	}
	for i:= 0; i < jump; i++{
		elevator.OrderQueue[MAX_ORDERS-1-i] = Order{}
	}
	return elevator*/

	


func UpdateFloor(updateFloorChan chan int){
	for{
		select{
			case <- updateFloorChan:
				updateFloorChan <- ReadFloor()
		}
	}
}

func UpdateLights(timedLightUpdate chan []Elevator) {
	//Lager et todim array over eksterne ordre!
	var lightArray [2][4]int
	var elevatorsForLight []Elevator


	for{
		
		elevatorsForLight = <- timedLightUpdate
		//fmt.Println(elevatorsForLight)
		for i := 0; i < len(lightArray); i++ {
			for j := 0; j < len(lightArray[0]); j++ {
				lightArray[i][j] = 0
			}
		}


		//Iterer gjennom alle heiser
		//fmt.Println(elevatorsForLight)
		for i := 0; i < len(elevatorsForLight); i++ {
			//fmt.Println(elevatorsForLight[i].OrderQueue)
			//Iterer gjennom alle ordre
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

func InitElev(localElevator Elevator)(Elevator, string){
	fmt.Println("Initiating elevator...")

	reachedFloor := ReadFloor()
	if reachedFloor == 0{
		motorChannel <- MOVE_UP
		fmt.Println("Searching for floor.")
	}
	
	for{
		time.Sleep(time.Millisecond*10)
		reachedFloor = ReadFloor()
		if  reachedFloor > 0{
			motorChannel <- MOVE_STOP
			localElevator.Direction = MOVE_STOP
			localElevator.LastFloor = reachedFloor
			break
		}
	}
	fmt.Println("Local elevator initiated in floor: ", localElevator.LastFloor)
	return localElevator, "idle"
}

func TimedUpdate(timedUpdateChan chan string){
	for{
		timedUpdateChan <- "Update"
		time.Sleep(time.Millisecond * 50)
	}
}

func ControlHandler(localOrderChan chan Order, updateDriverChan chan Elevator, receiveDriverUpdateChan chan Elevator, updateFloorChan chan int, timedLightUpdate chan []Elevator, localUpdateDriverChan chan Elevator, updateFromDriverChan chan Elevator, readyForUpdateChan chan bool){
	fmt.Println("ControlHandler started.")
	//variables
	var(
		reachedFloor int
		waitTime time.Time
	)


	//channels
	motorChannel = make(chan MoveDir)
	setOrderLightChannel = make(chan []Order)
	stopButtonChannel = make(chan bool)
	changeInElevator := make(chan bool, 1)
	timedUpdateChan := make(chan string)
	
	//Function calls
	IoInit()

	
	go ClearAllLights()
	go MotorControl()
	go GetOrderButton(localOrderChan)
	go SetOrderLights()
	go TimedUpdate(timedUpdateChan)
	go UpdateFloor(updateFloorChan)
	go UpdateLights(timedLightUpdate)

	

	//testvariables
	//OrderQueue := []Order{Order{1, ORDER_INTERNAL}, Order{1, ORDER_UP}, Order{2, ORDER_UP}, Order{2, ORDER_INTERNAL}, Order{3, ORDER_UP}, Order{3, ORDER_INTERNAL}, Order{4, ORDER_INTERNAL}, Order{4, ORDER_DOWN}, Order{3, ORDER_DOWN},Order{2,ORDER_DOWN}}	
	
	localElevator := <- updateDriverChan

	localElevator, state := InitElev(localElevator)
	oldstate := state
	for{
		time.Sleep(time.Millisecond * 1)

		select{
		case <-changeInElevator:
			updateFromDriverChan <- localElevator
    		localElevator = <- updateFromDriverChan

    	case localElevator = <- localUpdateDriverChan:	
				setOrderLightChannel <- localElevator.OrderQueue

		case <-timedUpdateChan:
			readyForUpdateChan <- true
			localElevator = <- updateDriverChan
		

		default: //STATE MACHINE!
			
			if oldstate != state{
				fmt.Println(oldstate, "--->",state)
				oldstate = state
			}

			switch state{
			//case "start":
				
			
			case "moving":
				
				if ReadFloor() > 0{ 
					state = "floor"
				}/*else if((localElevator.Direction == MOVE_DOWN) && (reachedFloor < localElevator.OrderQueue[0].Floor)){
					localElevator = SetNewDirection(localElevator)
				}else if ((localElevator.Direction == MOVE_UP) && (reachedFloor > localElevator.OrderQueue[0].Floor)){
					localElevator = SetNewDirection(localElevator)
				}*/

			case "floor":
				reachedFloor = ReadFloor()
				localElevator.LastFloor = reachedFloor

				if reachedFloor == localElevator.OrderQueue[0].Floor{
					motorChannel <- MOVE_STOP
					ClearOrderLight(Order{ReadFloor(),ORDER_INTERNAL})
					
					state = "arrived"

					fmt.Println("Erasing ", localElevator.OrderQueue[0] , "from floor ", reachedFloor, "STATE:" , state)
					
					localElevator = ReachedFloorClearOrders(localElevator, changeInElevator) //bør oppdatere queue
					
					

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
				reachedFloor = ReadFloor()
				localElevator.LastFloor = reachedFloor

				ClearOrderLight(Order{ReadFloor(),ORDER_INTERNAL})

				if len(localElevator.OrderQueue)>1{
					if localElevator.OrderQueue[0].Floor > 0{
						if localElevator.OrderQueue[0].Floor == reachedFloor{
							
							localElevator = ReachedFloorClearOrders(localElevator, changeInElevator) 
										
						}else{
							localElevator = SetNewDirection(localElevator)
							motorChannel <- localElevator.Direction
							state = "moving"
						}	
					}else if(ReadFloor() == 0){
						fmt.Println("ERROR!!! The elevator has stopped between floors.")
					}

				}else{
					fmt.Println("Corrupted localElevator!!!")
					localElevator = ReachedFloorClearOrders(localElevator,  changeInElevator)
				}
			}
		}
	}

}
