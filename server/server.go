package server

import (
	"fmt"
	"net"
)

func Start(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Println("Server is listening on", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		client := NewClient(conn)

		fmt.Printf("New client connected: %s\n", conn.RemoteAddr())

		go client.Handle()
	}
}
