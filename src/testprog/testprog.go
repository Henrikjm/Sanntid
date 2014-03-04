package main

import (
	//"network"
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

func main() {
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
	type Queue struct {
		IP string
		timeRef time.Time
	}

	
	var workingVariable Queue
	
	workingVariable.IP = *GetLocalIp()
	workingVariable.timeRef = time.Now()
	fmt.Println("TimeNow = ", workingVariable.timeRef, "IP = ",workingVariable.IP)
	
	
}
