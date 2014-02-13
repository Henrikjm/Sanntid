package p2pNetwork

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

// A node spams this message on UDP to allow other nodes to connect to it
type broadcastMessage struct {
	IP   string
	Port int
}

func (bMsg broadcastMessage) toTCPAddr() (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp4", bMsg.IP+":"+strconv.Itoa(bMsg.Port))
}

type consistencyMessage struct {
	IPlist []string
}

type consistencyResponse struct {
	ConsistencyOK bool // If this is true, consistency check was OK
}

type consistencySync struct {
	ConsistencySyncOK bool // If this is true, all nodes returned OK on consistency check
}

// Init a p2p network by broadcasting your own listen port
// via multicast UDP. Wait for incoming connections and check
// that the p2p network is consistent 
// (i.e everyone is connected to each other)
func (self *NetworkNode) initP2Pconnection() error {
	var timeout time.Duration = 1000 * time.Millisecond
	// Broadcast on UDP and wait for TCP connections from other clients
	var exitLoop bool = false
	for exitLoop == false {
		time.Sleep(10 * time.Millisecond)
		self.sendUDPBroadcast()
		select {
		case newTCPConn := <-self.tcpConnectChan:
			self.newConnection(newTCPConn)
		case <-time.After(timeout):
			exitLoop = true
		}
	}
	self.tcpListener.Close()
	return self.checkNetworkConsistency()
}

// Atomic action to start a network consistency check.
// Each connected node recieves a list of IP addresses
// currently connected to this node. 
// If the network is consistent this list should be equal on all nodes.
// See also initiatePeerConnection
func (self *NetworkNode) checkNetworkConsistency() error {
	var timeout time.Duration = 3000 * time.Millisecond
	var buf []byte = make([]byte, 1024)
	consMessage := self.makeConsistencyMessage()
	// For each connection, send consistency message and await reply
	for _, conn := range self.Connections {
		jsonConsMessage, err := json.Marshal(consMessage)
		checkError(err)
		_, err = conn.tcpConn.Write(jsonConsMessage)
		checkError(err)
		conn.tcpConn.SetReadDeadline(time.Now().Add(timeout)) // Add a timeout to the read operation
		n, err := conn.tcpConn.Read(buf)                      // This can time-out, meaning that the consistency check failed
		checkError(err)
		var reply consistencyResponse
		err = json.Unmarshal(buf[0:n], &reply)
		checkError(err)
		if !reply.ConsistencyOK {
			return fmt.Errorf("Consistency check failed, client ip: %s\n", conn.tcpConn.RemoteAddr().String())
		}
	}
	// Everyone replied OK on consistency, send sync message to everyone
	var syncMessage consistencySync
	syncMessage.ConsistencySyncOK = true
	for _, conn := range self.Connections {
		message, err := json.Marshal(syncMessage)
		checkError(err)
		_, err = conn.tcpConn.Write(message)
		checkError(err)
		go conn.handleConn()
	}
	return nil
}

// Atomic action to initiate a new TCP connection
// (i.e add another node to the peer-to-peer network)
// We expect a consistency message when we connect via TCP
// The consistency message is checked and a reply is sent.
// The connection is rejected if the consistency check fails
// See also the function checkNetworkConsistency
func (self *NetworkNode) initiatePeerConnection(remoteAddr *net.TCPAddr) (*net.TCPConn, error) {
	var timeout time.Duration = 3000 * time.Millisecond
	buf := make([]byte, 1024)
	TCPConn, err := net.DialTCP("tcp4", nil, remoteAddr)
	checkError(err)

	// Wait for a consistency message from the peer
	TCPConn.SetReadDeadline(time.Now().Add(timeout))
	n, err := TCPConn.Read(buf) // This can time-out, meaning that the consistency check failed
	checkError(err)
	var message consistencyMessage
	err = json.Unmarshal(buf[0:n], &message)
	checkError(err)
	status := self.checkConsistencyMessage(message)
	if status == false {
		TCPConn.Close()
		return nil, fmt.Errorf("Error in consistency check")
	}
	// Send a consistency response
	var reply consistencyResponse
	reply.ConsistencyOK = true
	replyJSON, err := json.Marshal(reply)
	checkError(err)
	_, err = TCPConn.Write(replyJSON)
	checkError(err)
	// Wait for sync message
	TCPConn.SetReadDeadline(time.Now().Add(timeout))
	n, err = TCPConn.Read(buf) // This can time-out, meaning that the consistency check failed
	checkError(err)
	var syncMessage consistencySync
	err = json.Unmarshal(buf[0:n], &syncMessage)
	checkError(err)
	if syncMessage.ConsistencySyncOK == false {
		TCPConn.Close()
		return nil, fmt.Errorf("Error in consistency sync")
	}

	// Everything OK, return new valid TCP connection
	return TCPConn, nil
}

func (self *NetworkNode) checkConsistencyMessage(message consistencyMessage) bool {
	myMessage := self.makeConsistencyMessage()
	myMessage.IPlist = append(myMessage.IPlist, GetLocalIP())
	sort.Strings(myMessage.IPlist)
	sort.Strings(message.IPlist)
	return strings.Join(myMessage.IPlist, "") == strings.Join(message.IPlist, "")
}

func pollTcpConnections(listener *net.TCPListener, newConnections chan<- *net.TCPConn) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil { // listener closed
			return
		}
		newConnections <- conn
	}
}

// Make a broadcast message, convert to JSON, and send it via UDP
func (self *NetworkNode) sendUDPBroadcast() {
	var err error
	var bcastMsg broadcastMessage
	bcastMsg.IP = self.tcpListenAddress.IP.String()
	bcastMsg.Port = self.tcpListenAddress.Port
	initMsg, err := json.Marshal(bcastMsg)
	checkError(err)
	_, err = self.udpCon.WriteToUDP(initMsg, self.udpListenAddress)
	checkError(err)
}

func (self *NetworkNode) makeConsistencyMessage() (message consistencyMessage) {
	message.IPlist = make([]string, 0)
	for ipAddr, _ := range self.Connections {
		message.IPlist = append(message.IPlist, ipAddr)
	}
	return message
}
