package main

import(
	"fmt"
	"time"
)

func count(){
	i := 0
	for{
		i += 1
		time.Sleep(time.Millisecond * 100)
		fmt.Println(i)
	}
}


func main(){
	
	var inc, out chan string
	var i, o string

	go count()

	for{
		select{
		case i = <- inc:
			fmt.Println(i)
		case o = <- out:
			fmt.Println(o)
		}
	}
}