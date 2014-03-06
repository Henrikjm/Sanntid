package main

import(
	"fmt"
)

type Coin struct{
	Coin int
	Dollar int
	Order Order
}

type Order struct{
	hei string
	hopp string
}

func main(){
	
	coin := Coin{}
	fmt.Println(coin)

}