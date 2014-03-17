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
func ClearOrderLight(){
	for{
		order := <- clearOrderLightChannel
		//fmt.Println(order)
		switch {
		case order.Floor == 1 && order.Orientation == ORDER_UP:
			ClearBit(LIGHT_UP1)
		case order.Floor == 2 && order.Orientation == ORDER_UP:
			ClearBit(LIGHT_UP2)
		case order.Floor == 3 && order.Orientation == ORDER_UP:
			ClearBit(LIGHT_UP3)
		case order.Floor == 2 && order.Orientation == ORDER_DOWN:
			ClearBit(LIGHT_DOWN2)
		case order.Floor == 3 && order.Orientation == ORDER_DOWN:
			ClearBit(LIGHT_DOWN3)
		case order.Floor == 4 && order.Orientation == ORDER_DOWN:
			ClearBit(LIGHT_DOWN4)
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
}

func ClearAllOrderLights() {
	clearOrderLightChannel <- Order{1, ORDER_UP}
	clearOrderLightChannel <- Order{2, ORDER_UP}
	clearOrderLightChannel <- Order{3, ORDER_UP}
	clearOrderLightChannel <- Order{2, ORDER_DOWN}
	clearOrderLightChannel <- Order{3, ORDER_DOWN}
	clearOrderLightChannel <- Order{4, ORDER_DOWN}
	clearOrderLightChannel <- Order{1, ORDER_INTERNAL}
	clearOrderLightChannel <- Order{1, ORDER_INTERNAL}
	clearOrderLightChannel <- Order{1, ORDER_INTERNAL}
	clearOrderLightChannel <- Order{1, ORDER_INTERNAL}
}

func SetDoorOpenLight() { SetBit(DOOR_OPEN) }
func SetStopLight()     { SetBit(LIGHT_STOP) }
func SetOrderLights(){
	for{
		orders := <- setOrderLightChannel
		for i:= 0; i < MAX_ORDERS; i++{
			switch {
			case orders[i].Floor == 1 && orders[i].Orientation == ORDER_UP:
				SetBit(LIGHT_UP1)
			case orders[i].Floor == 2 && orders[i].Orientation== ORDER_UP:
				SetBit(LIGHT_UP2)
			case orders[i].Floor == 3 && orders[i].Orientation == ORDER_UP:
				SetBit(LIGHT_UP3)
			case orders[i].Floor == 2 && orders[i].Orientation == ORDER_DOWN:
				SetBit(LIGHT_DOWN2)
			case orders[i].Floor == 3 && orders[i].Orientation == ORDER_DOWN:
				SetBit(LIGHT_DOWN3)
			case orders[i].Floor == 4 && orders[i].Orientation == ORDER_DOWN:
				SetBit(LIGHT_DOWN4)
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

//READS
func GetStopButton(){
	for{
		stopButtonChannel <- ReadBit(STOP)
	}
}
func GetObstruction() { ReadBit(OBSTRUCTION) }

func GetOrderButton(localOrderChan chan Order){
	for{
		time.Sleep(time.Millisecond*1)
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
			
		} else if newDir == MOVE_DOWN {
			SetBit(MOTORDIR)
			WriteAnalog(MOTOR, MAXSPEED)
			
		} else {
			WriteAnalog(MOTOR, MINSPEED)
		}
	}
}



//VARIABLES
var newDir MoveDir
var motorChannel chan MoveDir
var readFloorChannel chan int
var setOrderLightChannel chan []Order
var clearOrderLightChannel chan Order
var stopButtonChannel chan bool

