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

	go func() {
		for msg := range c.MsgChan {
			c.Conn.Write(msg)
			c.BytesDownloaded += int64(len(msg))

			fmt.Printf("Client %s: already used %d bytes of data\n", c.Address, c.TotalBytes())

			if c.TotalBytes() >= 100 {
				c.Conn.Write([]byte("Byte limit reached. Disconnecting.\n"))
				return
			}
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

		c.BytesUploaded += int64(len(message))

		fmt.Printf("Client %s: already used %d bytes of data\n", c.Address, c.TotalBytes())

		if c.TotalBytes() >= 100 {
			c.Conn.Write([]byte("Byte limit reached. Disconnecting.\n"))
			break
		}

		fmt.Printf("Received from %s: %s", c.Address, message)
	}
}

func (c *Client) TotalBytes() int64 {
	return c.BytesUploaded + c.BytesDownloaded
}
