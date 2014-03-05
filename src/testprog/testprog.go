package main

import (
	"network"
	//"driver"
	"fmt"
	"time"
	"net"
	"strings"
)

func GetLocalIp() *string {
	conn, _ := net.Dial("udp4", "google.com:80")
	//CheckError(err, "ERROR: LocalIp: dialing to google.com:80")
	return &strings.Split(conn.LocalAddr().String(), ":")[0]
}

func RecieveAlive(alivePort string, aliveChan *chan map[string]time.Time){
	data := make([]byte, 1024)
	ownAddr := *GetLocalIp();
	conn := network.MakeListenerConn(alivePort)
	for {		
		_, addr_, err := conn.ReadFromUDP(data)
		CheckError(err, "ERROR ReadFromUDP")

		if (string(data) == "ImAlive!") && (addr.String() != ownAddr){
			*aliveChan <- addr.String()//add/update alive map
		}
	}

}

func UpdateAlive(aliveChan *chan string, updateChan *chan map[string]time.Time) {
	for {
		select{
			case incomingIP := <-*aliveChan:
				aliveMap[incomingIP] = time.Now()
			case <-updateChan:
				*updateChan<-aliveMap
			default:
				for i, value := range aliveMap {//Iterate through alive-map and delete timed-out machines
					if time.Now().Sub(value) > 500000000 {
						delete(aliveMap, i)
					}
				}
				if lengthOfMap != len(aliveMap) {
					lengthOfMap = len(aliveMap)
					*updateChan <- aliveMap
				}
		}
	}
}

func main() {
	const N_ELEVATORS int = 4
	/*go func(){
	driver.IoInit()
	driver.SetBit(driver.LIGHT_STOP)
	driver.SetBit(driver.LIGHT_COMMAND1)
	driver.SetBit(driver.LIGHT_DOWN2)
	driver.SetBit(driver.MOTORDIR)
	driver.WriteAnalog(driver.MOTOR, 4000)
	fmt.Println("Going to sleep.")
	time.Sleep(1 * time.Second)
	fmt.Println("Waking.")
	driver.ClearBit(driver.MOTORDIR)
	driver.WriteAnalog(driver.MOTOR, 4000)
	fmt.Println("Going to sleep.")
	time.Sleep(1 * time.Second)
	fmt.Println("Waking.")
	driver.WriteAnalog(driver.MOTOR, 0)
	driver.ClearBit(driver.MOTORDIR)

	driver.ClearBit(driver.LIGHT_STOP)
	driver.ClearBit(driver.LIGHT_COMMAND1)
	driver.ClearBit(driver.LIGHT_DOWN2)

	driver.IoInit()
	}()
	*/
	/*alivePort := "33042"
	incoming := make(chan string)

	go network.ImAliveUDP(alivePort)
	go network.ListenToNetworkUDP(network.MakeListenerConn(alivePort), incoming)

	for{
		fmt.Println(<-incoming, "bgerboeg")
		}
	var exit string
	fmt.Scanln(&exit)
	*/

	aliveArray := make(map[string]time.Time)
	IP := *GetLocalIp()
	aliveArray[IP] = time.Now()
	if n,ok := aliveArray[IP]; ok{
		fmt.Println(time.Now().Sub(n))
	}
	time.Sleep(time.Second * 1)
	fmt.Println(time.Now().Sub(aliveArray[*GetLocalIp()]))
	fmt.Println(len(aliveArray))
	//aliveArray := make([]Alive)


/*
	workingVariable.timeRef = time.Now()
	time.Sleep(time.Millisecond * 490)
	
	reference := time.Now()

	fmt.Println(reference.Sub(workingVariable.timeRef), 
		reference.Sub(workingVariable.timeRef) > 500000000)

	for i := 0; i < len(aliveArray); i++ {
		fmt.Println()
	}



	aliveArray[0] = workingVariable
	aliveArray[1] = keke
	fmt.Println(aliveArray[0].timeRef)

	fmt.Println(aliveArray[1].timeRef.Minute, aliveArray[1].timeRef.Date)
	fmt.Println(len(aliveArray))

	*/
}