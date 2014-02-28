package main

import (

	."fmt"
	"os/exec"
)

func main() {
	// i := 0
	// select{
	// case i = 0:
	// 	Println("check")
	// default:
	// 	Println("poopcheck")
	// }
	i := 0
	if i == 0{
		Println("kek")
	}else{
		Println("vet ikke")
	}
	cmd := exec.Command("mate-terminal", "-x")
	cmd.Run()
}
