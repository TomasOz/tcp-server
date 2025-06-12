package server

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	Conn            net.Conn
	Address         string
	MsgChan         chan []byte
	BytesUploaded   int64
	BytesDownloaded int64
	mu              sync.Mutex
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn:    conn,
		Address: conn.RemoteAddr().String(),
		MsgChan: make(chan []byte, 10),
	}
}

func (c *Client) Handle(broadcastCh chan Message) {
	defer func() {
		fmt.Printf("Client %s fully disconnected.\n", c.Address)
		removeClient(c.Address)
		c.Conn.Close()
	}()

	go c.receiveMessages()

	c.broadcastMessages(broadcastCh)
}

func (c *Client) broadcastMessages(broadcastCh chan Message) {
	reader := bufio.NewReader(c.Conn)

	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Client %s disconnected.\n", c.Address)
			break
		}

		broadcastCh <- Message{
			Sender: c.Address,
			Data:   string(message),
		}

		if c.addBytesAndCheck(int64(len(message)), 0) {
			c.Conn.Write([]byte("Byte limit reached. Disconnecting.\n"))
			break
		}

		fmt.Printf("Received from %s: %s", c.Address, message)
	}
}

func (c *Client) receiveMessages() {
	for msg := range c.MsgChan {
		c.Conn.Write(msg)

		if c.addBytesAndCheck(0, int64(len(msg))) {
			c.Conn.Write([]byte("Byte limit reached. Disconnecting.\n"))
			return
		}
	}
}

func (c *Client) addBytesAndCheck(upload, download int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.BytesUploaded += upload
	c.BytesDownloaded += download

	return c.BytesUploaded+c.BytesDownloaded >= 100
}

func removeClient(addr string) {
	clientsMu.Lock()
	delete(clients, addr)
	clientsMu.Unlock()
}
