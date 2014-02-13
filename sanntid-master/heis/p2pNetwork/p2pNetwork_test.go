package p2pNetwork

import (
	"log"
	"testing"
	"time"
)

// Make a few network nodes and do elementary testing on these
func TestInit(t *testing.T) {
	net1, err := NewNetworkNode("127.0.0.1", "9001")
	if err != nil {
		t.Fatalf("Error in net1 init: %s\n", err)
	} else {
		t.Log("net1 init A-OK")
	}
	<-net1.ConnectionUpSig
	time.Sleep(100 * time.Millisecond)

	net2, err := NewNetworkNode("127.0.0.1", "9002")
	if err != nil {
		t.Fatalf("Error in net2 init: %s\n", err)
	} else {
		t.Log("net2 init A-OK")
	}
	<-net2.ConnectionUpSig
	time.Sleep(100 * time.Millisecond)

	net3, err := NewNetworkNode("127.0.0.1", "9003")
	if err != nil {
		t.Fatalf("Error in net3 init: %s\n", err)
	} else {
		t.Log("net3 init A-OK")
	}
	<-net3.ConnectionUpSig
	time.Sleep(100 * time.Millisecond)

	time.Sleep(100 * time.Millisecond)
	t.Logf("node1: %v\n", net1.GetNodeList())
	t.Logf("node2: %v\n", net2.GetNodeList())
	t.Logf("node3: %v\n", net3.GetNodeList())

	/*
		net1.Close()
		time.Sleep(100*time.Millisecond)
		t.Logf("Net1 closed")
		time.Sleep(100*time.Millisecond)
		t.Logf("node2: %v\n", net2.GetNodeList())
		net2.Close()
		t.Logf("Net2 closed")
		time.Sleep(100*time.Millisecond)
		t.Logf("node3: %v\n", net3.GetNodeList())
		net3.Close()
		t.Logf("Net3 closed")
	*/
}

func TestGetIP(t *testing.T) {
	t.Log("IP: " + GetLocalIP())
}

func testCloseAndOpen(t *testing.T) {
	net1, err := NewNetworkNode("127.0.0.1", "9001")
	if err != nil {
		t.Fatalf("Error in net1 init: %s\n", err)
	} else {
		t.Log("net1 init A-OK")
	}
	<-net1.ConnectionUpSig
	time.Sleep(100 * time.Millisecond)

	net2, err := NewNetworkNode("127.0.0.1", "9002")
	if err != nil {
		t.Fatalf("Error in net2 init: %s\n", err)
	} else {
		t.Log("net2 init A-OK")
	}
	<-net2.ConnectionUpSig

	time.Sleep(100 * time.Millisecond)
	t.Logf("node1: %v\n", net1.GetNodeList())
	t.Logf("node2: %v\n", net2.GetNodeList())

	net1.Close()
	log.Printf("net1 connection closed\n")

	time.Sleep(100 * time.Millisecond)

	addr := <-net2.ConnectionLost
	time.Sleep(100 * time.Millisecond)
	t.Logf("Connection lost: %v\n", addr)
	t.Logf("node2: %v\n", net2.GetNodeList())

	net1.Reconnect()
	<-net1.ConnectionUpSig
	time.Sleep(100 * time.Millisecond)
	t.Logf("After reconnect, node1: %v, node2: %v\n", net1.GetNodeList(), net2.GetNodeList())

}

func testSingleNode(t *testing.T) {
	net3, err := NewNetworkNode(GetLocalIP(), "9001")
	if err != nil {
		t.Fatalf("Error in net3 init: %s\n", err)
	} else {
		t.Log("net3 init A-OK")
	}
	<-net3.ConnectionUpSig
	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 2; i++ {
		time.Sleep(time.Second)
		t.Logf("nodelist: %v\n", net3.GetNodeList())
	}
}

func TestConnectivity(t *testing.T) {
	t.Logf("Connectivity: %v\n", CheckConnectivity())
}
