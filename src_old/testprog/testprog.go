package main

import (
	
	//"network"
	// "driver"
	 "fmt"
	// "time"
	"encoding/json"

)


type KEEEG struct{

	i int = 0
	kek string = "lol"
}


func main() {
	
	
	//test := network.TestVariable{1,"he"}
	
	var test KEEEG
	stringB,_ := json.Marshal(test)
	fmt.Println(stringB, test)

}
	/*
	const N_ELEVATORS int = 4
	  type Response1 struct {
    	Page   int
	    Fruits []string
}
	type Response2 struct {
    	Page   int      `json:"page"`
    	Fruits []string `json:"fruits"`
	}
	res1D := &Response1{
        Page:   1,
        Fruits: []string{"apple", "peach", "pear"}}
    res1B, _ := json.Marshal(res1D)
    fmt.Println(string(res1B))*/
  
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

/*
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
	*/

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
