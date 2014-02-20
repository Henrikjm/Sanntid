package driver

import (
	"time"
)

//TYPE DEFINES
type OrderDir int
type MoveDir int

//VARIABLES
var motorChannel chan MoveDir

//CONSTANTS
const (
	N_BUTTONS int = 3
	N_FLOORS  int = 4

	MAXSPEED int = 4048
	MINSPEED int = 2048

	ORDER_UP OrderDir = iota
	ORDER_DOWN
	ORDER_INTERNAL

	MOVE_UP MoveDir = iota
	MOVE_DOWN
	MOVE_STOP
)

func InitElev() {
	ClearAllLights()
	motorChannel = make(chan MoveDir)
	go MotorControl()
	for {
		motorChannel <- MOVE_DOWN
		time.Sleep(time.Second * 1)
		motorChannel <- MOVE_UP
		time.Sleep(time.Second * 1)
		motorChannel <- MOVE_STOP
		time.Sleep(time.Second * 1)
	}
}

func MotorControl() {
	currentDir := MOVE_STOP
	WriteAnalog(MOTOR, MINSPEED)
	for {
		newDir := <-motorChannel

		if (newDir == MOVE_STOP) && (currentDir == MOVE_UP) {
			SetBit(MOTORDIR)
			time.Sleep(time.Millisecond * 10)
			WriteAnalog(MOTOR, MINSPEED)
		} else if (newDir == MOVE_STOP) && (currenDir == MOVE_DOWN) {
			ClearBit(MOTORDIR)
			time.Sleep(time.Millisecond * 10)
			WriteAnalog(MOTOR, MINSPEED)
		} else if newDir == MOVE_UP {
			ClearBit(MOTORDIR)
			WriteAnalog(MOTOR, MAXSPEED)
		} else if newDir == MOVE_DOWN {
			SetBit(MOTORDIR)
			WriteAnalog(MOTOR, MAXSPEED)
		} else {
			WriteAnalog(MOTOR, MINSPEED)
		}
		currentDir = newDir
	}
}

//LIGHTS

func ClearAllLights() {
	ClearAllOrderLights()
	ClearStopLight()
	ClearDoorOpenLight()
}
func ClearDoorOpenLight() { ClearBit(DOOR_OPEN) }
func ClearStopLight()     { ClearBit(LIGHT_STOP) }
func ClearOrderLight(floor int, dir OrderDir) {
	switch {
	case floor == 1 && dir == ORDER_UP:
		ClearBit(LIGHT_UP1)
	case floor == 2 && dir == ORDER_UP:
		ClearBit(LIGHT_UP2)
	case floor == 3 && dir == ORDER_UP:
		ClearBit(LIGHT_UP3)
	case floor == 2 && dir == ORDER_DOWN:
		ClearBit(LIGHT_DOWN2)
	case floor == 3 && dir == ORDER_DOWN:
		ClearBit(LIGHT_DOWN3)
	case floor == 4 && dir == ORDER_DOWN:
		ClearBit(LIGHT_DOWN4)
	case floor == 1 && dir == ORDER_INTERNAL:
		ClearBit(LIGHT_COMMAND1)
	case floor == 2 && dir == ORDER_INTERNAL:
		ClearBit(LIGHT_COMMAND2)
	case floor == 3 && dir == ORDER_INTERNAL:
		ClearBit(LIGHT_COMMAND3)
	case floor == 4 && dir == ORDER_INTERNAL:
		ClearBit(LIGHT_COMMAND4)
	}
}
func ClearAllOrderLights() {
	ClearOrderLight(1, ORDER_UP)
	ClearOrderLight(2, ORDER_UP)
	ClearOrderLight(3, ORDER_UP)
	ClearOrderLight(2, ORDER_DOWN)
	ClearOrderLight(3, ORDER_DOWN)
	ClearOrderLight(4, ORDER_DOWN)
	ClearOrderLight(1, ORDER_INTERNAL)
	ClearOrderLight(1, ORDER_INTERNAL)
	ClearOrderLight(1, ORDER_INTERNAL)
	ClearOrderLight(1, ORDER_INTERNAL)
}

func SetDoorOpenLight() { SetBit(DOOR_OPEN) }
func SetStopLight()     { SetBit(LIGHT_STOP) }
func SetOrderLight() {
	switch {
	case floor == 1 && dir == ORDER_UP:
		SetBit(LIGHT_UP1)
	case floor == 2 && dir == ORDER_UP:
		SetBit(LIGHT_UP2)
	case floor == 3 && dir == ORDER_UP:
		SetBit(LIGHT_UP3)
	case floor == 2 && dir == ORDER_DOWN:
		SetBit(LIGHT_DOWN2)
	case floor == 3 && dir == ORDER_DOWN:
		SetBit(LIGHT_DOWN3)
	case floor == 4 && dir == ORDER_DOWN:
		SetBit(LIGHT_DOWN4)
	case floor == 1 && dir == ORDER_INTERNAL:
		SetBit(LIGHT_COMMAND1)
	case floor == 2 && dir == ORDER_INTERNAL:
		SetBit(LIGHT_COMMAND2)
	case floor == 3 && dir == ORDER_INTERNAL:
		SetBit(LIGHT_COMMAND3)
	case floor == 4 && dir == ORDER_INTERNAL:
		SetBit(LIGHT_COMMAND4)
	}
}
func SetFloorIndicatorLight(floor int) {
	switch floor {
	case 1:
		Clear_bit(FLOOR_IND1)
		Clear_bit(FLOOR_IND2)
	case 2:
		Clear_bit(FLOOR_IND1)
		Set_bit(FLOOR_IND2)
	case 3:
		Set_bit(FLOOR_IND1)
		Clear_bit(FLOOR_IND2)
	case 4:
		Set_bit(FLOOR_IND1)
		Set_bit(FLOOR_IND2)
	}
}

//READS
func GetStopButton() int  { return ReadBit(STOP) }
func GetObstruction() int { return ReadBit(OBSTRUCTION) }
func GetOrderButton() int {

}
func ReadFloor() int {

}
