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
func ClearInternalOrderLight(order Order){
		switch {
		case order.Floor == 1 && order.Orientation == ORDER_INTERNAL:
			ClearBit(LIGHT_COMMAND1)
		case order.Floor == 2 && order.Orientation == ORDER_INTERNAL:
			ClearBit(LIGHT_COMMAND2)
		case order.Floor == 3 && order.Orientation == ORDER_INTERNAL:
			ClearBit(LIGHT_COMMAND3)
		case order.Floor == 4 && order.Orientation == ORDER_INTERNAL:
			ClearBit(LIGHT_COMMAND4)
		}
}
func ClearAllOrderLights() {
	ClearBit(LIGHT_UP1)
	ClearBit(LIGHT_UP2)
	ClearBit(LIGHT_UP3)
	ClearBit(LIGHT_DOWN2)
	ClearBit(LIGHT_DOWN3)
	ClearBit(LIGHT_DOWN4)
	ClearBit(LIGHT_COMMAND1)
	ClearBit(LIGHT_COMMAND2)
	ClearBit(LIGHT_COMMAND3)
	ClearBit(LIGHT_COMMAND4)
}
func SetDoorOpenLight() { SetBit(DOOR_OPEN) }
func SetStopLight()     { SetBit(LIGHT_STOP) }
func SetInternalOrderLights(){
	for{
		orders := <- setInternalOrderLightChannel
		for i:= 0; i < len(orders); i++{
			switch {
			case orders[i].Floor == 1 && orders[i].Orientation == ORDER_INTERNAL:
				SetBit(LIGHT_COMMAND1)
			case orders[i].Floor == 2 && orders[i].Orientation == ORDER_INTERNAL:
				SetBit(LIGHT_COMMAND2)
			case orders[i].Floor == 3 && orders[i].Orientation == ORDER_INTERNAL:
				SetBit(LIGHT_COMMAND3)
			case orders[i].Floor == 4 && orders[i].Orientation == ORDER_INTERNAL:
				SetBit(LIGHT_COMMAND4)
			}
		}
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

func GetStopButton(){
	for{
		stopButtonChannel <- ReadBit(STOP)
	}
}
func GetObstruction() { ReadBit(OBSTRUCTION) }
func GetOrderButton(localOrderChan chan Order){
	for{
		time.Sleep(time.Millisecond*30)
		switch{
		case ReadBit(FLOOR_UP1):
			localOrderChan <- Order{1,ORDER_UP}
		case ReadBit(FLOOR_UP2):
			localOrderChan <- Order{2, ORDER_UP}
		case ReadBit(FLOOR_UP3):
			localOrderChan <- Order{3, ORDER_UP}
		case ReadBit(FLOOR_DOWN2):
			localOrderChan <- Order{2, ORDER_DOWN}
		case ReadBit(FLOOR_DOWN3):
			localOrderChan <- Order{3, ORDER_DOWN}
		case ReadBit(FLOOR_DOWN4):
			localOrderChan <- Order{4, ORDER_DOWN}
		case ReadBit(FLOOR_COMMAND1):
			localOrderChan <- Order{1, ORDER_INTERNAL}
		case ReadBit(FLOOR_COMMAND2):
			localOrderChan <- Order{2, ORDER_INTERNAL}
		case ReadBit(FLOOR_COMMAND3):
			localOrderChan <- Order{3, ORDER_INTERNAL}
		case ReadBit(FLOOR_COMMAND4):
			localOrderChan <- Order{4, ORDER_INTERNAL}
		}
	}
}
func ReadFloor()int{
	time.Sleep(time.Millisecond * 5)
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
		newDir = <-motorChannel
		if (newDir == MOVE_STOP) && (currentDir == MOVE_UP) {
			WriteAnalog(MOTOR, 1000)
			SetBit(MOTORDIR)
			time.Sleep(time.Millisecond * 30)
			ClearBit(MOTORDIR)
			WriteAnalog(MOTOR, MINSPEED)
			currentDir = newDir
			
		} else if (newDir == MOVE_STOP) && (currentDir == MOVE_DOWN) {
			WriteAnalog(MOTOR, 1000)
			ClearBit(MOTORDIR)
			time.Sleep(time.Millisecond * 30)
			SetBit(MOTORDIR)
			WriteAnalog(MOTOR, MINSPEED)
			currentDir = newDir
			
		} else if newDir == MOVE_UP {
			ClearBit(MOTORDIR)
			WriteAnalog(MOTOR, MAXSPEED)
			currentDir = newDir
			
		} else if newDir == MOVE_DOWN {
			SetBit(MOTORDIR)
			WriteAnalog(MOTOR, MAXSPEED)
			currentDir = newDir
			
		} else {
			WriteAnalog(MOTOR, MINSPEED)
			currentDir = newDir
		}
	}
}



//VARIABLES
var newDir MoveDir
var motorChannel chan MoveDir
var readFloorChannel chan int
var setInternalOrderLightChannel chan []Order
var stopButtonChannel chan bool

