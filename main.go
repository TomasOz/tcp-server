package main

import (
	"fmt"
	"tcpserver/server"
)

func main() {
	err := server.Start(":9000")
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	fmt.Println("Server is listening on port 9000...")
}
