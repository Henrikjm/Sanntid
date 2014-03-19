package queue

import(
	"fmt"
	"time"
	."types"
	"strings"
	"strconv"
	"math"
)

func GetElevatorCost(elevator Elevator, order Order, updateFloorChan chan int) int {
	score := 0
	//correct direction plays 4 points difference
	if elevator.Direction == MOVE_STOP{
		score += GetNumberOfFloorsForOrder(updateFloorChan, order)

	}else if ((elevator.Direction == MOVE_UP) && (order.Floor > elevator.LastFloor)) || ((elevator.Direction == MOVE_DOWN) && (order.Floor < elevator.LastFloor)){ //hvis bestilling er i riktig retning
		score -= 4
	}else{
		score += 4
	}
	// each order in queue before this order plays 1 point (NOTE: the internal and both the external orders play part consequently)
	score += GetNumberOfStopsBeforeOrder(elevator, order)
	return score
}

func GetNumberOfFloorsForOrder(updateFloorChan chan int, order Order) int{
	updateFloorChan <- 0
	reachedFloor := <-updateFloorChan
	numberOfFloors := int(math.Abs(float64(reachedFloor - order.Floor)))
	return numberOfFloors
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
	var oldPri int
	if elevator.OrderQueue[0].Floor == 0{
		return 0
	}
	newPri := GetInsertOrderPriority(elevator, order)
	for i := 0; i < MAX_ORDERS; i++{
		oldPri = GetInsertOrderPriority(elevator, elevator.OrderQueue[i])
		fmt.Println("newPri = ", newPri, "oldPri = ", oldPri)
		if oldPri >= newPri{
				//newPri 1 and means that the order is "on-the-way" of the tour and reretour respectively, we optimize with respect to the current direction
			if newPri == oldPri{
				if (newPri == 1 || newPri == 4) && ((elevator.LastFloor < elevator.OrderQueue[0].Floor && elevator.OrderQueue[i].Floor > order.Floor) || (elevator.LastFloor > elevator.OrderQueue[0].Floor && elevator.OrderQueue[i].Floor < order.Floor)){
					return i
				//newPri 2 or 3 means that the order is "on-the-way" of the retour, we optimize with respect to direction
				}else if (newPri == 2 || newPri == 3) && ((elevator.LastFloor < elevator.OrderQueue[0].Floor && order.Floor > elevator.OrderQueue[i].Floor) || (elevator.LastFloor > elevator.OrderQueue[0].Floor && order.Floor < elevator.OrderQueue[i].Floor)){
					return i
				}
			}else{
				return i
			}
		}
	}
	return MAX_ORDERS-1
}

