package tcp

import (
	"net"
	"testing"
	"time"
)

func TestAcceptingNewClient(t *testing.T) {
	var isAccepted bool

	server := NewTcpServer()
	server.OnNewClient(func(client *Client) {
		isAccepted = true
	})

	server.Start("localhost:6666")

	time.Sleep(10 * time.Millisecond)

	clientConn, err := net.Dial("tcp", "localhost:6666")
	if err != nil {
		t.Fatal("Failed tp connect to server")
	}

	clientConn.Close()

	time.Sleep(10 * time.Millisecond)
	server.Stop()

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

	server.Start("localhost:6666")

	time.Sleep(10 * time.Millisecond)

	clientConn, err := net.Dial("tcp", "localhost:6666")
	if err != nil {
		t.Fatal("Failed tcp connect to server")
	}

	clientConn.Close()

	time.Sleep(10 * time.Millisecond)
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

	time.Sleep(10 * time.Millisecond)

	clientConn, err := net.Dial("tcp", "localhost:6666")
	if err != nil {
		t.Fatal("Failed tp connect to server")
	}

	clientConn.Write([]byte("Ping message\n"))
	clientConn.Close()

	time.Sleep(10 * time.Millisecond)
	server.Stop()

	if !isReceived {
		t.Fatal("Message not received")
	}

	if messageReceived != "Ping message\n" {
		t.Fatal("Massage not equal")
	}
}
