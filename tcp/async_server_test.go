package tcp

import (
	"errors"
	"net"
	"testing"
	"time"
)

const waitTimeout = 100 * time.Millisecond

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

	server, err := NewServer("localhost:6667")
	if err != nil {
		t.Fatal(err)
	}

	server.OnNewClient(func(client *Client) {
		isAccepted = true
	})

	go server.Listen()

	client, err := startClient("localhost:6667")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(waitTimeout)
	client.Close()
	server.Shutdown()

	if !isAccepted {
		t.Fatal("Server not accepted client")
	}
}

func TestDisconnectingClient(t *testing.T) {
	var isDisconnected bool

	server, err := NewServer("localhost:6668")
	if err != nil {
		t.Fatal(err)
	}

	server.OnClientDisconnected(func(*Client, error) {
		isDisconnected = true
	})

	go server.Listen()

	client, err := startClient("localhost:6668")
	if err != nil {
		t.Fatal(err)
	}

	client.Close()

	time.Sleep(waitTimeout)
	server.Shutdown()

	if !isDisconnected {
		t.Fatal("Client not disconnected")
	}
}

func TestMessageReceivedFromClient(t *testing.T) {
	var messageReceived string
	var isReceived bool

	server, err := NewServer("localhost:6669")
	if err != nil {
		t.Fatal(err)
	}

	server.OnClientMessageReceived(func(client *Client, message string) {
		isReceived = true
		messageReceived = message
	})

	go server.Listen()

	client, err := startClient("localhost:6669")
	if err != nil {
		t.Fatal(err)
	}
	client.Write([]byte("Ping message\n"))

	time.Sleep(waitTimeout)
	client.Close()
	server.Shutdown()

	if !isReceived {
		t.Fatal("Message not received")
	}

	if messageReceived != "Ping message\n" {
		t.Fatal("Massage not equal")
	}
}
