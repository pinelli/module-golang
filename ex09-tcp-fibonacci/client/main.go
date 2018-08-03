package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"math/big"
	"time"
)

type Msg struct{
	n big.Int
}

func connect(addr string) net.Conn{
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Cannot send message: ", err)
		os.Exit(1)
	}
	return conn
}

func send(conn net.Conn, msg *big.Int) {
	encoder := json.NewEncoder(conn)
	err := encoder.Encode(&msg)
	if err != nil {
		fmt.Println("Error encoding message")
		os.Exit(1)
		conn.Close()
	}
}

func receive(conn net.Conn, msg *big.Int) {
	// encoder := json.NewEncoder(conn)
	// err := encoder.Encode(&msg)
	// if err != nil {
	// 	fmt.Println("Error encoding message")
	// 	os.Exit(1)
	// }

	//res *big.Int
	d := json.NewDecoder(conn)

	err := d.Decode(&msg)

	if err != nil {
		fmt.Println("Error decoding message")
		os.Exit(1)
	}
	fmt.Println("Message--------:", msg)
	a := msg
	fmt.Println("A:", a)

}

func main(){
	conn := connect("localhost:9090")
	// num := big.NewInt(5)

	// var res *big.Int = big.NewInt(0)
	// send(conn, num)
	// receive(conn, res)

	time.Sleep(200 * time.Millisecond)
	//-------------------------------------------------
	num1 := big.NewInt(100)
	var res1 *big.Int = big.NewInt(0)

	send(conn, num1)
	receive(conn, res1)

	fmt.Println("Res: ", res1)
	select{}
}
