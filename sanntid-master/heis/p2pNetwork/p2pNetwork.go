package p2pNetwork

import (
	"net"
	"time"
)

const UDPMulticastPort = "8726"
const UDPMulticastAddr = "224.0.1.60"

// A network node maintains a map of connections on the p2p network
// When a connection is lost, it is notified on the ConnectionLost channel before it is deleted. 
// If the lost connection is equal to the local IP address, the connection has been closed.
type NetworkNode struct {
	Connections      map[string]*NodeConnection // The key is the remote IP address
	ConnectionUpSig  chan int                   // This is signaled after connection is up and running, used for syncronization
	ConnectionLost   chan string                // string is IP address to lost connection
	NewConnection    chan string                // string is IP address to new connection
	tcpConnectChan   chan *net.TCPConn          // new raw tcpConnections
	tcpListener      *net.TCPListener
	udpCon           *net.UDPConn
	udpListenAddress *net.UDPAddr
	tcpListenAddress *net.TCPAddr
}

// Represents a single node connection, send data to it on
// SendChan and recieve data on RecvChan.
type NodeConnection struct {
	SendChan    chan []byte
	RecvChan    chan []byte
	tcpConn     *net.TCPConn
	closeSig    chan int
	connLostSig chan int
}

// Create a new network node and try to add it to the existing network
func NewNetworkNode(listenAddress string, listenPort string) (*NetworkNode, error) {
	n := new(NetworkNode)
	n.Connections = make(map[string]*NodeConnection, 0)
	n.ConnectionUpSig = make(chan int)
	n.ConnectionLost = make(chan string)
	n.NewConnection = make(chan string)
	n.tcpConnectChan = make(chan *net.TCPConn, 0)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", listenAddress+":"+listenPort)
	checkError(err)
	n.tcpListenAddress = tcpAddr

	// Setup TCP listener
	n.tcpListener, err = net.ListenTCP("tcp4", n.tcpListenAddress)
	checkError(err)
	go pollTcpConnections(n.tcpListener, n.tcpConnectChan)

	// Setup UDP listener
	localAddr, err := net.ResolveUDPAddr("udp4", ":0")
	checkError(err)
	n.udpCon, err = net.ListenUDP("udp4", localAddr)
	checkError(err)

	multicastAddr, err := net.ResolveUDPAddr("udp4", UDPMulticastAddr+":"+UDPMulticastPort)
	checkError(err)
	n.udpListenAddress = multicastAddr

	err = n.initP2Pconnection()
	if err != nil {
		return nil, err
	}
	go n.normalOperation()
	return n, nil
}

// Try to reconnect after losing connection, blocks
func (self *NetworkNode) Reconnect() error {
	for CheckConnectivity() == false {
		time.Sleep(300 * time.Millisecond)
	}
	var err error
	self.tcpListener, err = net.ListenTCP("tcp4", self.tcpListenAddress)
	checkError(err)
	go pollTcpConnections(self.tcpListener, self.tcpConnectChan)
	time.Sleep(300 * time.Millisecond)
	err = self.initP2Pconnection()
	if err == nil {
		go self.normalOperation()
		return nil
	} else {
		return err
	}
	return nil
}

func GetLocalIP() string {
	addresses, _ := net.InterfaceAddrs()
	finalIP := "127.0.0.1"
	for _, address := range addresses {
		strAddr := address.String()
		prefix := strAddr[0:3]
		if prefix == "129" {
			finalIP = strAddr[0:15]
			break
		}
	}
	return finalIP
}

// Gets the local TCP listen address
func (self *NetworkNode) GetListenAddr() string {
	return self.tcpListenAddress.String()
}

// Get the IP addresses of connected nodes
func (self *NetworkNode) GetNodeList() []string {
	consMessage := self.makeConsistencyMessage()
	return consMessage.IPlist
}

// Returns true if you are connected to the network
func CheckConnectivity() bool {
	connectivityOK := make(chan int)
	timeout := 100 * time.Millisecond
	go func() {
		_, err := net.LookupHost("www.ntnu.no")
		checkError(err)
		if err == nil {
			connectivityOK <- 1
		}
	}()

	select {
	case <-connectivityOK:
		return true
	case <-time.After(timeout):
		return false
	}

	return true
}
