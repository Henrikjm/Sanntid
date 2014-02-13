package main

import (
	"time"
)

func main() {
	IoInit()
	SetBit(MOTORDIR)
	WriteAnalog(MOTOR, 100)
	time.sleep(2 * time.Second)
	WriteAnalog(Motor, -100)
	time.sleep(2 * time.Second)
	WriteAnalog(Motor, 0)
	ClearBit(MOTORDIR)
}
