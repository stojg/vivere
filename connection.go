package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/binary"
	"fmt"
	"io"
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

var newConn = make(chan *ClientConn)

func (cc *ClientConn) Close() {
	cc.ws.Close()
	cc.open = false
}

// ReadMessage picks the current message from the inbuffer and
func (cc *ClientConn) ReadMessage(reader io.Reader) (cmd UserCommand, err error) {

	pkt := cc.inBuf[0:]

	n, err := reader.Read(pkt)
	if err != nil {
		return cmd, fmt.Errorf("ws.Read() - error during read: '%s'", err)
	}

	pkt = pkt[0:n]
	buffer := bytes.NewBuffer(pkt)

	err = binary.Read(buffer, binary.LittleEndian, &cc.tick)
	if err != nil {
		return cmd, fmt.Errorf("binary.Read() - Couldn't read tick: '%s'", err)
	}

	err = binary.Read(buffer, binary.LittleEndian, &cmd.Sequence)
	if err != nil {
		return cmd, fmt.Errorf("binary.Read() - Couldn't read sequence: '%s'", err)
	}

	err = binary.Read(buffer, binary.LittleEndian, &cmd.Actions)
	if err != nil {
		return cmd, fmt.Errorf("binary.Read() - Couldn't read command: '%s'", err)
	}
	return cmd, nil
}

func wsHandler(ws *websocket.Conn) {
	clientConn := &ClientConn{}
	clientConn.ws = ws
	clientConn.open = true
	clientConn.cmdBuf = make(chan UserCommand, 5)

	// Push the new connection to the newConn channel
	newConn <- clientConn

	// Read messages from the client
	for {
		cmd, err := clientConn.ReadMessage(ws)
		if err != nil {
			log.Printf("[!] connection: %s", err)
			break
		}
		clientConn.cmdBuf <- cmd
	}
}
