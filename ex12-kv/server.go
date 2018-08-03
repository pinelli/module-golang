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

func (server *Server) Send(cmdStr, valStr string, conn net.Conn) {
	cmd := []byte(cmdStr)
	val := []byte(valStr)

	for i := 0; i < cmdLen; i++ {
		_, e := conn.Write(cmd[i : i+1])
		if e != nil {
			return
		}
	}
	for i := 0; i < valLen; i++ {
		_, e := conn.Write(val[i : i+1])
		if e != nil {
			return
		}
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

func readBytes(conn net.Conn) (string, string, string, error) {
	cmd := make([]byte, cmdLen, cmdLen)
	key := make([]byte, keyLen, keyLen)
	val := make([]byte, valLen, valLen)

	for i := 0; i < cmdLen; i++ {
		_, e := conn.Read(cmd[i : i+1])
		if e != nil {
			return "", "", "", e
		}
	}
	for i := 0; i < keyLen; i++ {
		_, e := conn.Read(key[i : i+1])
		if e != nil {
			return "", "", "", e
		}
	}
	for i := 0; i < valLen; i++ {
		_, e := conn.Read(val[i : i+1])
		if e != nil {
			return "", "", "", e
		}
	}

	cmdStr := string(cmd)
	keyStr := string(key)
	valStr := string(val)

	return cmdStr, keyStr, valStr, nil
}

func
func (server *Server) handleConnection(conn net.Conn) {
	cmd, key, val, err := readBytes(conn)
	if err != nil {
		return
	}

	fmt.Println("Request->>", cmd, ":", key, ":", val)
	resVal, err := execute(cmd, key, val)
	server.Send(cmd, val, conn)

	conn.Close()
}
