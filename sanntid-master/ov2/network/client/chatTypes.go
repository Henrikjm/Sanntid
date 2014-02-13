package main

type KeepAliveMessage int

type ClientConnection struct {
	ListChan chan ClientList
	MessageChan chan Message
	StatusChan chan StatusID
	ClientID string
}

type ClientList []string

type Message struct {
	To string
	From string
	Message string
}

type Registration struct {
	ClientID string
}

type StatusID int
const (
	ConnectionLost = iota
	UserIDRecieved
)
