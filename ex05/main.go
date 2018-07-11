package main

import "fmt"

type Dialog struct {
     client *Client
     dialog chan bool
}
type BarberShop struct{
  chairs []Chair
  queue chan Dialog
  barberAccepts chan bool
}

func main(){
	shop:=BarberShop{}
	shop.queue = make(chan Dialog, 2) 
  shop.barberAccepts = make(chan bool, 2)

	shop.chairs = append(shop.chairs, NewChair(0))
	shop.chairs = append(shop.chairs, NewChair(1))

	client1:=NewClient("Sasha", shop)
	client2:=NewClient("Vasia", shop)
	client3:=NewClient("Andrey", shop)

	go BarberWork(&shop)
	go client1.visit()
  go client2.visit()
  go client3.visit()

	select{}
	fmt.Print("OK")
}
