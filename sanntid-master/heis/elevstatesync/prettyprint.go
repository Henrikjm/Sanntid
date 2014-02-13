package elevstatesync

// Contains functions to print out data nicely for debugging

import (
	"fmt"
	"strings"
)

const marginPadding = "\t\t\t "

func (self syncType) String() (s string) {
	switch self {
	case InitialSync:
		s = "InitialSync"
	case ElevatorSynced:
		s = "ElevatorSynced"
	case ButtonPushed:
		s = "ButtonPushed"
	case ElevatorStopped:
		s = "ElevatorStopped"
	case FloorReached:
		s = "FloorReached"
	case OrderComplete:
		s = "OrderComplete"
	}
	return s
}

func (self syncData) String() string {
	var data string
	switch self.DataType {
	case InitialSync:
		data = self.InitialSync.String()
	case ElevatorSynced:
		data = self.ElevatorSynced.String()
	case ButtonPushed:
		data = self.ButtonPushed.String()
	case ElevatorStopped:
		data = ""
	case FloorReached:
		data = self.FloorReached.String()
	case OrderComplete:
		data = ""
	}
	return self.DataType.String() + " " + data
}

func (self *floorReached_T) String() string {
	return fmt.Sprintf("%d,%s", self.Floor, self.Dir)
}

func (self syncMessage) String() string {
	return fmt.Sprintf("ID:%d, Sender:%s \n %s %s", self.Id, strings.Split(self.Sender, ".")[3], marginPadding, self.SyncData.String())
}

func (self ElevatorStatus) String() (s string) {
	switch self {
	case ONLINE:
		s = "ONLINE"
	case STOPPED:
		s = "STOPPED"
	case OFFLINE:
		s = "OFFLINE"
	}
	return s
}

func (self OrderSlice) String() (s string) {
	for _, order := range self {
		if order.Active {
			s += "X"
		} else {
			s += "_"
			s += " "
		}
	}
	return s
}

func (self ElevatorState) String() string {
	s := fmt.Sprintf("Floor:%d Direction: %s Status: %s\n", self.Floor, self.Dir, self.Status) // + InternalOrders.String() + "\n" ExternalOrders.String()
	return s
}

func (self ElevatorMap) String() (s string) {
	for ip, state := range self {
		s += marginPadding + strings.Split(ip, ".")[3] + ": " + state.String()
	}
	return s
}

func (self PanelSlice) String() (s string) {
	for _, light := range self {
		if light {
			s += "*"
		} else {
			s += "_"
		}
		s += " "
	}
	return s
}

func (self sharedState) String() string {
	return fmt.Sprintf("Elevators: \n %s %s Panel: %s \n %s \t %s", self.Elevators, marginPadding, self.ExternalPanelState[0].String(), marginPadding, self.ExternalPanelState[1].String())
}
