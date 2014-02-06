package main

import(
	."fmt"
	."net"
)

var err error
var chanCon chan *UDPConn

func ListenToNetwork(chanCon chan *UDPConn){
	Println("Start UDP server")

	udpAddr, err := ResolveUDPAddr("udp4", ":20666") //resolving
	if err != nil{
		println("ERROR: Resolving error")
	}
	conn, err := ListenUDP("udp", udpAddr) //initiating listening
	if err != nil{
		println("ERROR: Listening error")
	}
   chanCon <- conn
	data := make([]byte,1024)
	for{
		_, addr, err := conn.ReadFromUDP(data) //kan bruke addr til å sjekke hvor melding kommer fra f.eks if addr not [egen i.p]
		if err != nil{
			println("ERROR: while reading")
		}
		Println("Recieved from: ", addr,"\nMessage: ",string(data))
	}	
}

func SendToNetwork(chanCon chan *UDPConn){
   sendAddr, err := ResolveUDPAddr("udp4","129.241.187.255:20018") //Spesifiserer adresse
	//connection,err := DialUDP("udp",nil, sendAddr) //setter opp "socket" for sending
	if err != nil {
		println("ERROR while resolving UDP addr")
	}
	connection := <- chanCon
	for{
		var(inn string)
		_,_ = Scanf("&d",inn)
		testmsg := []byte("testing")
		if connection ==  nil{
			println("ERROR, connection = nil")
		}	
		connection.WriteToUDP(testmsg, sendAddr)
	}	
}



func main(){

var(shit string)
for{
	Scanf(shit)
	println(shit)
}
/*chanCon := make(chan *UDPconn, 1)
go ListenToNetwork(chanCon)
go SendToNetwork(chanCon)
*/

}