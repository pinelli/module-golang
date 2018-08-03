package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"time"
)

type Server struct {
	Host  string
	Cache map[int]*big.Int
}

type Response struct {
	N    *big.Int
	Time time.Duration
}

func CreateServer(host string) *Server {

	server := &Server{host, make(map[int]*big.Int)}
	return server
}

func (server *Server) Run() {
	go server.Listen()
}

func (server *Server) Send(num *big.Int, time time.Duration, conn net.Conn) {
	resp := Response{num, time}
	fmt.Printf("Sending...%#v", resp)

	encoder := json.NewEncoder(conn)
	err := encoder.Encode(resp)

	if err != nil {
		fmt.Println("Error encoding message")
		conn.Close()
		return
	}
}

func (server *Server) Listen() {
	ln, err := net.Listen("tcp", server.Host)
	if err != nil {
		fmt.Println("Cannot create server")
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Accepted connection:", conn.RemoteAddr())
		go server.handleConnection(conn)
	}
}

func (server *Server) handleConnection(c net.Conn) {
	d := json.NewDecoder(c)
	for {
		fmt.Println("iteration")
		var msg *big.Int

		err := d.Decode(&msg)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Client", c.RemoteAddr(), "disconnected")
			} else {
				fmt.Println("Error decoding message", err)
				c.Close()
			}
			return
		}

		val, _ := strconv.Atoi(msg.String())
		fmt.Println("Message:", msg)
		if val < 0 {
			server.Send(big.NewInt(-1), time.Duration(0), c)
		} else if cached, exist := server.Cache[val]; exist {
			server.Send(cached, time.Duration(0), c)
		} else {
			start := time.Now()
			fib := fibonacci(val)
			server.Cache[val] = fib
			fmt.Println("MyCache: ", server.Cache)
			elapsed := time.Since(start)
			server.Send(fib, elapsed, c)

		}
	}
}
