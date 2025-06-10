package server

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	Conn    net.Conn
	Address string
	MsgChan chan []byte
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn:    conn,
		Address: conn.RemoteAddr().String(),
		MsgChan: make(chan []byte, 10),
	}
}

func (c *Client) Handle(broadcastCh chan Message) {
	defer c.Conn.Close()

	go func() {
		for msg := range c.MsgChan {
			c.Conn.Write(msg)
		}
	}()

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

		fmt.Printf("Received from %s: %s", c.Address, message)
	}
}
