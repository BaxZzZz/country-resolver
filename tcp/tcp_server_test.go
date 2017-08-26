package tcp

import (
	"net"
	"testing"
	"time"
	"errors"
)

const waitTimeout = 100 * time.Millisecond
const testAddress = "localhost:6666"

func startClient(address string) (net.Conn, error) {
	var reconnectCount uint
	for {
		if reconnectCount > 5 {
			return nil, errors.New("Failed to connect to server")
		}
		clientConn, err := net.Dial("tcp", address)
		if err == nil {
			return clientConn, nil
		}
		time.Sleep(waitTimeout)
		reconnectCount++
	}
}

func TestAcceptingNewClient(t *testing.T) {
	var isAccepted bool

	server := NewTcpServer()
	server.OnNewClient(func(client *Client) {
		isAccepted = true
	})

	err := server.Start(testAddress)
	if err != nil {
		t.Fatalf("Failed to start server, %v", err)
	}

	client, err := startClient(testAddress)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(waitTimeout)
	server.Stop()
	client.Close()

	if !isAccepted {
		t.Fatal("Server not accepted client")
	}
}

func TestDisconnectingClient(t *testing.T) {
	var isDisconnected bool

	server := NewTcpServer()
	server.OnClientDisconnected(func(*Client, error) {
		isDisconnected = true
	})

	err := server.Start("localhost:6666")
	if err != nil {
		t.Fatalf("Failed to start server, %v", err)
	}

	client, err := startClient(testAddress)
	if err != nil {
		t.Fatal(err)
	}
	client.Close()

	time.Sleep(waitTimeout)
	server.Stop()

	if !isDisconnected {
		t.Fatal("Client not disconnected")
	}
}

func TestMessageReceivedFromClient(t *testing.T) {
	var messageReceived string
	var isReceived bool

	server := NewTcpServer()
	server.OnClientMessageReceived(func(client *Client, message string) {
		isReceived = true
		messageReceived = message
	})

	server.Start("localhost:6666")
	time.Sleep(waitTimeout)

	client, err := startClient(testAddress)
	if err != nil {
		t.Fatal(err)
	}
	client.Write([]byte("Ping message\n"))

	time.Sleep(waitTimeout)
	server.Stop()
	client.Close()

	if !isReceived {
		t.Fatal("Message not received")
	}

	if messageReceived != "Ping message\n" {
		t.Fatal("Massage not equal")
	}
}
