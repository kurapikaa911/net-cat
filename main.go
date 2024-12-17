package main

import (
	"fmt"
	"net"
	"os"

	"netcat/utils"
)

var (
	ADRESS = "0.0.0.0:"
	PORT   = "8989"
)

func main() {
	if len(os.Args) != 1 {
		if len(os.Args) == 2 {
			PORT = os.Args[1]
		} else {
			fmt.Fprintln(os.Stderr, "[USAGE]: ./TCPChat $port")
			return
		}
	}

	ln, err := net.Listen("tcp", ADRESS+PORT)
	if err != nil {
		PORT = "8989"
		fmt.Println("invalid port...using default port", PORT)
		ln, err = net.Listen("tcp", ADRESS+PORT)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}

	fmt.Println("Listening on the port :" + PORT)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go utils.HandleConnection(&conn, PORT)
	}
}
