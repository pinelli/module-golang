package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

type Response struct {
	N    *big.Int
	Time time.Duration
}

func connect(addr string) net.Conn {
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

func receive(conn net.Conn) (*big.Int, time.Duration) {
	d := json.NewDecoder(conn)
	var resp Response
	err := d.Decode(&resp)

	if err != nil {
		//fmt.Println("Error decoding message")
		os.Exit(1)
	}
	return resp.N, resp.Time
}

func reader(jobs chan int64) {

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		txt := scanner.Text()

		var job int64
		fmt.Sscanf(txt, "%d\n", &job)

		jobs <- job
	}
	close(jobs)
}

func main() {
	conn := connect("localhost:9090")

	var jobs = make(chan int64)
	go reader(jobs)

	for job := range jobs {
		send(conn, big.NewInt(job))
		n, t := receive(conn)

		fmt.Fprintln(os.Stdout, t, n)
	}
}
