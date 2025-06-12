package server

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	Conn            net.Conn
	Address         string
	MsgChan         chan []byte
	BytesUploaded   int64
	BytesDownloaded int64
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

		c.BytesUploaded += int64(len(message))

		if c.TotalBytes() >= 100 {
			c.Conn.Write([]byte("Byte limit reached. Disconnecting.\n"))
			break
		}

		fmt.Printf("Received from %s: %s", c.Address, message)
	}
}

func (c *Client) receiveMessages() {
	for msg := range c.MsgChan {
		c.Conn.Write(msg)
		c.BytesDownloaded += int64(len(msg))

		if c.TotalBytes() >= 100 {
			c.Conn.Write([]byte("Byte limit reached. Disconnecting.\n"))
			return
		}
	}
}

func (c *Client) TotalBytes() int64 {
	return c.BytesUploaded + c.BytesDownloaded
}

func removeClient(addr string) {
	clientsMu.Lock()
	delete(clients, addr)
	clientsMu.Unlock()
}
