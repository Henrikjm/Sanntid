package elevstatesync

import (
	"encoding/json"
	"github.com/anlif/sanntid/heis/elevdriver"
	"github.com/anlif/sanntid/heis/p2pNetwork"
	"hash/adler32"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
)

type Order struct {
	Active bool
}

type ElevatorStatus int

const (
	ONLINE ElevatorStatus = iota
	STOPPED
	OFFLINE
)

type OrderSlice [elevdriver.N_FLOORS]Order

type ElevatorState struct {
	Floor          int
	Dir            elevdriver.MoveDirection
	Status         ElevatorStatus
	InternalOrders OrderSlice
	ExternalOrders [2]OrderSlice
}

type ElevatorMap map[string]*ElevatorState // Key is IP address to elevator
type PanelSlice [elevdriver.N_FLOORS]bool

// This exported state is synced between all elevators
type sharedState struct {
	Elevators          ElevatorMap
	ExternalPanelState [2]PanelSlice // Order panel state, should be the same on all elevators
}

func newElevatorState() *ElevatorState {
	elev := new(ElevatorState)
	elev.Floor = 1
	elev.Dir = elevdriver.MOVE_UP
	elev.Status = ONLINE
	return elev
}

func (self *SyncState) SyncElevatorOnline() (ElevatorState, error) {
	var myState *ElevatorState = self.Elevators[self.localIP]
	var newData syncData
	myState.Status = ONLINE
	newData.DataType = ElevatorSynced
	newData.ElevatorSynced = myState
	self.syncDataChan <- newData
	return <-self.NewStateChan, nil
}

func (self *SyncState) SyncButtonPush(button elevdriver.Button) (ElevatorState, error) {
	var newData syncData
	newData.DataType = ButtonPushed
	newData.ButtonPushed = &button
	self.syncDataChan <- newData
	return <-self.NewStateChan, nil
}

func (self *SyncState) SyncFloorReached(floor int, dir elevdriver.MoveDirection) (ElevatorState, error) {
	var newData syncData
	newData.DataType = FloorReached
	newData.FloorReached = &floorReached_T{floor, dir}
	self.syncDataChan <- newData
	return <-self.NewStateChan, nil
}

func (self *SyncState) SyncOrderComplete() (ElevatorState, error) {
	var newData syncData
	newData.DataType = OrderComplete
	newData.OrderComplete = nil
	self.syncDataChan <- newData
	return <-self.NewStateChan, nil
}

func (self *SyncState) SyncElevatorStopped() (ElevatorState, error) {
	var newData syncData
	newData.DataType = ElevatorStopped
	self.syncDataChan <- newData
	return <-self.NewStateChan, nil
}

// This object contains all shared and local state, state changes from other elevators are syncronized sequentially through the NewStateChan
type SyncState struct {
	sharedState
	NewStateChan chan ElevatorState
	network      *p2pNetwork.NetworkNode
	localIP      string
	synced       bool
	syncDataChan chan syncData
	messageIndex int
}

func NewSyncState(portNum string) (*SyncState, error) {
	var stateObj *SyncState = new(SyncState)
	var err error = nil
	stateObj.Elevators = make(map[string]*ElevatorState)
	stateObj.syncDataChan = make(chan syncData)
	stateObj.NewStateChan = make(chan ElevatorState)

	myIP := p2pNetwork.GetLocalIP()
	stateObj.network, err = p2pNetwork.NewNetworkNode(myIP, portNum)
	if err != nil {
		return nil, err
	}
	<-stateObj.network.ConnectionUpSig
	stateObj.localIP = p2pNetwork.GetLocalIP()

	log.SetFlags(log.Lshortfile)
	log.SetPrefix(strings.Split(myIP, ".")[3] + ":")

	if len(stateObj.network.Connections) == 0 {
		stateObj.synced = true // Alone on the network -> synced by default
		stateObj.Elevators[stateObj.localIP] = newElevatorState()
		go func() { stateObj.NewStateChan <- *stateObj.Elevators[stateObj.localIP] }()
	} else {
		stateObj.synced = false // Need to wait for sync from other elevator
		onlineElevs := stateObj.network.GetNodeList()
		sort.Strings(onlineElevs)
		smallestIP := onlineElevs[0]
		newData := <-stateObj.network.Connections[smallestIP].RecvChan
		var newSyncMessage syncMessage
		err := json.Unmarshal(newData, &newSyncMessage)
		checkError(err)
		go stateObj.updateState(newSyncMessage)
	}

	go stateObj.runSync()

	return stateObj, nil
}

