package types

type(
	MoveDir int
	OrderDir int
)

const(
	N_BUTTONS int = 3
	N_FLOORS  int = 4
	MAX_ORDERS int = 10
	N_ELEVATORS int = 2


	MAXSPEED int = 4048
	MINSPEED int = 2048

	ORDER_UP OrderDir = iota //matched with FLOOR for actuall order
	ORDER_DOWN
	ORDER_INTERNAL

	MOVE_UP MoveDir = iota //defines av elevators direction 
	MOVE_DOWN
	MOVE_STOP
)

type Elevator struct{
	//constant
	Ip string
	//subject to change (will trigger select)
	OrderQueue []Order
	Direction MoveDir
	LastFloor int
}

type Order struct{
	Floor int
	Orientation OrderDir
}

type Cost struct{
	Cost int
	Order Order
	Ip string
}

type Change struct{
	Type string
	Ip string
}

