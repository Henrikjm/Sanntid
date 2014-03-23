package queue

import(
	"fmt"
	"time"
	."types"
	"strings"
	"strconv"
	"math"
	"encoding/json"
    "io/ioutil"
)

func GetElevatorCost(elevator Elevator, order Order, updateFloorChan chan int) int {

	if len(elevator.OrderQueue)<1{
		fmt.Println("CORRUPTED LOCAL ELEVATOR!!!!")
		for len(elevator.OrderQueue)<MAX_ORDERS{
			elevator.OrderQueue = append(elevator.OrderQueue, Order{0,0})
		}
	}

	score := 0
	//correct direction plays 5 points difference
	if elevator.Direction == MOVE_STOP{
		score += GetNumberOfFloorsForOrder(updateFloorChan, order)

	}else if ((elevator.Direction == MOVE_UP) && (order.Floor > elevator.LastFloor)) || ((elevator.Direction == MOVE_DOWN) && (order.Floor < elevator.LastFloor)){ //hvis bestilling er i riktig retning
		score -= 5
	}else{
		score += 5
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

	for j:= 0; j < placement; j++{ //Removing common objective orders from score
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
	for i := 0; i < len(elevator.OrderQueue); i++{
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

func InsertOrder(elevator Elevator, order Order) Elevator{
	fmt.Println("|||||||||||||Inserting order||||||||||")
	if order.Floor == 0{
		fmt.Println("ERROR in InsertOrder: order.Floor == 0")
		return elevator
	}
	placement := GetInsertOrderPlacement(elevator, order)
	fmt.Println("Placement of ", order, " = ", placement)

	var temp, insert Order
	insert = order
	for i := placement; i <len(elevator.OrderQueue); i++{
		temp = elevator.OrderQueue[i]
		elevator.OrderQueue[i] = insert
		insert = temp
	}
	return elevator
}

func GetLocalElevatorIndex(elevList []Elevator, localIp string)int{
	for i := 0; i < len(elevList); i++{
		if elevList[i].Ip == localIp{
			fmt.Println("LocalElevatorIndex of ", elevList[i].Ip, " = ", i)
			return i
		}
	}
	fmt.Println("ERROR local Ip not found i elevators")
	return -1
}

func HandleDeadElev(elevList []Elevator, ip string, deadOrderToUDPChan chan Order) []Elevator{
	var i int
	var deadElevQueue []Order
	for i = 0 ; i < len(elevList); i++{
		if elevList[i].Ip == ip{
			deadElevQueue = elevList[i].OrderQueue
			elevList[i] = Elevator{}
			break
		}
	}
	for i = 0; i < len(deadElevQueue); i++{
		if deadElevQueue[i].Orientation != ORDER_INTERNAL && deadElevQueue[i].Orientation != 0{
			deadOrderToUDPChan <- deadElevQueue[i]
		}
	}
	return elevList
}

func HandleNewElevator(elevList []Elevator, ip string) []Elevator{
	for i := 0; i < len(elevList); i++{
		if elevList[i].Ip == ""{ 
			elevList[i].Ip = ip
			elevList[i].OrderQueue = make([]Order,MAX_ORDERS)
			break
		}
	}
	return elevList
}





func IsNotInElevator(elevator Elevator, order Order) bool {
	for i :=0; i < len(elevator.OrderQueue); i++{
		if elevator.OrderQueue[i] == order{
			return false
		}
	}
	return true
}

func checkForInternalOrderBackup(elevator Elevator) Elevator{

	dat, err := ioutil.ReadFile("internalOrderBackupFile")
	var readOrders []int

	if err != nil {
		internalOrders := []int{0,0,0,0}
		d1,_ := json.Marshal(internalOrders)
	    err := ioutil.WriteFile("internalOrderBackupFile", d1, 0644)
	    fmt.Println("No internal orders stored. Writing to new file...")
	    if err != nil{
	    	fmt.Println("ERROR!!! writing to internalOrderBackupFile.")
	    }

	    dat, err := ioutil.ReadFile("internalOrderBackupFile")

		json.Unmarshal(dat, &readOrders)
	    fmt.Println("Made new file with empty orders: ", readOrders)
	}else{

		json.Unmarshal(dat, &readOrders)
		fmt.Println("Found file for internal orders. The orders were: ", readOrders)
	}
	if len(readOrders) != 4 {
		fmt.Println("CORRUPTED READ!!! Parsing empty order list to elevator.")
		readOrders = []int{0,0,0,0}
	}

	for i,value := range readOrders{
		if value == 1 {
			elevator = InsertOrder(elevator, Order{i+1,ORDER_INTERNAL})
		}
	}

	return elevator

}

func UpdateInternalOrderBackupFile(elevator Elevator){
	internalOrders := []int{0,0,0,0}

	for _, orderInstance := range elevator.OrderQueue{
		if orderInstance.Orientation == ORDER_INTERNAL{
			internalOrders[orderInstance.Floor-1] = 1
		}
	}

	d1,_ := json.Marshal(internalOrders)
    err := ioutil.WriteFile("internalOrderBackupFile", d1, 0644)

    if err != nil{
    	fmt.Println("ERROR!!! writing to internalOrderBackupFile.")
    }
}



func QueueHandler(receiveElevatorChan chan Elevator, updateNetworkChan chan Elevator, newOrderFromUDPChan chan Order, deadOrderToUDPChan chan Order, sendCostChan chan Cost, recieveCostChan chan map[string]Cost, 
	changedElevatorChan chan Change, localIpChan chan string, localOrderChan chan Order, updateDriverChan chan Elevator, receiveDriverUpdateChan chan Elevator, orderToNetworkChan chan Order,
	 updateFloorChan chan int, timedLightUpdate chan []Elevator, localUpdateDriverChan chan Elevator, updateFromDriverChan chan Elevator, readyForUpdateChan chan bool){

	fmt.Println("QueueHandler started.")	



	//Making situation picture
	elevators := make([]Elevator, N_ELEVATORS) 
	localIpChan <- "LocalIp"
	localIp := <- localIpChan
	elevators = HandleNewElevator(elevators, localIp) 
	localElevatorIndex := GetLocalElevatorIndex(elevators, localIp)
	elevators[localElevatorIndex] = checkForInternalOrderBackup(elevators[localElevatorIndex])
	receiveDriverUpdateChan <- elevators[0]
	elevators[localElevatorIndex] = <- receiveDriverUpdateChan 

	updateDriverChan <- elevators[localElevatorIndex]

	fmt.Println("QueueHandler initiated.")



	for{
		time.Sleep(time.Millisecond * 1)

		select{
		
		case localOrder := <- localOrderChan:
			fmt.Println("localOrders")
			if IsNotInElevator(elevators[localElevatorIndex], localOrder){
				if localOrder.Orientation == ORDER_INTERNAL{
					elevators[localElevatorIndex] = InsertOrder(elevators[localElevatorIndex], localOrder)
					//localUpdateDriverChan <- elevators[localElevatorIndex] //Bør legge inn localUpdate i control
				}else{
					orderToNetworkChan <- localOrder
				}
			}		

		case newOrder := <-newOrderFromUDPChan:
			//fmt.Println("RecievedNewOrder")
			localCost := Cost{GetElevatorCost(elevators[localElevatorIndex], newOrder,  updateFloorChan), newOrder, elevators[localElevatorIndex].Ip}
			sendCostChan <- localCost


		case receivedCostMap := <- recieveCostChan:
			//fmt.Println("RecievedCost")
			best := Cost{}
			best.Cost = 20
			for _, recievedCost :=  range receivedCostMap{
				if recievedCost.Cost < best.Cost{
					best = recievedCost
				}
			}

			highestIp := 0 //strconv.Atoi(dummyIpStr)
			newIp := 0
			for _,recievedCost := range receivedCostMap{
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


		case changedElevator := <- changedElevatorChan:
			//fmt.Println("ChangedElevator")
			if localElevatorIndex == 0{
				if changedElevator.Type == "new"{
					fmt.Println("NEW elevator: ", changedElevator.Ip )
					if changedElevator.Ip != localIp{
					elevators = HandleNewElevator(elevators, changedElevator.Ip)
				}
				}else if changedElevator.Type == "dead"{
					fmt.Println("DEAD elevator: ", changedElevator.Ip )
					elevators = HandleDeadElev(elevators, changedElevator.Ip, deadOrderToUDPChan)
				}
			}

		case updateElevatorFromNetwork := <- receiveElevatorChan: // Recieves updates from all-over, updates accordingly
			//fmt.Println("RecievedElevator")
			//fmt.Println(updateElevator.Ip, updateElevator.OrderQueue[0])
			if updateElevatorFromNetwork.Ip != localIp{
				for i := 0; i < N_ELEVATORS; i++{
					if elevators[i].Ip == updateElevatorFromNetwork.Ip{
						elevators[i] = updateElevatorFromNetwork
						break
					}
				}
			}

		//Updating the other module
		case elevators[localElevatorIndex] = <- updateFromDriverChan:
			//fmt.Println("timedUpdateChanDriver")
			updateFromDriverChan <- elevators[localElevatorIndex]

		case <- readyForUpdateChan: // Timed update to network
			//fmt.Println("timedUpdate")
			updateNetworkChan <- elevators[localElevatorIndex]
			timedLightUpdate <- elevators
			updateDriverChan <- elevators[localElevatorIndex]
			go func(elevator Elevator){UpdateInternalOrderBackupFile(elevator)}(elevators[localElevatorIndex])

	}
}
}