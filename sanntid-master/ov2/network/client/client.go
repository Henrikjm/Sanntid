package main

import (
	"time"
	"fmt"
	"net"
	"os"
	"bufio"
	"encoding/json"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s username\n", os.Args[0])
		os.Exit(1)
	}
	userName := os.Args[1]

	serverAddr := "129.241.187.146:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverAddr)
	checkError(err)

	clientList := make(ClientList, 50)

	messageChannel := make(chan Message,0)
	var outMessage, inMessage Message
	var imAlive KeepAliveMessage
	imAlive = 1337
	buf := make([]byte, 1024)
	go sendMessage(messageChannel)

	for {
		conn,err := dialServer(userName, tcpAddr)
		if(err != nil) {
			continue
		}

		for {
			if(conn == nil) {
				break;
			}
			select {
			case outMessage = <-messageChannel:
				//send message
				if validUser(clientList, outMessage.To){
					messageString,_ := json.Marshal(outMessage)
					_,err = conn.Write([]byte(messageString))
					checkError(err)
				}else{
					fmt.Printf("Invalid username: %s\n",outMessage.To)
				}
			case <-time.After(10*time.Millisecond):
				conn.SetReadDeadline(time.Now())
				n, err := conn.Read(buf)
				opErr, isOpError := err.(*net.OpError)
				if n > 0 {
					err = json.Unmarshal([]byte(buf[0:n]), &clientList)
					if err == nil {
						fmt.Println("Client list:")
						fmt.Println(strings.Join(clientList, "\n"))
					} else { //message was not client list, could be message
						//recieve message
						err = json.Unmarshal([]byte(buf[0:n]), &inMessage)
						checkError(err)
						fmt.Printf("Message from %s: %s\n", inMessage.From, inMessage.Message)
					}
				} else if err != nil && !(isOpError && opErr.Timeout()) {
					conn.Close()
					fmt.Println("Connection lost!")
					conn = nil
					break
				} else {
					imAliveString,_:= json.Marshal(imAlive)
					_,err = conn.Write([]byte(imAliveString))
					checkError(err)
				}
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func sendMessage(messageChannel chan Message){
	var outMessage Message
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		in,_ := input.ReadString('\n')
		userInput := strings.Split(in," ")

		if userInput[0] == "m" {
			outMessage.To = userInput[1]
			outMessage.Message = strings.Join(userInput[2:]," ")
			//TODO: handle nonvalid username
			messageChannel <- outMessage
			fmt.Println("sending message")
		} else {
			fmt.Println("Usage: m <reciever> <message>")
		}
	}
}

func validUser(clients ClientList, userName string) bool{
	for _,client := range(clients) {
		if client == userName {
			return true
		}
	}
	return false
}

func dialServer(userName string, tcpAddr *net.TCPAddr) (net.Conn, error) {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if(err != nil) {
		return conn, err
	}
	jsonUID, err := json.Marshal(userName)
	_,err = conn.Write(jsonUID)
	return conn,err
}
