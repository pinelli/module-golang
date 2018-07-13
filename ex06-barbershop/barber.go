package main

import "fmt"
import "time"
import "os"

func Shave(name string) {
	fmt.Fprintln(os.Stderr, "Barber shaves:", name)
}

func BarberWork(shop *BarberShop) {
	for true {
		select {
		case client := <-(*shop).queue:
			client.dialog <- true //barber: accept
			<-client.dialog       //client: chair released

			Shave(client.client.name)
			client.dialog <- true //shaving is finished
		default:
			fmt.Fprintln(os.Stderr, "Barber sleeps....Zzz.....")
			time.Sleep(200 * time.Millisecond)
		}
	}
}
