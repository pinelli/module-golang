package main

import "fmt"
import "time"
import "os"

type Client struct {
	name  string
	shop  BarberShop
	chair Chair //chair that the client took
}

func NewClient(name string, shop BarberShop) Client {
	return Client{name, shop, Chair{}}
}

func (this *Client) waitForBarber() {
	dialog := make(chan bool)

	fmt.Println(this.name, "is waiting for barber")
	d := Dialog{this, dialog}
	this.shop.queue <- d

	<-dialog //barber: accepted
	this.chair.release()
	dialog <- true //client: chair released

	<-dialog //barber: shaving is finished
}

func (this *Client) takeChair() bool {
	for _, ch := range this.shop.chairs {
		if ch.take() {
			fmt.Fprintln(os.Stderr, this.name, "took", ch.id, "chair")
			this.chair = ch
			return true
		}
	}
	return false
}

func (this *Client) visit() {
	for true {
		if this.takeChair() {
			this.waitForBarber()
			//this.chair.release()
			fmt.Fprintln(os.Stderr, "Client", this.name, "is leaving")
		} else {
			fmt.Fprintln(os.Stderr, "No place for client:", this.name)
		}
		time.Sleep(200 * time.Millisecond)
	}

}
