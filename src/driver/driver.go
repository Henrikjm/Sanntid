package driver

import (
	"time"
	."types"
)


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
func SetOrderLight(order Order) {
	switch {
	case order.Floor == 1 && order.Orientation == ORDER_UP:
		SetBit(LIGHT_UP1)
	case order.Floor == 2 && order.Orientation== ORDER_UP:
		SetBit(LIGHT_UP2)
	case order.Floor == 3 && order.Orientation == ORDER_UP:
		SetBit(LIGHT_UP3)
	case order.Floor == 2 && order.Orientation == ORDER_DOWN:
		SetBit(LIGHT_DOWN2)
	case order.Floor == 3 && order.Orientation == ORDER_DOWN:
		SetBit(LIGHT_DOWN3)
	case order.Floor == 4 && order.Orientation == ORDER_DOWN:
		SetBit(LIGHT_DOWN4)
	case order.Floor == 1 && order.Orientation == ORDER_INTERNAL:
		SetBit(LIGHT_COMMAND1)
	case order.Floor == 2 && order.Orientation == ORDER_INTERNAL:
		SetBit(LIGHT_COMMAND2)
	case order.Floor == 3 && order.Orientation == ORDER_INTERNAL:
		SetBit(LIGHT_COMMAND3)
	case order.Floor == 4 && order.Orientation == ORDER_INTERNAL:
		SetBit(LIGHT_COMMAND4)
	}
}

func SetFloorIndicatorLight(floor int) {
	switch floor {
	case 1:
		ClearBit(FLOOR_IND1)
		ClearBit(FLOOR_IND2)
	case 2:
		ClearBit(FLOOR_IND1)
		SetBit(FLOOR_IND2)
	case 3:
		SetBit(FLOOR_IND1)
		ClearBit(FLOOR_IND2)
	case 4:
		SetBit(FLOOR_IND1)
		SetBit(FLOOR_IND2)
	}
}

//READS
func GetStopButton()  { ReadBit(STOP) }
func GetObstruction() { ReadBit(OBSTRUCTION) }

func GetOrderButton(localOrdersChan chan Order){
	for{
		time.Sleep(time.Millisecond*5)
		switch{
		case ReadBit(FLOOR_UP1):
			localOrdersChan <- Order{1,ORDER_UP}
		case ReadBit(FLOOR_UP2):
			localOrdersChan <- Order{2, ORDER_UP}
		case ReadBit(FLOOR_UP3):
			localOrdersChan <- Order{3, ORDER_UP}
		case ReadBit(FLOOR_DOWN2):
			localOrdersChan <- Order{2, ORDER_DOWN}
		case ReadBit(FLOOR_DOWN3):
			localOrdersChan <- Order{3, ORDER_DOWN}
		case ReadBit(FLOOR_DOWN4):
			localOrdersChan <- Order{4, ORDER_DOWN}
		case ReadBit(FLOOR_COMMAND1):
			localOrdersChan <- Order{1, ORDER_INTERNAL}
		case ReadBit(FLOOR_COMMAND2):
			localOrdersChan <- Order{2, ORDER_INTERNAL}
		case ReadBit(FLOOR_COMMAND3):
			localOrdersChan <- Order{3, ORDER_INTERNAL}
		case ReadBit(FLOOR_COMMAND4):
			localOrdersChan <- Order{4, ORDER_INTERNAL}
		}
	}




}

func ReadFloor()int{
	switch{
	case ReadBit(SENSOR1):
		return 1
	case ReadBit(SENSOR2):
		return 2
	case ReadBit(SENSOR3):
		return 3
	case ReadBit(SENSOR4):
		return 4
	}
	return 0
}



func MotorControl() {
	currentDir := MOVE_STOP
	WriteAnalog(MOTOR, MINSPEED)
	for {

		newDir := <-motorChannel

		if (newDir == MOVE_STOP) && (currentDir == MOVE_UP) {
			SetBit(MOTORDIR)
			time.Sleep(time.Millisecond * 20)
			WriteAnalog(MOTOR, MINSPEED)
			
		} else if (newDir == MOVE_STOP) && (currentDir == MOVE_DOWN) {
			ClearBit(MOTORDIR)
			time.Sleep(time.Millisecond * 20)
			WriteAnalog(MOTOR, MINSPEED)
			
		} else if newDir == MOVE_UP {
			ClearBit(MOTORDIR)
			WriteAnalog(MOTOR, MAXSPEED)
			time.Sleep(time.Second*1)
		} else if newDir == MOVE_DOWN {
			SetBit(MOTORDIR)
			WriteAnalog(MOTOR, MAXSPEED)
			time.Sleep(time.Second*1)
		} else {
			WriteAnalog(MOTOR, MINSPEED)
		}
		currentDir = newDir
		time.Sleep(time.Millisecond*5)
	}
}



//VARIABLES
var motorChannel chan MoveDir
var readFloorChannel chan int

