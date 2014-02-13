package elevstatesync

import (
	"encoding/json"
	"testing"
	"time"
)

func testJsonPack(t *testing.T) {
	var myState ElevatorState
	t.Logf("state: %v\n", myState)

	marshaled, err := json.Marshal(myState)
	t.Logf("marshaled: %v\n, error: %v", string(marshaled), err)
}

// Make three nodes and send a sync message
func TestSendSyncMessage(t *testing.T) {
	s1, err := NewSyncState("9007")
	if err != nil {
		t.FailNow()
	}
	t.Logf("s1 node list: %v\n", s1.network.GetNodeList())

	//t.Logf("s1: %v\n", s1)
	time.Sleep(100 * time.Millisecond)

	s2, err := NewSyncState("9008")
	if err != nil {
		t.FailNow()
	}
	t.Logf("s2 node list: %v\n", s2.network.GetNodeList())
	time.Sleep(100 * time.Millisecond)

	s3, err := NewSyncState("9009")
	if err != nil {
		t.FailNow()
	}
	t.Logf("s3 node list: %v\n", s3.network.GetNodeList())

	time.Sleep(2 * time.Second)

	t.Logf("s1 %d, s2 %d, s3 %d\n", s1.makeStateHash(), s2.makeStateHash(), s3.makeStateHash())

}
