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
		
		time.Sleep(time.Millisecond*50)

		timedLightUpdate <- elevatorsForLight
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


func ControlHandler(localOrderChan chan Order, updateDriverChan chan Elevator, receiveDriverUpdateChan chan Elevator, updateFloorChan chan int, timedLightUpdate chan []Elevator, localUpdateDriverChan chan Elevator){
	fmt.Println("ControlHandler started.")
	//variables
	var(
		localElevator Elevator
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
	go UpdateFloor(updateFloorChan)
	go UpdateLights(timedLightUpdate)	

	//testvariables
	//OrderQueue := []Order{Order{1, ORDER_INTERNAL}, Order{1, ORDER_UP}, Order{2, ORDER_UP}, Order{2, ORDER_INTERNAL}, Order{3, ORDER_UP}, Order{3, ORDER_INTERNAL}, Order{4, ORDER_INTERNAL}, Order{4, ORDER_DOWN}, Order{3, ORDER_DOWN},Order{2,ORDER_DOWN}}	
	
	dummyElev = <- updateDriverChan
	localElevator = dummyElev
	state = "start"
	
	
	for{
		time.Sleep(time.Millisecond * 1)

		select{

		case localElevator = <- updateDriverChan:			
			setOrderLightChannel <- localElevator.OrderQueue

		case <- timedUpdateQueueChan:
			go func(receiveDriverUpdateChan chan Elevator){
			receiveDriverUpdateChan <- localElevator
			}(receiveDriverUpdateChan)
			

		case localElevator = <- localUpdateDriverChan:
				setOrderLightChannel <- localElevator.OrderQueue

		default: //STATE MACHINE!
			switch state{
			case "start":
				
				localElevator.Direction = MOVE_STOP
				localElevator.LastFloor = 0
				reachedFloor = ReadFloor()
				if reachedFloor == 0{
					motorChannel <- MOVE_UP
				}
				fmt.Println("Searching for floor.")
				for{
					time.Sleep(time.Millisecond*10)
					reachedFloor = ReadFloor()
					if  reachedFloor > 0{
					motorChannel <- MOVE_STOP
					break
					}
				}
				localElevator.LastFloor = reachedFloor
				setOrderLightChannel <- localElevator.OrderQueue
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
					localElevator.LastFloor = reachedFloor
					state = "floor"
				}
			case "floor":
				
				if localElevator.LastFloor == localElevator.OrderQueue[0].Floor{
					motorChannel <- MOVE_STOP
					localElevator.Direction = MOVE_STOP
					state = "arrived"
					ReachedFloorClearOrders(&localElevator)
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
				clearOrderLightChannel <- Order{localElevator.LastFloor,ORDER_INTERNAL}
				//clearOrderLightChannel <- Order{elevator.LastFloor,ORDER_UP}
				//clearOrderLightChannel <- Order{elevator.LastFloor,ORDER_DOWN}
				//fmt.Println(elevator)
				if localElevator.OrderQueue[0].Floor > 0{
					if localElevator.OrderQueue[0].Floor == localElevator.LastFloor{
						ReachedFloorClearOrders(&localElevator) //changing variable n
					}else{
						SetNewDirection(&localElevator)
						//fmt.Println("RAW TO MOTORCHANNEL: ", localElevator.Direction)
						motorChannel <- localElevator.Direction
						state = "moving"
						
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
