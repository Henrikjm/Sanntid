package main

import(
	"fmt"
)

type Elevator struct{
	elevIp string
	workQueue []int
}


func main(){
	e1 := Elevator{"hello", []int{1,2,3,4}}
	e2 := Elevator{"yo", []int{2,3,4,5}}

	es := []Elevator{&e1,&e2}

	func(es []Elevator){
		var temp *Elevator
		temp = &es[0]
		&es[0] = &es[1]
		&es[1] = temp
	}(es)

	fmt.Println(&es)


}