func (self *SyncState) runSync() {
	for {
		select {
		case newConnIP := <-self.network.NewConnection:
			log.Printf("New connection: %s\n", newConnIP)
			// Only the elevator with the smallest IP address sends initial sync message (only one is needed)
			smallestIP := self.localIP
			for _, ip := range self.getOnlineElevators() {
				if ip < smallestIP {
					smallestIP = ip
				}
			}
			if smallestIP == self.localIP {
				self.sendInitialSync(newConnIP)
			}
		case <-self.network.ConnectionUpSig:
			log.Printf("Reconnected!")
		case connLostIP := <-self.network.ConnectionLost:
			log.Printf("Lost self.network connection: %v\n", connLostIP)
			self.Elevators[connLostIP].Status = OFFLINE
			self.delegateOrders(connLostIP)
			if connLostIP == self.localIP {
				log.Printf("Lost Network Connection")
				self.clearPanelState()
				self.synced = false
				go func() {
					time.Sleep(1 * time.Second)
					err := self.network.Reconnect()
					checkError(err)
				}()
			}
            self.NewStateChan <- *self.Elevators[self.localIP]
		case newData := <-self.syncDataChan:
			var newSyncMessage syncMessage = self.createSyncMessage(newData)
			if self.synced == true {
				self.sendSyncMessage(newSyncMessage)
			}
			self.updateState(newSyncMessage)
		default:
			time.Sleep(time.Microsecond)
			for _, conn := range self.network.Connections {
				select {
				case newData := <-conn.RecvChan:
					var newSyncMessage syncMessage
					err := json.Unmarshal(newData, &newSyncMessage)
					checkError(err)
					self.updateState(newSyncMessage)
				default:
					time.Sleep(time.Microsecond)
				}
			}
		}
	}
}

