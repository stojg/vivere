package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/binary"
	"log"
)

// ClientConn is the current connection and the current command
type ClientConn struct {
	ws         *websocket.Conn
	inBuf      [1500]byte
	currentCmd UserCommand
	cmdBuf     chan UserCommand
}

var newConn = make(chan *ClientConn)

var clients = make(map[PlayerId]*ClientConn)

func wsHandler(ws *websocket.Conn) {
	clientConn := &ClientConn{}
	clientConn.ws = ws
	// keep 5 User commands in the buffer
	clientConn.cmdBuf = make(chan UserCommand, 5)

	// Create a new UserCommand
	cmd := UserCommand{}

	log.Println("wsHandler: new client connection")

	// Push the new connection to the newConn channel
	newConn <- clientConn

	// Infinite loop that reads UserCommands from the client
	for {
		pkt := clientConn.inBuf[0:]
		n, err := ws.Read(pkt)
		// Oh noes, client probably disconnected during read
		if err != nil {
			log.Printf("wsHandler: Error during read '%s'\n", err)
			break
		}
		// Reassign all packets into the pkt buffer
		pkt = pkt[0:n]
		buf := bytes.NewBuffer(pkt)
		err = binary.Read(buf, binary.LittleEndian, &cmd)
		// Oh noes, couldn't read the user command
		if err != nil {
			log.Printf("wsHandler: error '%s'\n", err)
			break
		}
		// Push the cmd to the clientCommand channel
		clientConn.cmdBuf <- cmd
	}
}
