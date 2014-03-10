package main

//import "net"
import "fmt"
//import "time"
import "encoding/json"

type KEEEG struct{

	I int
	Kek string
}
func main(){


	
	
	//test := network.TestVariable{1,"he"}
	
	test := KEEEG{3, "lol"}
	fmt.Println("Test før:", test)


	var svar KEEEG
	stringB,_ := json.Marshal(test)
	fmt.Println("Test i byte:", stringB)
	if err := json.Unmarshal(stringB, &svar); err != nil {
        panic(err)
    }

    fmt.Println(stringB, svar)


	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)
   	var dat map[string]interface{}
    if err := json.Unmarshal(byt, &dat); err != nil {
        panic(err)
    }



    /*
	udpAddr, err := net.ResolveUDPAddr("udp4", ":"+"20202")
	//CheckError(err, "ERROR while resolving UDPaddr for ListenToNetwork")
	fmt.Println("Establishing ListenToNetwork")
	conn, err := net.ListenUDP("udp4", udpAddr)
	fmt.Println("Listening on port ", udpAddr.String())
	//CheckError(err, "Error while establishing listening connection")

	conn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Millisecond))
	data := make([]byte, 1024)
	_,_,err = conn.ReadFromUDP(data)
	fmt.Println("Error says: ", err)
	kek := err.Error()
	port := "20202"
	fmt.Println( kek == "read udp4 0.0.0.0:"+port+": i/o timeout")
	*/
}