func (self *SyncState) updateState(message syncMessage) {
	log.Printf("NEW SYNC MESSAGE: %v\n", message)
	switch message.SyncData.DataType {
	case InitialSync:
		if self.synced == true {
			panic("Recieved sync message when already synced from " + message.Sender) // This should not happen
		} else {
			newState := message.SyncData.InitialSync
			if _, inNewList := newState.Elevators[self.localIP]; !inNewList { // Not in new elevator list, create initial state
				var myState *ElevatorState = newElevatorState()
				self.Elevators = newState.Elevators
				self.Elevators[self.localIP] = myState
			} else if _, inLocalList := self.Elevators[self.localIP]; inLocalList { // Already in local list, keep local orders
				newState.Elevators[self.localIP] = self.Elevators[self.localIP]
				self.Elevators = newState.Elevators
			} else { // Already in new list and not in local, overwrite all local state
				self.Elevators = newState.Elevators
			}
			self.ExternalPanelState = newState.ExternalPanelState
			self.Elevators[self.localIP].Status = ONLINE
			self.synced = true
			self.sendElevatorSynced()
		}
	case ElevatorSynced:
		self.Elevators[message.Sender] = message.SyncData.ElevatorSynced
		log.Printf("%s elevator synced: %s\n", self.localIP, message.Sender)
	case ButtonPushed:
		var button elevdriver.Button = *message.SyncData.ButtonPushed
		if button.Dir == elevdriver.ORDER_INTERNAL {
			//add order internally if it is not already added
			self.Elevators[message.Sender].addOrder(elevdriver.OrderT(button))
		} else if (self.synced == true) && (self.ExternalPanelState[button.Dir][button.Floor-1] == false) {
			//assign external order to some elevator if it is not already added, and ther are online elevators
			if self.Elevators[self.localIP].Status != ONLINE && len(self.getRunningElevators()) == 0 {
				log.Println("No running elevators, external order not registered")
			} else {
				ip := self.calculateWhoServes(elevdriver.OrderT(button))
				log.Printf("IP returned from calc: %s", ip)
				self.Elevators[ip].addOrder(elevdriver.OrderT(button))
				self.ExternalPanelState[button.Dir][button.Floor-1] = true
			}
		}
	case ElevatorStopped:
		// Remove external orders from elevator and delegate them to others
		elev := self.Elevators[message.Sender]
		elev.Status = STOPPED
		self.delegateOrders(message.Sender)
	case FloorReached:
		var newState floorReached_T = *message.SyncData.FloorReached
		elev := self.Elevators[message.Sender]
		elev.Floor = newState.Floor
		elev.Dir = newState.Dir
	case OrderComplete:
		elev := self.Elevators[message.Sender]
		closestOrder, gotOrders := elev.GetClosestOrder(elev.Floor, elev.Dir)
		for gotOrders && closestOrder.Floor == elev.Floor {
			removedExternal := elev.removeOrder(closestOrder)
			if removedExternal {
				self.ExternalPanelState[closestOrder.Dir][closestOrder.Floor-1] = false
			}
			closestOrder, gotOrders = elev.GetClosestOrder(elev.Floor, elev.Dir)
		}
	}
	self.updateOrderLights()
	log.Printf("HASH AFTER UPDATE: %d\n", self.makeStateHash())
	log.Printf("internal:\t%v\n", self.Elevators[self.localIP].InternalOrders.String())
	log.Printf("external:\t%v\n", self.Elevators[self.localIP].ExternalOrders[0].String())
	log.Printf("\t \t%v\n", self.Elevators[self.localIP].ExternalOrders[1].String())
	log.Printf("panel:\t%v\n", self.ExternalPanelState[0])
	log.Printf("\t \t%v\n", self.ExternalPanelState[1])
	self.NewStateChan <- *self.Elevators[self.localIP]
}

func (self *SyncState) clearPanelState() {
	for floor := 1; floor <= elevdriver.N_FLOORS; floor++ {
		self.ExternalPanelState[elevdriver.ORDER_UP][floor-1] = false
		self.ExternalPanelState[elevdriver.ORDER_DOWN][floor-1] = false
	}
}

func (self *SyncState) sendElevatorSynced() {
	var myState *ElevatorState = self.Elevators[self.localIP]
	var newData syncData
	newData.DataType = ElevatorSynced
	newData.ElevatorSynced = myState

	self.sendSyncMessage(self.createSyncMessage(newData))
}

func (self *SyncState) getOnlineElevators() []string {
	elevatorList := make([]string, 0)
	for ip, state := range self.Elevators {
		if state.Status != OFFLINE && ip != self.localIP {
			elevatorList = append(elevatorList, ip)
		}
	}
	return elevatorList
}

func (self *SyncState) getRunningElevators() []string {
	elevatorList := make([]string, 0)
	for ip, state := range self.Elevators {
		if state.Status == ONLINE && ip != self.localIP {
			elevatorList = append(elevatorList, ip)
		}
	}
	return elevatorList
}

func (self *SyncState) makeStateHash() uint32 {
	marshaledState, err := json.Marshal(self.sharedState)
	checkError(err)
	return adler32.Checksum(marshaledState)
}

func checkError(err error) {
	if err != nil {
		_, _, line, _ := runtime.Caller(1)
		log.Printf("Fatal error: %s, line: %d", err.Error(), line)
		os.Exit(1)
	}
}
