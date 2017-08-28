package tcp

import (
	"bufio"
	"log"
	"net"
)

// Client holds info about connection
type Client struct {
	connection net.Conn
	server     *AsyncServer
}

// Async TCP server
type AsyncServer struct {
	listener                net.Listener
	address                 string
	doStop                  chan bool
	isStopped               chan bool
	newClientHandler        func(client *Client)
	clientDisconnectHandler func(client *Client, err error)
	clientMessageHandler    func(client *Client, message string)
}

func (client *Client) readMessage() {
	reader := bufio.NewReader(client.connection)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			client.connection.Close()
			client.server.clientDisconnectHandler(client, err)
			return
		}
		client.server.clientMessageHandler(client, message)
	}
}

// Send text message to client
func (client *Client) SendMessage(message string) error {
	_, err := client.connection.Write([]byte(message))
	return err
}

// Close client connection
func (client *Client) Close() error {
	return client.connection.Close()
}

// Get client IP address
func (client *Client) GetRemoteIpAddress() (string, error) {
	ip, _, err := net.SplitHostPort(client.connection.RemoteAddr().String())
	if err != nil {
		return "", err
	}

	return ip, nil
}

// Read client data from channel
func (server *AsyncServer) Listen() {
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			select {
			case <-server.doStop:
				server.isStopped <- true
				return
			default:
				log.Printf("Accept failed: %v", err)
			}
		}

		client := &Client{
			connection: connection,
			server:     server,
		}

		server.newClientHandler(client)
		go client.readMessage()
	}
}

// Shutdown working server
func (server *AsyncServer) Shutdown() error {
	server.doStop <- true
	return server.listener.Close()
	<-server.isStopped
	return nil
}

// Set handler for handling new connection with client
func (server *AsyncServer) OnNewClient(callbackFunc func(*Client)) {
	server.newClientHandler = callbackFunc
}

// Set handler for handling disconnection with client
func (server *AsyncServer) OnClientDisconnected(callbackFunc func(*Client, error)) {
	server.clientDisconnectHandler = callbackFunc
}

// Set handler for handling client message
func (server *AsyncServer) OnClientMessageReceived(callbackFunc func(*Client, string)) {
	server.clientMessageHandler = callbackFunc
}

// Creates new TCP server instance
func NewServer(address string) (*AsyncServer, error) {
	log.Println("Start TCP server on " + address)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	server := &AsyncServer{
		listener:  listener,
		address:   address,
		doStop:    make(chan bool, 1),
		isStopped: make(chan bool),
	}

	server.OnNewClient(func(*Client) {})
	server.OnClientDisconnected(func(*Client, error) {})
	server.OnClientMessageReceived(func(*Client, string) {})

	return server, nil
}
