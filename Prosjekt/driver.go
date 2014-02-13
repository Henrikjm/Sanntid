package driver

import (
	"time"
)

func func main() {
	IoInit()
	SetBit(MOTORDIR)
	WriteAnalog(MOTOR, 100)
	sleep(2*time.Second)
	WriteAnalog(Motor,-100)
	sleep(2*time.Second)
	WriteAnalog(Motor,0)
	ClearBit(MOTORDIR)
}