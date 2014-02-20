package main

import (
	"driver"
	"time"
)

func main() {
	driver.IoInit()
	driver.SetBit(MOTORDIR)
	driver.WriteAnalog(MOTOR, 100)
	time.Sleep(2 * time.Second)
	driver.WriteAnalog(MOTOR, -100)
	time.Sleep(2 * time.Second)
	driver.WriteAnalog(MOTOR, 0)
	driver.ClearBit(MOTORDIR)
}
