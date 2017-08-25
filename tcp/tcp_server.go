package tcp

import (
	"bufio"
	"log"
	"net"
)

type Client struct {
	connection net.Conn
	server     *tcpServer
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

func (client *Client) SendMessage(message string) error {
	_, err := client.connection.Write([]byte(message))
	return err
}

func (client *Client) Close() error {
	return client.connection.Close()
}

func (client *Client) GetRemoteIpAddress() (string, error) {
	ip, _, err := net.SplitHostPort(client.connection.RemoteAddr().String())
	if err != nil {
		return "", err
	}

	return ip, nil
}

type tcpServer struct {
	listener                net.Listener
	address                 string
	done                    chan bool
	newClientHandler        func(client *Client)
	clientDisconnectHandler func(client *Client, err error)
	clientMessageHandler    func(client *Client, message string)
}

func (server *tcpServer) accept() {
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			select {
			case <-server.done:
			default:
				log.Printf("Accept failed: %v", err)
			}
			return
		}

		client := &Client{
			connection: connection,
			server:     server,
		}

		server.newClientHandler(client)
		go client.readMessage()
	}
}

func (server *tcpServer) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server.listener = listener
	server.address = address

	go server.accept()

	return err
}

func (server *tcpServer) Stop() error {
	server.done <- true
	return server.listener.Close()
}

func (server *tcpServer) OnNewClient(callbackFunc func(*Client)) {
	server.newClientHandler = callbackFunc
}

func (server *tcpServer) OnClientDisconnected(callbackFunc func(*Client, error)) {
	server.clientDisconnectHandler = callbackFunc
}

func (server *tcpServer) OnClientMessageReceived(callbackFunc func(*Client, string)) {
	server.clientMessageHandler = callbackFunc
}

func NewTcpServer() *tcpServer {
	server := &tcpServer{
		done: make(chan bool, 1),
	}

	server.OnNewClient(func(*Client) {})
	server.OnClientDisconnected(func(*Client, error) {})
	server.OnClientMessageReceived(func(*Client, string) {})

	return server
}
