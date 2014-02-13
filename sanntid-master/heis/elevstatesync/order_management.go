package elevstatesync

// This file contains logic to add, remove and delegate orders

import (
	"github.com/anlif/sanntid/heis/elevdriver"
	"log"
	"math"
	"sort"
)

func (self *ElevatorState) GetClosestOrder(startFloor int, startDir elevdriver.MoveDirection) (closestOrder elevdriver.OrderT, gotOrders bool) {
	if startFloor == 1 || startDir == elevdriver.MOVE_STOP {
		startDir = elevdriver.MOVE_UP
	}
	if startFloor == elevdriver.N_FLOORS {
		startDir = elevdriver.MOVE_DOWN
	}

	startState := elevdriver.OrderT{startFloor, elevdriver.OrderDirection(startDir)}
	currentState := startState

	for {
		if hasOrder, order := self.hasOrder(currentState.Floor, currentState.Dir); hasOrder {
			return order, true
		} else {
			currentState = currentState.Increment()
		}
		if currentState == startState { // No orders found
			return elevdriver.OrderT{startState.Floor, elevdriver.ORDER_INTERNAL}, false
		}
	}
	return elevdriver.OrderT{1, elevdriver.ORDER_INTERNAL}, false
}

func (self *ElevatorState) hasOrder(floor int, dir elevdriver.OrderDirection) (bool, elevdriver.OrderT) {
	if self.InternalOrders[floor-1].Active {
		return true, elevdriver.OrderT{floor, elevdriver.ORDER_INTERNAL}
	} else if self.ExternalOrders[dir][floor-1].Active {
		return true, elevdriver.OrderT{floor, dir}
	}
	return false, elevdriver.OrderT{floor, elevdriver.ORDER_INTERNAL}
}

func (self *ElevatorState) removeOrder(order elevdriver.OrderT) (removedExternal bool) {
	removedExternal = false
	if order.Dir == elevdriver.ORDER_INTERNAL {
		self.InternalOrders[order.Floor-1].Active = false
	} else {
		self.ExternalOrders[order.Dir][order.Floor-1].Active = false
		removedExternal = true
	}
	return removedExternal
}

func (self *ElevatorState) addOrder(order elevdriver.OrderT) {
	newOrder := Order{true}
	if isValidOrder(order.Floor, order.Dir) {
		if order.Dir == elevdriver.ORDER_INTERNAL { //internal order
			self.InternalOrders[order.Floor-1] = newOrder
		} else {
			self.ExternalOrders[order.Dir][order.Floor-1] = newOrder
		}
	} else {
		panic("Ordered order.Floor out of range")
	}
}

func isValidOrder(floor int, dir elevdriver.OrderDirection) bool {
	return (!((floor < 1) || (floor > 4) || ((floor == 4) && (dir == elevdriver.ORDER_UP)) || ((floor == 1) && (dir == elevdriver.ORDER_DOWN))) && ((dir == elevdriver.ORDER_INTERNAL) || (dir == elevdriver.ORDER_UP) || (dir == elevdriver.ORDER_DOWN)))
}

func (self *SyncState) delegateOrders(ip string) {
	elev := self.Elevators[ip]
	var orders []elevdriver.OrderT
	var removedOrder elevdriver.OrderT
	for floor := 1; floor <= elevdriver.N_FLOORS; floor++ {
		if elev.ExternalOrders[elevdriver.ORDER_UP][floor-1].Active {
			removedOrder = elevdriver.OrderT{floor, elevdriver.ORDER_UP}
			elev.removeOrder(removedOrder)
			orders = append(orders, removedOrder)
		}
		if elev.ExternalOrders[elevdriver.ORDER_DOWN][floor-1].Active {
			removedOrder = elevdriver.OrderT{floor, elevdriver.ORDER_DOWN}
			elev.removeOrder(removedOrder)
			orders = append(orders, removedOrder)
		}
	}
	if self.Elevators[self.localIP].Status != ONLINE && len(self.getRunningElevators()) == 0 {
		log.Println("No running elevators, external orders lost")
		self.clearPanelState()
	} else {
		for _, removedOrder := range orders {
			delegatedIP := self.calculateWhoServes(removedOrder)
			self.Elevators[delegatedIP].addOrder(removedOrder)
		}
	}
}

func (self *SyncState) updateOrderLights() {
	for floor := 1; floor <= elevdriver.N_FLOORS; floor++ {
		for dir := elevdriver.ORDER_UP; dir <= elevdriver.ORDER_DOWN; dir++ {
			if isValidOrder(floor, dir) {
				if self.ExternalPanelState[dir][floor-1] {
					elevdriver.SetLight(floor, dir)
				} else {
					elevdriver.ClearLight(floor, dir)
				}
			}
		}
		if self.Elevators[self.localIP].InternalOrders[floor-1].Active {
			elevdriver.SetLight(floor, elevdriver.ORDER_INTERNAL)
		} else {
			elevdriver.ClearLight(floor, elevdriver.ORDER_INTERNAL)
		}
	}
}

func (self *SyncState) calculateWhoServes(order elevdriver.OrderT) string {
	log.Printf("HASH in calc %d\n", self.makeStateHash())
	costMap := make(map[string]int)
	onlineElevators := self.getRunningElevators()
	if self.Elevators[self.localIP].Status == ONLINE {
		onlineElevators = append(onlineElevators, self.localIP)
	}
	sort.Strings(onlineElevators)
	log.Printf("online elevs: %v\n", onlineElevators)
	log.Printf("elevs: \n %s", self.Elevators)
	for _, ip := range onlineElevators {
		costMap[ip] = self.Elevators[ip].getCost(order)
	}
	log.Printf("costMap %v", costMap)
	var min int = int(math.MaxInt32)
	var lowestBidder string
	for _, ip := range onlineElevators {
		if costMap[ip] < min {
			lowestBidder = ip
			min = costMap[ip]
		}
	}
	return lowestBidder
}

// Get the cost for an elevator to complete an order
// The heuristic used counts the number of floor passings needed to complete the order in addition to how many orders it has to complete before this order
func (self *ElevatorState) getCost(order elevdriver.OrderT) (penalty int) {
	currentPos := elevdriver.OrderT{self.Floor, elevdriver.OrderDirection(self.Dir)}
	startPos := currentPos
	penalty = abs(startPos.Floor - order.Floor)
	if self.Dir == elevdriver.MOVE_STOP {
		return penalty
	}
	for {
		if hasOrder, _ := self.hasOrder(currentPos.Floor, currentPos.Dir); hasOrder {
			penalty += 1
		}
		currentPos = currentPos.Increment()
		if currentPos == startPos {
			break
		}
	}
	return penalty
}

func abs(i int) int {
	if i < 0 {
		i = -i
	}
	return i
}
