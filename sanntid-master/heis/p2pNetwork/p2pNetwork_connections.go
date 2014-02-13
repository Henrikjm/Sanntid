package p2pNetwork

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

type messageType int

const (
	keep_alive_message messageType = iota + 1337
	user_message
)

type networkMessage struct {
	Type messageType
	Data []byte
}

func (self *NetworkNode) newConnection(newTCPConn *net.TCPConn) (*NodeConnection, string) {
	newConn := new(NodeConnection)
	newConn.tcpConn = newTCPConn
	newConn.SendChan = make(chan []byte)
	newConn.RecvChan = make(chan []byte)
	newConn.closeSig = make(chan int)
	newConn.connLostSig = make(chan int)

	newConnIP, _, err := net.SplitHostPort(newTCPConn.RemoteAddr().String())
	checkError(err)
	log.Printf("New connection IP: %s\n", newConnIP)
	self.Connections[newConnIP] = newConn
	return newConn, newConnIP
}

func (self *NetworkNode) closeNetwork() error {
	for addr, conn := range self.Connections { // Close all connections
		select {
		case conn.closeSig <- 1:
		case _, _ = <-conn.connLostSig: // can be closed already
		}
		close(conn.closeSig)
		delete(self.Connections, addr)
	}
	log.Printf("Closed Network\n")
	return nil
}

// Runtime routine of a network node, entry point after init & consistency check
func (self *NetworkNode) normalOperation() {
	self.ConnectionUpSig <- 1
	tcpConnChan := make(chan *net.TCPConn)
	quitChan := make(chan bool)

	go self.listenToUDP(tcpConnChan, quitChan)

	for {
		select {
		case newTCPConn := <-tcpConnChan:
			newConn, newIP := self.newConnection(newTCPConn)
			go newConn.handleConn()
			self.NewConnection <- newIP
		default:
			// Check if any connections are lost
			time.Sleep(time.Millisecond)
			for addr, conn := range self.Connections {
				select {
				case <-conn.connLostSig:
					if CheckConnectivity() == false { // Network cable unplugged
						self.ConnectionLost <- GetLocalIP()
						self.closeNetwork()
						quitChan <- true
						return
					} else {
						self.ConnectionLost <- addr
						delete(self.Connections, addr)
					}
				default:
				}
			}
		}
	}
}

// Runtime routine of a single node connection, handles sending and recieving messages asyncronously
func (self *NodeConnection) handleConn() {
	defer func() {
		err := self.tcpConn.Close()
		checkError(err)
		log.Printf("Closed connection to %s\n", self.tcpConn.RemoteAddr().String())
	}()

	var keepAliveMessage networkMessage = networkMessage{keep_alive_message, nil}
	var timeout time.Duration = 200 * time.Millisecond // We should recieve 'keep-alive' messages faster than this timeout
	var lastKeepAlive time.Time = time.Now()

	reader := bufio.NewReader(self.tcpConn)
	keepAliveMsgJSON, err := json.Marshal(keepAliveMessage)
	checkError(err)
	keepAliveMsgJSON = append(keepAliveMsgJSON, '\n') // Newline is needed to split incoming messages
	sendKeepAlive := time.After(timeout / 3)
	for {
		select {
		case <-self.closeSig:
			close(self.connLostSig)
			return
		case newData := <-self.SendChan:
			newPacket := networkMessage{user_message, newData}
			newPacketJSON, err := json.Marshal(newPacket)
			checkError(err)
			newPacketJSON = append(newPacketJSON, '\n') // We use ReadLine on recieve, need delimiter
			go func() {
				_, berr := self.tcpConn.Write(newPacketJSON)
				checkError(berr)
			}()
		case <-sendKeepAlive:
			// Send 'keep-alive' message
			go func() {
				_, err = self.tcpConn.Write(keepAliveMsgJSON)
				if err != nil { // Connection lost
					self.connLostSig <- 1
					close(self.connLostSig)
					return
				}
				sendKeepAlive = time.After(timeout / 3)
			}()
		default:
			// Try to recieve data, but do not block too long
			err := self.tcpConn.SetReadDeadline(time.Now())
			checkError(err)
			buf, err := reader.ReadBytes('\n')
			if err == nil { // Got data
				var message networkMessage
				err = json.Unmarshal(buf, &message)
				if err == nil {
					switch message.Type {
					case keep_alive_message:
						lastKeepAlive = time.Now()
					case user_message:
						go func(mess networkMessage) { self.RecvChan <- mess.Data }(message)
					default:
						panic("Type error in network message recieve")
					}
				} else if err != nil {
					checkError(err)
				}
			} else if err != nil && !isTimeoutErr(err) { // Error was not a timeout, connection lost
				log.Printf("Connection lost: %s\n", err)
				self.connLostSig <- 1
				close(self.connLostSig)
				return
			}

			if time.Now().Sub(lastKeepAlive) > timeout { // Connection timed out
				log.Printf("Timeout from %s\n", self.tcpConn.RemoteAddr().String())
				self.connLostSig <- 1
				close(self.connLostSig)
				return
			}
		}
	}
}

// Listen to UDP for new nodes, and initiate new connections
func (self *NetworkNode) listenToUDP(TCPChan chan<- *net.TCPConn, quit chan bool) {
	buf := make([]byte, 512)

	for {
		select {
		case <-quit:
			return
		default:
			udpConnection, err := net.ListenMulticastUDP("udp", nil, self.udpListenAddress)
			checkError(err)
			udpConnection.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
			n, _, err := udpConnection.ReadFromUDP(buf)
			if err == nil {
				var newConn broadcastMessage
				err = json.Unmarshal(buf[0:n], &newConn)
				checkError(err)
				if newConn.IP != self.tcpListenAddress.IP.String() || newConn.Port != self.tcpListenAddress.Port {
					tcpAddr, err := newConn.toTCPAddr()
					newTCPConn, err := self.initiatePeerConnection(tcpAddr)
					checkError(err)
					TCPChan <- newTCPConn
				}
			} else if !isTimeoutErr(err) {
				checkError(err)
			}
			err = udpConnection.Close()
			checkError(err)
		}
	}
}

func checkError(err error) {
	if err != nil {
		_, _, line, _ := runtime.Caller(1)
		log.Printf("Fatal error: %s, line: %d", err.Error(), line)
		os.Exit(1)
	}
}

func isTimeoutErr(err error) bool {
	opErr, isOpError := err.(*net.OpError)
	return (isOpError && opErr.Timeout())
}
