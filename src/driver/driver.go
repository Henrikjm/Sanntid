package driver

import (
	"time"
)

func main() {
	IoInit()
	SetBit(MOTORDIR)
	WriteAnalog(MOTOR, 100)
	time.Sleep(2 * time.Second)
	WriteAnalog(MOTOR, -100)
	time.Sleep(2 * time.Second)
	WriteAnalog(MOTOR, 0)
	ClearBit(MOTORDIR)
}
