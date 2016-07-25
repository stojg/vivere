// Package client provides structs and function that deals with websocket connections.
// The messages are sent in a binary format and should always start with a 'header' containing
// - timestamp float64
// - messageType MessageType
package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"time"
)

type MessageType uint8

// ClientHandler represent a group of connection and wires them all up
// to individual Clients
type ClientHandler struct {
	newClientChan chan *Client
}

// NewClientHandler returns a new instance of a ClientHandler
func NewClientManager() *ClientHandler {
	ch := &ClientHandler{}
	ch.newClientChan = make(chan *Client)
	return ch
}

// Websocket is a websocket handler func that should be used like this
//     http.Handle("/ws/", websocket.Handler(ch.Websocket))
// This will create a new Client for each new connection and setup
// the necessary go channel to communicate a new client connection
// to other parts of the engine
func (ch *ClientHandler) Websocket(ws *websocket.Conn) {
	client := &Client{}
	client.ws = ws
	client.open = true
	client.cmdBuf = make(chan ClientCommand, 5)
	client.pingStartTime = 0
	ch.NewClients() <- client
	// Buffer message from the client
	for {
		cmd, err := client.Read(client.ws)
		if err == io.EOF {
			log.Printf("[!] client disconnected", err)
			break
		}
		if err != nil {
			log.Printf("[!] connection: %s", err)
			break
		}

		client.cmdBuf <- cmd
	}
}

// NewClients returns channel that will push a *Client through the channel
func (ch *ClientHandler) NewClients() chan *Client {
	return ch.newClientChan
}

// ClientCommand what the client sends to the server, it represents actions
// that the user issued, for example clicking the up arrow key
type ClientCommand struct {
	Data  uint32
	Type MessageType
	Sequence uint32
	Duration float64
}

// Client represents a open websocket connection, ie a user.
type Client struct {
	ws            *websocket.Conn
	cmdBuf        chan ClientCommand
	inBuf         [1500]byte
	currentCmd    ClientCommand
	open          bool
	pingStartTime float64
	ping          float64
	serverTime    float64
}

func (c *Client) Input() chan ClientCommand {
	return c.cmdBuf
}

// Write provides a io.reader interface for writing a message to the client
func (c *Client) Write(p []byte) (n int, err error) {
	n = len(p)
	err = nil
	if n == 0 {
		err = fmt.Errorf("Write() - zero length message")
		return
	}
	err = websocket.Message.Send(c.ws, p)
	if err != nil {
		err = fmt.Errorf("Write() - error during send %v", err)
		return
	}
	return
}

// Read reads a buffer in binary format, extracts the server time and message type
// and passes the rest of the message on to a message handler
func (c *Client) Read(reader io.Reader) (cmd ClientCommand, err error) {
	var msgType MessageType

	pkt := c.inBuf[0:]

	n, err := reader.Read(pkt)
	if err != nil {
		return cmd, fmt.Errorf("ws.Read() - error during read: '%s'", err)
	}

	pkt = pkt[0:n]
	buffer := bytes.NewBuffer(pkt)

	err = binary.Read(buffer, binary.LittleEndian, &c.serverTime)
	if err != nil {
		return cmd, fmt.Errorf("binary.Read() - Couldn't read c.serverTime #: '%s'", err)
	}

	err = binary.Read(buffer, binary.LittleEndian, &msgType)
	if err != nil {
		return cmd, fmt.Errorf("binary.Read() - Couldn't read msgType #: '%s'", err)
	}

	return c.input(buffer, msgType)
}

// NewMessage returns a buffer writer ready for binary writing including
// the message type and the current server timestamp
func (client *Client) NewMessage(msgType MessageType) *bytes.Buffer {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, client.timestamp())
	binary.Write(buf, binary.LittleEndian, msgType)
	return buf
}

func (client *Client) Input() chan ClientCommand {
	return client.cmdBuf
}

// Returns the current unix nano time as float64 since
// there is some problem reading int64 in javascript
func (client *Client) timestamp() float64 {
	return float64(time.Now().UnixNano())
}

// Updatea
func (client *Client) Update(buf *bytes.Buffer, msgType MessageType) {
	message := client.NewMessage(msgType)
	message.Write(buf.Bytes())
	client.Write(message.Bytes())
}

// Ping sends a ping request to the client
//func (client *Client) Ping() {
//	if client.pingStartTime != 0 {
//		return
//	}
//	client.pingStartTime = client.timestamp()
//	buf := client.NewMessage(MSG_PING)
//	client.Write(buf.Bytes())
//}

// pong updates the client connection with the latest ping response from
// the client.
func (client *Client) pingResponse(reader io.Reader) {
	client.ping = (client.timestamp() - client.pingStartTime) / 1e6
	client.pingStartTime = 0
}

// ReadMessage picks the current message from the inbuffer and
func (c *Client) input(reader io.Reader, msgType MessageType) (cmd ClientCommand, err error) {
	buffer := reader
	cmd.Type = msgType
	err = binary.Read(buffer, binary.LittleEndian, &cmd.Sequence)
	if err != nil {
		return cmd, fmt.Errorf("binary.Read() - Couldn't read uint32 Sequence #: '%s'", err)
	}
	err = binary.Read(buffer, binary.LittleEndian, &cmd.Duration)
	if err != nil {
		return cmd, fmt.Errorf("binary.Read() - Couldn't read float64 Duration: '%s'", err)
	}
	err = binary.Read(buffer, binary.LittleEndian, &cmd.Data)
	if err != nil {
		return cmd, fmt.Errorf("binary.Read() - Couldn't read Data: '%s'", err)
	}
	return cmd, nil
}
