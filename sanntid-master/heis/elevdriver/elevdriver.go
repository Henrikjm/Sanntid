package elevdriver

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

type MoveDirection int

const (
	MOVE_UP MoveDirection = iota
	MOVE_DOWN
	MOVE_STOP
)

type OrderDirection int

const (
	ORDER_UP OrderDirection = iota
	ORDER_DOWN
	ORDER_INTERNAL
)

func (self MoveDirection) String() (s string) {
	switch self {
	case MOVE_UP:
		s = "move_up"
	case MOVE_DOWN:
		s = "move_down"
	case MOVE_STOP:
		s = "move internal"
	default:
		s = "move dir out of range"
	}
	return s
}
func (self OrderDirection) String() (s string) {
	switch self {
	case ORDER_UP:
		s = "order up"
	case ORDER_DOWN:
		s = "order down"
	case ORDER_INTERNAL:
		s = "order internal"
	default:
		s = "order dir out of range"
	}
	return s
}

const N_FLOORS = 4

type Button struct {
	Floor int
	Dir   OrderDirection
}
type OrderT Button

func (self Button) String() string {
	return fmt.Sprintf("Button %d,%s", self.Floor, self.Dir)
}

func (self OrderT) Increment() (dummy OrderT) {
	up := ORDER_UP
	down := ORDER_DOWN
	stop := ORDER_INTERNAL
	var priorityQueue [2][N_FLOORS]OrderT = [2][N_FLOORS]OrderT{
		[N_FLOORS]OrderT{{2, up}, {3, up}, {4, down}, {4, stop}},
		[N_FLOORS]OrderT{{1, stop}, {1, up}, {2, down}, {3, down}}}
	if (self.Dir == up) && (self.Floor == 4) {
		self.Dir = down
	} else if (self.Dir == down) && (self.Floor == 1) {
		self.Dir = up
	}

	if self.Dir == stop {
		panic("error in Increment: order out of range")
	} else {
		dummy = priorityQueue[self.Dir][self.Floor-1]
	}
	return dummy
}

const MAX_SPEED = 4024
const MIN_SPEED = 2048

func Init() (startFloor int, startDir MoveDirection) {
	val := IoInit()
	if !val {
		fmt.Printf("Driver initiated\n")
	} else {
		fmt.Printf("Driver not initiated\n")
	}

	clearLights()

	buttonChan = make(chan Button)
	floorChan = make(chan int)
	motorChan = make(chan MoveDirection)
	stopButtonChan = make(chan bool)
	obsChan = make(chan bool)

	go func() {
		// capture ctrl+c and stop elevator
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		sig := <-c
		log.Printf("captured %v, stopping elevator and exiting..", sig)
		Write_analog(MOTOR, MIN_SPEED)
		clearLights()
		os.Exit(1)
	}()

	go listen()
	go motorHandler()

	//ensure elevator is in valid state
	if Read_bit(SENSOR4) {
		startDir = MOVE_DOWN
		startFloor = 4
	} else {
		startDir = MOVE_UP
		MotorUp()
		startFloor = GetFloor() //blocking until floor reached
		MotorStop()
	}
	SetFloor(startFloor)
	return startFloor, startDir
}

func clearLights() {
	ClearDoor()
	ClearStopButton()
	ClearLight(1, ORDER_UP)
	ClearLight(2, ORDER_UP)
	ClearLight(3, ORDER_UP)
	ClearLight(2, ORDER_DOWN)
	ClearLight(3, ORDER_DOWN)
	ClearLight(4, ORDER_DOWN)
	ClearLight(1, ORDER_INTERNAL)
	ClearLight(2, ORDER_INTERNAL)
	ClearLight(3, ORDER_INTERNAL)
	ClearLight(4, ORDER_INTERNAL)
}

var buttonChan chan Button
var floorChan chan int
var motorChan chan MoveDirection
var stopButtonChan chan bool
var obsChan chan bool

func motorHandler() {
	currentDir := MOVE_STOP
	Write_analog(MOTOR, MIN_SPEED)
	for {
		newDir := <-motorChan
		if (newDir == MOVE_STOP) && (currentDir == MOVE_UP) {
			Set_bit(MOTORDIR)
			time.Sleep(10 * time.Millisecond) //break
			Write_analog(MOTOR, MIN_SPEED)
		} else if (newDir == MOVE_STOP) && (currentDir == MOVE_DOWN) {
			Clear_bit(MOTORDIR)
			time.Sleep(10 * time.Millisecond) //break
			Write_analog(MOTOR, MIN_SPEED)
		} else if newDir == MOVE_UP {
			Clear_bit(MOTORDIR)
			Write_analog(MOTOR, MAX_SPEED)
		} else if newDir == MOVE_DOWN {
			Set_bit(MOTORDIR)
			Write_analog(MOTOR, MAX_SPEED)
		} else {
			Write_analog(MOTOR, MIN_SPEED)
		}
		currentDir = newDir
	}
}

