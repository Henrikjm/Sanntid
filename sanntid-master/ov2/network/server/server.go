package main

import (
	"bufio"
	"container/list"
	"encoding/json"
	"log"
	"net"
	"sort"
	"time"
)

func main() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("ip4", service)
	checkError(err)
	log.Println("server addr: ", tcpAddr.String())

	pendingConnections := list.New()
	activeConnections := make(map[string]ClientConnection)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	rawConnections := make(chan net.Conn, 10)

	go pollTcpConnections(listener, rawConnections)

	for {
		// Check for new connections (non-blocking)
		select {
		case tcpConn := <-rawConnections:
			currentCon := NewClientConnection()
			pendingConnections.PushBack(currentCon)
			go currentCon.Routine(tcpConn)
		case <-time.After(10*time.Millisecond):
		}

		// Iterate pending connections to see if any have a valid ID
		for elem := pendingConnections.Front(); elem != nil; elem = elem.Next() {
			// Cast value to a ClientConnection
			v := elem.Value.(*ClientConnection)
			select {
			case status := <-v.StatusChan:
				if status == ConnectionLost {
					pendingConnections.Remove(elem)
				} else if status == UserIDRecieved {
					// add user to active connections, TODO: handle collisions
					activeConnections[v.ClientID] = *v
					log.Printf("%s\n", "Added "+v.ClientID+" to active")
					sendClientList(activeConnections)
					pendingConnections.Remove(elem)
				}
			default:
				continue
			}
		}

		// Iterate active connections and check channels
		for _, clientConnection := range activeConnections {
			select {
			case message := <-clientConnection.MessageChan:
				if _, exists := activeConnections[message.To]; exists {
					activeConnections[message.To].MessageChan <- message
				} else {
					// TODO: handle message error
					log.Printf("Message error, invalid ID? from:%s, to:%s\n",
						message.From, message.To)
				}
			case status := <-clientConnection.StatusChan:
				if status == ConnectionLost {
					delete(activeConnections, clientConnection.ClientID)
					log.Print("Lost connection: " + clientConnection.ClientID + "\n")
					sendClientList(activeConnections)
				}
			default:
				continue
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %s", err.Error())
	}
}

func sendClientList(activeConnections map[string]ClientConnection) {
	// Make a sorted list of clients
	clientList := make([]string, len(activeConnections))
	i := 0
	for clientID, _ := range activeConnections {
		clientList[i] = clientID
		i = i + 1
	}
	sort.Strings(clientList)

	// Send list to everyone
	for cID, client := range activeConnections {
		select {
		case client.ListChan <- clientList:
		case <-time.After(10 * time.Millisecond):
			log.Printf("TIMEOUT; Skipped sending client list to %s\n", cID)
			continue
		}
	}
}

func pollTcpConnections(listener net.Listener, rawConnections chan net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			checkError(err)
		}
		rawConnections <- conn
	}
}

func NewClientConnection() *ClientConnection {
	c := &ClientConnection{}
	c.ListChan = make(chan []string)
	c.MessageChan = make(chan Message)
	c.StatusChan = make(chan StatusID)
	c.ClientID = "INVALID ID"
	return c
}

func (c *ClientConnection) Routine(conn net.Conn) {
	defer conn.Close()
	defer func() { c.StatusChan <- ConnectionLost }()

	reader := bufio.NewReader(conn)
	buf := make([]byte, 1024)
	var timeout <-chan time.Time

	for {
		if c.ClientID == "INVALID ID" { // Client is uninitialized
			n, err := conn.Read(buf)
			checkError(err)
			err = json.Unmarshal(buf[0:n], &c.ClientID)
			checkError(err)
			c.StatusChan <- UserIDRecieved
			timeout = time.After(10*time.Second)
		} else {
			select {
			case newList := <-c.ListChan:
				// New client list, send it
				message, err := json.Marshal(newList)
				checkError(err)
				_, err = conn.Write(message)
				checkError(err)
				// log.Println("Sent client list to " + c.ClientID)
			case message := <-c.MessageChan:
				// New message, send it
				jsonMessage, err := json.Marshal(message)
				checkError(err)
				_, err = conn.Write(jsonMessage)
				checkError(err)
				// log.Println("Sent message: ", message)
			case <-timeout:
				conn.Close();
				log.Printf("User %s timed out", c.ClientID);
				return
			default:
				// Read incoming TCP data from client & parse (should be a Message or keepAlive)
				conn.SetReadDeadline(time.Now())
				n, err := reader.Read(buf)
				opErr, isOpError := err.(*net.OpError)
				if err == nil && n > 0 {
					message := Message{}
					err = json.Unmarshal(buf[0:n], &message)
					if err != nil {
						var heartbeat KeepAliveMessage
						err = json.Unmarshal(buf[0:n], &heartbeat)
						if err != nil {
							log.Printf("Error parsing JSON from client: %s\n", err.Error())
						} else { // Update timeout
							timeout = time.After(2*time.Second)
						}
					} else {
						message.From = c.ClientID
						c.MessageChan <- message
					}
				} else if err != nil && !(isOpError && opErr.Timeout()) {
					conn.Close()
					return
				}
			}
		}
	}
}

