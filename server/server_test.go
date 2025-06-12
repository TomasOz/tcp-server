package server

import (
	"net"
	"strings"
	"testing"
	"time"
)

func TestMessageBroadcast(t *testing.T) {
	ln, _ := net.Listen("tcp", ":0")
	defer ln.Close()

	clients = make(map[string]*Client)
	broadcastCh = make(chan Message, 10)

	go func() {
		go broadcastLoop()
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			client := NewClient(conn)
			clientsMu.Lock()
			clients[client.Address] = client
			clientsMu.Unlock()
			go client.Handle(broadcastCh)
		}
	}()

	connection1, _ := net.Dial("tcp", ln.Addr().String())
	defer connection1.Close()
	connection2, _ := net.Dial("tcp", ln.Addr().String())
	defer connection2.Close()

	time.Sleep(100 * time.Millisecond)

	connection1.Write([]byte("test\n"))

	buffer := make([]byte, 100)
	connection2.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := connection2.Read(buffer)

	if err != nil || string(buffer[:n]) != "test\n" {
		t.Error("Message broadcast failed")
	}
}

func TestByteLimitEnforcement(t *testing.T) {
	ln, _ := net.Listen("tcp", ":0")
	defer ln.Close()

	clients = make(map[string]*Client)
	broadcastCh = make(chan Message, 10)

	go func() {
		go broadcastLoop()
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			client := NewClient(conn)
			clientsMu.Lock()
			clients[client.Address] = client
			clientsMu.Unlock()
			go client.Handle(broadcastCh)
		}
	}()

	conn, _ := net.Dial("tcp", ln.Addr().String())
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 8; i++ {
		conn.Write([]byte("0123456789012\n"))
		time.Sleep(20 * time.Millisecond)
	}

	buffer := make([]byte, 100)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, _ := conn.Read(buffer)

	response := string(buffer[:n])
	if !strings.Contains(response, "limit") {
		t.Error("Expected byte limit notification")
	}
}
