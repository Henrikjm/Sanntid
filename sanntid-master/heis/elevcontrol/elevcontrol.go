package elevcontrol

import(
	"time"
	"github.com/anlif/sanntid/heis/elevdriver"
	"github.com/anlif/sanntid/heis/elevstatesync"
)

type Elevator struct{
	currentState elevstatesync.ElevatorState
	stateSyncher *elevstatesync.SyncState
	control elevControl
}

func (self *Elevator) Start() () {
	var err error
	// Get a valid initial state from the driver
	startFloor, startDir := elevdriver.Init()

	// Init state syncronization
	self.stateSyncher, err = elevstatesync.NewSyncState("9001")
	checkError(err)
	self.currentState = <-self.stateSyncher.NewStateChan
	self.currentState, err = self.stateSyncher.SyncFloorReached(startFloor, startDir)
	checkError(err)

	// Init statemachine and start generating events
	self.control.Init()
	time.Sleep(100*time.Millisecond)
	go self.generateEvents()
}

func (self *Elevator) generateEvents() {
	buttonChan := elevdriver.GetButtonChan()
	obsChan := elevdriver.GetObsChan()
	floorChan := elevdriver.GetFloorChan()
	stopChan := elevdriver.GetStopButtonChan()
	for {
		select {
		case <-stopChan:
			self.handleControlEvent(eMStopPushed)
		case <-obsChan:
			self.handleControlEvent(obstrSignal)
		case floor := <-floorChan:
			if floor == 1 {
				self.currentState.Dir = elevdriver.MOVE_UP
			} else if floor == 4 {
				self.currentState.Dir = elevdriver.MOVE_DOWN
			}
			var err error
			self.currentState, err = self.stateSyncher.SyncFloorReached(floor, self.currentState.Dir)
			checkError(err)
			elevdriver.SetFloor(floor)
			self.currentState.Floor = floor
			self.handleControlEvent(floorReached)
		case button := <-buttonChan:
			self.currentState, _ = self.stateSyncher.SyncButtonPush(button)
		case self.currentState = <-self.stateSyncher.NewStateChan:
		case <-time.After(50*time.Millisecond):
			if self.moreOrdersPending() {
				self.handleControlEvent(ordersPending)
			}
		}
	}
}
