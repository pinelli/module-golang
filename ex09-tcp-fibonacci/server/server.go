package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"math/big"
)

type Server struct {
	Host    string
}

func CreateServer(host string) *Server {
	server := &Server{host}
	return server
}

func (server *Server) Run() {
	go server.Listen()
}

func (server *Server) Send(num *big.Int, time int, conn net.Conn) {
	fmt.Println("Sending...", num)
	encoder := json.NewEncoder(conn)
	err := encoder.Encode(&num)

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
		if err.Error() == "EOF"{
			fmt.Println("Client" ,c.RemoteAddr(), "disconnected")
		}else{
			fmt.Println("Error decoding message", err)
			c.Close()	
		}
		return
	}
	fmt.Println("Message:", msg)
	server.Send(msg, 5, c)
	}

	//c.Close()
}