package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 9000...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v", err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			fmt.Printf("New connection from %s\n", c.RemoteAddr())
		}(conn)
	}
}
