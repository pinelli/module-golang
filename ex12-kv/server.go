package main

import (
	"fmt"
	"net"
	"os"
)

const cmdLen = 1
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

func readBytes(conn net.Conn) (string, string, string, error) {
	cmd := make([]byte, cmdLen, cmdLen)
	key := make([]byte, keyLen, keyLen)
	val := make([]byte, valLen, valLen)

	for i := 0; i < cmdLen; i++ {
		_, e := conn.Read(cmd[i : i+1]) //_, err
		if e != nil {
			return "", "", "", e
		}
	}
	for i := 0; i < keyLen; i++ {
		_, e := conn.Read(key[i : i+1]) //_, err
		if e != nil {
			return "", "", "", e
		}
	}
	for i := 0; i < valLen; i++ {
		_, e := conn.Read(val[i : i+1]) //_, err
		if e != nil {
			return "", "", "", e
		}
	}

	cmdStr := string(cmd)
	keyStr := string(key)
	valStr := string(val)

	return cmdStr, keyStr, valStr, nil
}

func (server *Server) handleConnection(conn net.Conn) {
	cmd, key, val, err := readBytes(conn)
	if err != nil {
		return
	}
	fmt.Println("Request->>", cmd, ":", key, ":", val)
	//	server.Send(runes, c)
	conn.Close()
}
