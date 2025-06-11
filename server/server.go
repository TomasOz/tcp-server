package server

import (
	"fmt"
	"net"
	"sync"
)

var (
	clients     = make(map[string]*Client)
	clientsMu   = sync.Mutex{}
	broadcastCh = make(chan Message)
)

type Message struct {
	Sender string
	Data   string
}

func Start(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer ln.Close()

	go broadcastLoop()

	fmt.Println("Server is listening on", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		client := NewClient(conn)

		fmt.Printf("New client connected: %s\n", conn.RemoteAddr())

		clientsMu.Lock()
		clients[client.Address] = client
		clientsMu.Unlock()

		go client.Handle(broadcastCh)
	}
}

func broadcastLoop() {
	for msg := range broadcastCh {
		clientsMu.Lock()
		for addr, client := range clients {
			if addr != msg.Sender {
				client.MsgChan <- []byte(msg.Data)
			}
		}
		clientsMu.Unlock()
	}
}

func removeClient(addr string) {
	clientsMu.Lock()
	delete(clients, addr)
	clientsMu.Unlock()
}
