package main

import (
	"fmt"
	"net"
	"os"
)

const commandLen = 1
const keyLen = 3
const valLen = 3

type Server struct {
	Host string
}

func CreateServer(host string) *Server {

	server := &Server{host}
	return server
}

func (server *Server) Run() {
	go server.Listen()
}

/*
func (server *Server) Send(runes []rune, conn net.Conn) {
	if err != nil {
		fmt.Println("Error encoding message")
		conn.Close()
		return
	}
}
*/

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

func readBytes(conn net.Conn) (string, error) {
	res := make([]byte, 25, 25)
	for i := 0; i < 3; i++ {
		_, err := conn.Read(res[i : i+1]) //_, err
		if err != nil {
			return "", err
		}
	}
	str := string(res)
	return str, nil
}

func (server *Server) handleConnection(conn net.Conn) {
	str, err := readBytes(conn)
	if err != nil {
		return
	}
	fmt.Println("Request:", str)
	//	server.Send(runes, c)
	conn.Close()
}
