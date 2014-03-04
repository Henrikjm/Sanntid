// channel workin in go

package main

import (
	"fmt"
)

func main() {

	var (
		i int = 1
	)
	channel1 := make(chan int, 1)

	channel1 <- i
	fmt.Println("int from channel = ", <-channel1)
	channel1 <- 3
	channel1 <- i

	fmt.Println("int from channel = ", <-channel1)
	fmt.Println("int from channel = ", <-channel1)

	channel1 <- 1

}
