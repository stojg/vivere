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
	tick       uint32
	cmdBuf     chan UserCommand
	open       bool
}

func (cc *ClientConn) Close() {
	cc.ws.Close()
	cc.open = false
}

var newConn = make(chan *ClientConn)

func wsHandler(ws *websocket.Conn) {
	clientConn := &ClientConn{}
	clientConn.ws = ws
	clientConn.open = true
	clientConn.cmdBuf = make(chan UserCommand, 5)

	// Create a new UserCommand
	cmd := UserCommand{}

	// Push the new connection to the newConn channel
	newConn <- clientConn

	// Read messages from the client
	for {
		pkt := clientConn.inBuf[0:]

		n, err := ws.Read(pkt)
		if err != nil {
			log.Printf("[-] ws.Read() - error during read '%s'\n", err)
			break
		}

		pkt = pkt[0:n]
		buf := bytes.NewBuffer(pkt)

		err = binary.Read(buf, binary.LittleEndian, &clientConn.tick)
		if err != nil {
			log.Printf("[-] binary.Read() - Couldn't read tick '%s'\n", err)
			break
		}

		err = binary.Read(buf, binary.LittleEndian, &cmd)
		if err != nil {
			log.Printf("[-] binary.Read() - Couldn't read command '%s'\n", err)
			break
		}

		clientConn.cmdBuf <- cmd
	}
}
