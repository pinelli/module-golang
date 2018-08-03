package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

type Server struct {
	Host string
}

type Response struct {
	N    *big.Int
	Time time.Duration
}

func CreateServer(host string) *Server {
	server := &Server{host}
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
		fmt.Println("Message:", msg)
		server.Send(msg, time.Duration(5), c)
	}

	//c.Close()
}
