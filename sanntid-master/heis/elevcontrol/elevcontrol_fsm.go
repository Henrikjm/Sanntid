package elevcontrol

import(
	"github.com/anlif/sanntid/heis/elevdriver"
	"time"
	"log"
	"os"
	"runtime"
)

const STOP_TIME = 3*time.Second
const N_C_STATES = 6
const N_C_EVENTS = 4

type controlEvent int
const(
	ordersPending controlEvent = iota
	obstrSignal
	floorReached
	eMStopPushed
)

type controlState int
const(
	idle controlState = iota
	execOrdr
	obstr
	stop
	eMObstr
	eMStop
)

type elevControl struct {
	CurrentState controlState
}

type stateActionT struct{
	State controlState
	action func(self *Elevator) error
	guard func(self *Elevator) bool
}

func (self controlEvent) String() (s string) {
	switch self {
	case ordersPending:
		s = "Orders pending"
	case obstrSignal:
		s = "Obstruction Signal"
	case floorReached:
		s = "Floor Reached"
	case eMStopPushed:
		s = "Emergency stop button pushed"
	default:
		s = "Invalid state"
	}
	return s
}

func (self controlState) String() (s string) {
	switch self {
	case idle:
		s = "IDLE"
	case execOrdr:
		s = "ExecutingOrder"
	case obstr:
		s = "Obstruction"
	case stop:
		s = "Stopped at floor"
	case eMStop:
		s = "Emergency stop"
	case eMObstr:
		s = "EmStop + Obstruction"
	default:
		s = "Invalid state"
	}
	return s
}

var eMStopA func(*Elevator) error = (*Elevator).eMStopAction
var execA func(*Elevator) error = (*Elevator).execAction
var obstrA func(*Elevator) error = (*Elevator).obstrAction
var stopA func(*Elevator) error = (*Elevator).stopAction
var exitEMA func(self *Elevator) error = (*Elevator).exitEMAction
var exEMOA func(self *Elevator) error = (*Elevator).exitEMObsAction
var idleA func(self *Elevator) error = (*Elevator).idleAction
var noA func(*Elevator) error = (*Elevator).noAction
var noG func(*Elevator)(bool) = (*Elevator).noGuard
var moreOrdrG func(*Elevator)(bool) = (*Elevator).moreOrdersPending
var hasOrdG func(self *Elevator)(bool) = (*Elevator).hasOrdGuard
var doorClG func(self *Elevator)(bool) = (*Elevator).doorClosedGuard

var controlStateTable [N_C_EVENTS][N_C_STATES]stateActionT = [N_C_EVENTS][N_C_STATES]stateActionT{
/*event/state 								idle					execOrd 				obstr					stop					EM+obstr			eMStop*/
/*ordersPendn*/	[N_C_STATES]stateActionT{{execOrdr,execA,noG},	{execOrdr,noA,noG},		{obstr,noA,noG},		{execOrdr,execA,doorClG},	{eMObstr,noA,noG},	{eMStop,noA,noG}},
/*obstrSignal*/	[N_C_STATES]stateActionT{{obstr,obstrA,noG},	{obstr,obstrA,noG},		{idle,idleA,noG},		{obstr,obstrA,noG},			{eMStop,noA,noG},	{eMObstr,noA,noG}},
/*floorReached*/[N_C_STATES]stateActionT{{stop,stopA,noG},		{stop,stopA,hasOrdG},	{obstr,noA,noG},		{stop,noA,noG},/*error*/	{eMObstr,noA,noG},	{eMStop,noA,noG}},
/*eMStopPushed*/[N_C_STATES]stateActionT{{eMStop,eMStopA,noG},	{eMStop,eMStopA,noG},	{eMObstr,eMStopA,noG},	{eMStop,eMStopA,noG},		{obstr,exEMOA,noG}, {idle,exitEMA,noG}}}

func (self *elevControl) Init() {
	self.CurrentState = idle
}

func (self *Elevator) handleControlEvent(event controlEvent) {
	next := controlStateTable[event][self.control.CurrentState]
	if next.guard(self) {
		if next.State != self.control.CurrentState {
			self.control.CurrentState = next.State
			log.Printf("Changed state to %s", next.State.String())
		}
		next.action(self)
	}
}

// Dummy function to pass as guard
func (self *Elevator) noGuard()(bool) {
	return true
}

// Dummy function to pass as action
func (self *Elevator) noAction() (error) {
	return nil
}

/*------------------Actions-------------------------*/
func (self *Elevator) execAction() (error) {
	currFloor := self.currentState.Floor
	currDir := self.currentState.Dir
	currOrder,_ := self.currentState.GetClosestOrder(currFloor,currDir)

	if (currOrder.Floor > currFloor){
		self.currentState.Dir = elevdriver.MOVE_UP
		elevdriver.MotorUp()
	} else if (currOrder.Floor < currFloor){
		self.currentState.Dir = elevdriver.MOVE_DOWN
		elevdriver.MotorDown()
	} else if(currOrder.Floor == currFloor){
		elevdriver.MotorStop()
		self.handleControlEvent(floorReached)
	}
	return nil
}

func (self *Elevator) obstrAction() (error) {
	elevdriver.MotorStop()
	self.currentState, _ = self.stateSyncher.SyncElevatorStopped()
	return nil
}

func (self *Elevator) eMStopAction() (error) {
	// Stop, clear/redistribute orders + sync
	elevdriver.MotorStop()
	elevdriver.SetStopButton()
	var err error
	self.currentState, err = self.stateSyncher.SyncElevatorStopped()
	checkError(err)
	return nil
}
func (self *Elevator) stopAction() (error) {
	self.currentState, _ = self.stateSyncher.SyncOrderComplete()
	elevdriver.MotorStop()
	elevdriver.SetDoor()
	go func() {
		time.Sleep(STOP_TIME)
		elevdriver.ClearDoor()
	} ()
	return nil
}
func (self *Elevator) exitEMAction() (error) {
	self.currentState, _ = self.stateSyncher.SyncElevatorOnline()
	if self.currentState.Dir == elevdriver.MOVE_DOWN {
		elevdriver.MotorDown()
	} else {
		elevdriver.MotorUp()
	}
	elevdriver.ClearStopButton()
	return nil
}

func (self *Elevator) exitEMObsAction() (error) {
	elevdriver.ClearStopButton()
	return nil
}

func (self *Elevator) idleAction() (error) {
	self.currentState, _ = self.stateSyncher.SyncElevatorOnline()
	elevdriver.MotorStop()
	return nil
}

/*------------------Guards-------------------------*/
func (self *Elevator) doorClosedGuard() (bool) {
	return !elevdriver.GetDoor()
}

func (self *Elevator) moreOrdersPending() (bool){
	_, moreOrders := self.currentState.GetClosestOrder(self.currentState.Floor, self.currentState.Dir)
	return moreOrders
}

func (self *Elevator) hasOrdGuard() (bool) {
	order, moreOrders := self.currentState.GetClosestOrder(self.currentState.Floor, self.currentState.Dir)
	return (order.Floor == self.currentState.Floor) || !moreOrders
}

func checkError(err error) {
	if err != nil {
		_, _, line, _ := runtime.Caller(1)
		log.Printf("Fatal error: %s, line: %d", err.Error(), line)
		os.Exit(1)
	}
}