func listen() {
	var floorMap = map[int]int{
		SENSOR1: 1,
		SENSOR2: 2,
		SENSOR3: 3,
		SENSOR4: 4,
	}

	var buttonMap = map[int]Button{
		FLOOR_COMMAND1: {1, ORDER_INTERNAL},
		FLOOR_COMMAND2: {2, ORDER_INTERNAL},
		FLOOR_COMMAND3: {3, ORDER_INTERNAL},
		FLOOR_COMMAND4: {4, ORDER_INTERNAL},
		FLOOR_UP1:      {1, ORDER_UP},
		FLOOR_UP2:      {2, ORDER_UP},
		FLOOR_UP3:      {3, ORDER_UP},
		FLOOR_DOWN2:    {2, ORDER_DOWN},
		FLOOR_DOWN3:    {3, ORDER_DOWN},
		FLOOR_DOWN4:    {4, ORDER_DOWN},
	}

	buttonList := make(map[int]bool)
	for key, _ := range buttonMap {
		buttonList[key] = Read_bit(key)
	}

	floorList := make(map[int]bool)
	for key, _ := range floorMap {
		floorList[key] = Read_bit(key)
	}

	oldStop := false
	oldObs := false

	for {
		time.Sleep(1E7)
		for key, floor := range floorMap {
			newValue := Read_bit(key)
			if newValue != floorList[key] {
				newFloor := floor
				go func() {
					floorChan <- newFloor
				}()
			}
			floorList[key] = newValue
		}

		for key, btn := range buttonMap {
			newValue := Read_bit(key)
			if newValue && !buttonList[key] {
				newButton := btn
				go func() {
					buttonChan <- newButton
				}()
			}
			buttonList[key] = newValue
		}

		newStop := Read_bit(STOP)
		if newStop && !oldStop {
			go func() {
				stopButtonChan <- true
			}()
		}
		oldStop = newStop

		newObs := Read_bit(OBSTRUCTION)
		if newObs != oldObs {
			go func() {
				obsChan <- newObs
			}()
		}
		oldObs = newObs
	}

}

func SetLight(floor int, dir OrderDirection) {
	switch {
	case floor == 1 && dir == ORDER_INTERNAL:
		Set_bit(LIGHT_COMMAND1)
	case floor == 2 && dir == ORDER_INTERNAL:
		Set_bit(LIGHT_COMMAND2)
	case floor == 3 && dir == ORDER_INTERNAL:
		Set_bit(LIGHT_COMMAND3)
	case floor == 4 && dir == ORDER_INTERNAL:
		Set_bit(LIGHT_COMMAND4)
	case floor == 1 && dir == ORDER_UP:
		Set_bit(LIGHT_UP1)
	case floor == 2 && dir == ORDER_UP:
		Set_bit(LIGHT_UP2)
	case floor == 3 && dir == ORDER_UP:
		Set_bit(LIGHT_UP3)
	case floor == 2 && dir == ORDER_DOWN:
		Set_bit(LIGHT_DOWN2)
	case floor == 3 && dir == ORDER_DOWN:
		Set_bit(LIGHT_DOWN3)
	case floor == 4 && dir == ORDER_DOWN:
		Set_bit(LIGHT_DOWN4)
	}
}

func ClearLight(floor int, dir OrderDirection) {
	switch {
	case floor == 1 && dir == ORDER_INTERNAL:
		Clear_bit(LIGHT_COMMAND1)
	case floor == 2 && dir == ORDER_INTERNAL:
		Clear_bit(LIGHT_COMMAND2)
	case floor == 3 && dir == ORDER_INTERNAL:
		Clear_bit(LIGHT_COMMAND3)
	case floor == 4 && dir == ORDER_INTERNAL:
		Clear_bit(LIGHT_COMMAND4)
	case floor == 1 && dir == ORDER_UP:
		Clear_bit(LIGHT_UP1)
	case floor == 2 && dir == ORDER_UP:
		Clear_bit(LIGHT_UP2)
	case floor == 3 && dir == ORDER_UP:
		Clear_bit(LIGHT_UP3)
	case floor == 2 && dir == ORDER_DOWN:
		Clear_bit(LIGHT_DOWN2)
	case floor == 3 && dir == ORDER_DOWN:
		Clear_bit(LIGHT_DOWN3)
	case floor == 4 && dir == ORDER_DOWN:
		Clear_bit(LIGHT_DOWN4)
	}
}

func MotorUp() {
	motorChan <- MOVE_UP
}

func MotorDown() {
	motorChan <- MOVE_DOWN
}

func MotorStop() {
	motorChan <- MOVE_STOP
}

func GetButton() (int, OrderDirection) {
	btn := <-buttonChan
	return btn.Floor, btn.Dir
}

func GetButtonChan() chan Button {
	return buttonChan
}

func GetFloor() int {
	floor := <-floorChan
	return floor
}

func GetFloorChan() chan int {
	return floorChan
}

func SetFloor(floor int) {
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

func GetStopButton() {
	<-stopButtonChan
}

func GetStopButtonChan() chan bool {
	return stopButtonChan
}

func SetStopButton() {
	Set_bit(LIGHT_STOP)
}

func ClearStopButton() {
	Clear_bit(LIGHT_STOP)
}

func GetObs() bool {
	return <-obsChan
}

func GetObsChan() chan bool {
	return obsChan
}

func SetDoor() {
	Set_bit(DOOR_OPEN)
}

func ClearDoor() {
	Clear_bit(DOOR_OPEN)
}

func GetDoor() bool {
	return Read_bit(DOOR_OPEN)
}
