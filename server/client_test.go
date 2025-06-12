package server

import (
	"net"
	"testing"
	"time"
)

func TestClientCreation(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	c := NewClient(server)

	if c.Conn != server {
		t.Error("Client connection not set correctly")
	}
	if c.BytesUploaded != 0 || c.BytesDownloaded != 0 {
		t.Error("Client byte counters should start at 0")
	}
}

func TestAddBytesAndCheck(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	c := NewClient(server)

	limitReached := c.addBytesAndCheck(50, 30)

	if limitReached {
		t.Errorf("Limit should not be reached yet, got: %v", limitReached)
	}
	if c.BytesUploaded != 50 {
		t.Errorf("Expected 50 uploaded bytes, got %d", c.BytesUploaded)
	}
	if c.BytesDownloaded != 30 {
		t.Errorf("Expected 30 downloaded bytes, got %d", c.BytesDownloaded)
	}

	limitReached = c.addBytesAndCheck(25, 0)

	if !limitReached {
		t.Errorf("Expected limit to be reached, got: %v", limitReached)
	}
	if c.BytesUploaded != 75 {
		t.Errorf("Expected 75 uploaded bytes, got %d", c.BytesUploaded)
	}
	if c.BytesDownloaded != 30 {
		t.Errorf("Expected total bytes 30, got %d", c.BytesDownloaded)
	}
}

func TestRemoveClient(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	c := NewClient(server)
	c.Address = "test-address"

	clientsMu.Lock()
	clients[c.Address] = c
	clientsMu.Unlock()

	removeClient(c.Address)

	clientsMu.Lock()
	if _, exists := clients[c.Address]; exists {
		t.Error("Client was not removed from global map")
	}
	clientsMu.Unlock()
}

func TestClientByteCountingOnReceive(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	client := NewClient(serverConn)

	done := make(chan struct{})
	go func() {
		client.receiveMessages()
		close(done)
	}()

	testMessage := []byte("Hello World")
	client.MsgChan <- testMessage
	close(client.MsgChan)

	buffer := make([]byte, 100)
	clientConn.SetReadDeadline(time.Now().Add(time.Second))

	bytesRead, err := clientConn.Read(buffer)
	if err != nil {
		t.Fatal("Failed to read message:", err)
	}

	<-done

	if bytesRead != len(testMessage) {
		t.Errorf("Expected %d bytes, got %d", len(testMessage), bytesRead)
	}
	if client.BytesDownloaded != int64(len(testMessage)) {
		t.Errorf("Expected %d downloaded bytes, got %d", len(testMessage), client.BytesDownloaded)
	}
	if client.BytesUploaded != 0 {
		t.Errorf("Expected 0 uploaded bytes, got %d", client.BytesUploaded)
	}
}

func TestClientByteCountingOnBroadcast(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	client := NewClient(serverConn)
	broadcastCh := make(chan Message, 10)

	go client.broadcastMessages(broadcastCh)

	testMessage := "test message\n"
	clientConn.Write([]byte(testMessage))
	clientConn.Close()

	msg := <-broadcastCh
	if msg.Data != testMessage {
		t.Errorf("Expected %q, got %q", testMessage, msg.Data)
	}
	if client.BytesUploaded != int64(len(testMessage)) {
		t.Errorf("Expected %d uploaded bytes, got %d", len(testMessage), client.BytesUploaded)
	}
	if client.BytesDownloaded != 0 {
		t.Errorf("Expected 0 uploaded bytes, got %d", client.BytesDownloaded)
	}

}
