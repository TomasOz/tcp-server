package server

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	Conn    net.Conn
	Address string
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn:    conn,
		Address: conn.RemoteAddr().String(),
	}
}

func (c *Client) Handle() {
	defer c.Conn.Close()
	reader := bufio.NewReader(c.Conn)

	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Client %s disconnected.\n", c.Address)
			break
		}

		fmt.Printf("Received from %s: %s", c.Address, message)
	}
}