func GetInsertOrderPriority(elevator Elevator, order Order) int{
		if order.Floor == 0{
			fmt.Println("WARNING GetInsertOrderPriority: order.Floor == 0")
			return 5
		}else if elevator.OrderQueue[0].Floor > elevator.LastFloor{
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
		}else if elevator.OrderQueue[0].Floor < elevator.LastFloor{
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
	fmt.Println("|||||||||||||Inserting order||||||||||")
	if order.Floor == 0{
		fmt.Println("ERROR in InsertOrder: order.Floor == 0")
		return
	}
	placement := GetInsertOrderPlacement(elevator, order)
	fmt.Println("Placement of ", order, " = ", placement)
	/*
	if placement == -1{
		fmt.Println("WARNING in InsertOrder: order existing, insertion cancelled")
	}*/
	var temp, insert Order
	insert = order
	for i := placement; i <MAX_ORDERS; i++{
		temp = elevator.OrderQueue[i]
		elevator.OrderQueue[i] = insert
		insert = temp
	}
}

func GetLocalElevatorIndex(elevators []Elevator, localIp string)int{
	for i := 0; i < N_ELEVATORS; i++{
		if elevators[i].Ip == localIp{
			fmt.Println("LocalElevatorIndex of ", elevators[i].Ip, " = ", i)
			return i
		}
	}
	fmt.Println("ERROR local Ip not found i elevators")
	return -1
}

func HandleDeadElev(elevators []Elevator, ip string, deadOrderToUDPChan chan Order){
	var i int
	var deadElevQueue []Order
	for i = 0 ; i < N_ELEVATORS; i++{
		if elevators[i].Ip == ip{
			deadElevQueue = elevators[i].OrderQueue
			break
		}
	}
	for i = 0; i < len(deadElevQueue); i++{
		if deadElevQueue[i].Orientation != ORDER_INTERNAL && deadElevQueue[i].Orientation != 0{
			deadOrderToUDPChan <- deadElevQueue[i]
		}
	}
}

func HandleNewElevator(elevators []Elevator, ip string){
	for i := 0; i < N_ELEVATORS; i++{
		if elevators[i].Ip == ""{ //tom plass
			elevators[i].Ip = ip
			elevators[i].OrderQueue = make([]Order,MAX_ORDERS)
			break
		}
	}
}

func TimedUpdate(timedUpdateChan chan string){
	for{
		timedUpdateChan <- "Update"
		time.Sleep(time.Millisecond * 50)
	}
}

func IsNotInElevator(elevator Elevator, order Order) bool {
	for i :=0; i < MAX_ORDERS; i++{
		if elevator.OrderQueue[i] == order{
			return false
		}/*else if elevator.OrderQueue[i].Floor == 0{
				return true
		}*/
	}
	return true
}



func QueueHandler(receiveElevatorChan chan Elevator, updateNetworkChan chan Elevator, newOrderFromUDPChan chan Order, deadOrderToUDPChan chan Order, sendCostChan chan Cost, recieveCostChan chan map[string]Cost, 
	changedElevatorChan chan Change, localIpChan chan string, localOrderChan chan Order, updateDriverChan chan Elevator, receiveDriverUpdateChan chan Elevator, orderToNetworkChan chan Order,
	 updateFloorChan chan int, timedLightUpdate chan []Elevator, localUpdateDriverChan chan Elevator){
	
	fmt.Println("QueueHandler started.")
		
	//Variables
	var newOrder, localOrder Order
	var localCost, recievedCost Cost
	var receivedCostMap map[string]Cost
	var changedElevator Change
	var updateElevator Elevator
	timedUpdateChanNetwork := make(chan string)
	timedUpdateChanDriver := make(chan string)
	
	//OrderQueue := make([]Order, MAX_ORDERS) //OrderQueue := []Order{Order{1, ORDER_INTERNAL}, Order{1, ORDER_UP}, Order{2, ORDER_UP}, Order{2, ORDER_INTERNAL}, Order{3, ORDER_UP}, Order{3, ORDER_INTERNAL}, Order{4, ORDER_INTERNAL}, Order{4, ORDER_DOWN}, Order{3, ORDER_DOWN},Order{2,ORDER_DOWN}}	
	


	//testvars
	//timedUpdateChanDriver := make(chan string)
	//IP := "some IP"
	//elevator := Elevator{IP, OrderQueue, MOVE_STOP, 0}
	//elevators := []Elevator{elevator}
	//localElevatorIndex := 0



	//Making situation picture
	elevators := make([]Elevator, N_ELEVATORS) //empty list of elevators

	
	localIpChan <- "LocalIp"
	localIp := <- localIpChan //Gets the local IP
	
	HandleNewElevator(elevators, localIp) //Ads the Ip to empty slot of elevators
	localElevatorIndex := GetLocalElevatorIndex(elevators, localIp)
	receiveDriverUpdateChan <- elevators[0]
	updateElevator = <- receiveDriverUpdateChan //Ads information from elevator (driver)
	elevators[localElevatorIndex] = updateElevator

	go TimedUpdate(timedUpdateChanNetwork)
	go TimedUpdate(timedUpdateChanDriver)

	//Listening and handling
	fmt.Println("QueueHandler initiated.")
	for{
		time.Sleep(time.Millisecond * 10)
		
		select{
		// RULING OUT CHANNEL WAITING FOR NOW
		//-------------------------------------
		//receiving updates from other modules
		case localOrder = <- localOrderChan: //recieves local orders from driver, imedeatly insert localy and send update
			fmt.Println("localOrders")
			if IsNotInElevator(elevators[localElevatorIndex], localOrder){
				if localOrder.Orientation == ORDER_INTERNAL{
					InsertOrder(elevators[localElevatorIndex], localOrder)
					localUpdateDriverChan <- elevators[localElevatorIndex] //Bør legge inn localUpdate i control
				}else{
					orderToNetworkChan <- localOrder
				}
			}		

		case newOrder = <-newOrderFromUDPChan: //receives new order and replies with sending local Cost
			fmt.Println("RecievedNewOrder")
			localCost = Cost{GetElevatorCost(elevators[localElevatorIndex], newOrder,  updateFloorChan), newOrder, elevators[localElevatorIndex].Ip}
			sendCostChan <- localCost


		case receivedCostMap = <- recieveCostChan: //receives a map of costs and ip's
			fmt.Println("RecievedCost")
			best := Cost{}
			best.Cost = 20
			for _, recievedCost =  range receivedCostMap{
				if recievedCost.Cost < best.Cost{
					best = recievedCost
				}
			}

			//dummyIpStr:= strings.Trim(strings.SplitAfter(elevators[localElevatorIndex].Ip, "187")[1], ".")
			highestIp := 0 //strconv.Atoi(dummyIpStr)
			newIp := 0
			for _,recievedCost = range receivedCostMap{
				if recievedCost.Cost == best.Cost{
					dummyIpStr := strings.Trim(strings.SplitAfter(recievedCost.Ip, "187")[1], ".")
					newIp, _ = strconv.Atoi(dummyIpStr)
					fmt.Println(newIp, highestIp)
					if newIp > highestIp{
						highestIp = newIp
						best = recievedCost
					}
				}
			}
			fmt.Println("The best IP is: ", best.Ip)
			if best.Ip == elevators[localElevatorIndex].Ip &&  IsNotInElevator(elevators[localElevatorIndex], best.Order){ //Map er ikke sortert, så heiser velger forskjellig og tar samme ordre
				InsertOrder(elevators[localElevatorIndex], best.Order)
				updateNetworkChan <- elevators[localElevatorIndex]
			}

		
		case changedElevator = <- changedElevatorChan:
			fmt.Println("ChangedElevator")
			if localElevatorIndex == 0{
				if changedElevator.Type == "new"{
					fmt.Println("NEW elevator: ", changedElevator.Ip )
					if changedElevator.Ip != localIp{
					HandleNewElevator(elevators, changedElevator.Ip)
				}
				}else if changedElevator.Type == "dead"{
					fmt.Println("DEAD elevator: ", changedElevator.Ip )
					HandleDeadElev(elevators, changedElevator.Ip, deadOrderToUDPChan)
				}
			}



		case updateElevator = <- receiveDriverUpdateChan:
			fmt.Println("RecievedDriverUpdate")
			elevators[localElevatorIndex].Direction = updateElevator.Direction
			elevators[localElevatorIndex].LastFloor = updateElevator.LastFloor



		case updateElevator = <- receiveElevatorChan: // Recieves updates from all-over, updates accordingly
			fmt.Println("RecievedElevator")
			if updateElevator.Ip != localIp{
				for i := 0; i < N_ELEVATORS; i++{
					if elevators[i].Ip == updateElevator.Ip{
						//fmt.Println(elevators[i])
						elevators[i] = updateElevator
						break
					}
				}
			}
			
			


		//Updating the other module
		case <- timedUpdateChanNetwork: // Timed update to network
			fmt.Println("timedUpdateChanNetwork")
			fmt.Println(elevators)
			updateNetworkChan <- elevators[localElevatorIndex]
			fmt.Println("Waiting")
			
			fmt.Println(elevators)
		case <- timedUpdateChanDriver:
			fmt.Println("timedUpdateChanDriver")
			//fmt.Println("Sending :       ", elevators[localElevatorIndex])
			updateDriverChan <- elevators[localElevatorIndex]
		case <- timedLightUpdate:
			fmt.Println("timedLg")
			timedLightUpdate <- elevators
		
		
		}
	}
}
