package elevstatesync

// Functions to create and send syncronization messages to other network nodes

import (
	"encoding/json"
	"github.com/anlif/sanntid/heis/elevdriver"
)

type syncType int

const (
	InitialSync syncType = iota + 100
	ElevatorSynced
	ButtonPushed
	ElevatorStopped
	FloorReached
	OrderComplete
)

type syncData struct {
	DataType        syncType
	InitialSync     *sharedState       `json:",omitempty"`
	ElevatorSynced  *ElevatorState     `json:",omitempty"`
	ButtonPushed    *elevdriver.Button `json:",omitempty"`
	ElevatorStopped *int               `json:",omitempty"`
	FloorReached    *floorReached_T    `json:",omitempty"`
	OrderComplete   *int               `json:",omitempty"`
}

type floorReached_T struct {
	Floor int
	Dir   elevdriver.MoveDirection
}

type syncMessage struct {
	Id       int
	Sender   string
	SyncData syncData
	SendList []string
}

// Send initial sync message to a new connected elevator
func (self *SyncState) sendInitialSync(ipAddr string) {
	var stateData syncData
	var newMessage syncMessage
	stateData.DataType = InitialSync
	stateData.InitialSync = &self.sharedState

	newMessage = self.createSyncMessage(stateData)
	newMessage.SendList = []string{ipAddr}
	self.sendSyncMessage(newMessage)
}

func (self *SyncState) sendSyncMessage(message syncMessage) {
	marshaledSyncMsg, err := json.Marshal(message)
	checkError(err)
	for _, ip := range message.SendList {
		go func(senderIP string) { self.network.Connections[senderIP].SendChan <- marshaledSyncMsg }(ip)
	}
}

func (self *SyncState) createSyncMessage(data syncData) syncMessage {
	var newMessage syncMessage
	newMessage.Id = self.messageIndex
	newMessage.Sender = self.localIP
	newMessage.SendList = self.getOnlineElevators()
	newMessage.SyncData = data
	self.messageIndex++

	return newMessage
}
