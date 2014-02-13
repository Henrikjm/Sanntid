package main

import(
	"github.com/anlif/sanntid/heis/elevcontrol"
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Started!\n")
	var elev elevcontrol.Elevator
	elev.Start()
	for {
		time.Sleep(10*time.Second)
	}
}
