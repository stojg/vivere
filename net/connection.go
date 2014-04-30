package net

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
	cmdBuf     chan UserCommand
	inBuf      [1500]byte
	currentCmd UserCommand
	tick       uint32
	open       bool
}

func (cc *ClientConn) Ws() *websocket.Conn {
	return cc.ws
}

// UserCommand represent a recieved command (Action) from the user
type UserCommand struct {
	Actions  uint32
	Sequence uint32
	Msec     uint32
}

// Close the client connection
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

	err = binary.Read(buffer, binary.LittleEndian, &cmd.Msec)
	if err != nil {
		return cmd, fmt.Errorf("binary.Read() - Couldn't read msec: '%s'", err)
	}

	err = binary.Read(buffer, binary.LittleEndian, &cmd.Actions)
	if err != nil {
		return cmd, fmt.Errorf("binary.Read() - Couldn't read command: '%s'", err)
	}
	return cmd, nil
}


type ConnectionHandler struct {
	NewConnections chan *ClientConn
}

func NewConnectionHandler() *ConnectionHandler {
	ch := &ConnectionHandler{}
	ch.NewConnections = make(chan *ClientConn)
	return ch
}

func (ch *ConnectionHandler) NewConn() chan *ClientConn {
	return ch.NewConnections
}

func (ch *ConnectionHandler) WsHandler(ws *websocket.Conn) {
	clientConn := &ClientConn{}
	clientConn.ws = ws
	clientConn.open = true
	clientConn.cmdBuf = make(chan UserCommand, 5)

	// Push the new connection to the newConn channel
	ch.NewConnections <- clientConn

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

func Send(player *Player, buf *bytes.Buffer) (error) {
	if buf.Len() == 0 {
		return nil
	}
	err := websocket.Message.Send(player.conn.ws, buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func GetUpdates(players []*Player) {
	for _, player := range players {
		for {
			select {
			case cmd := <-player.conn.cmdBuf:
				player.conn.currentCmd = cmd
			default:
				goto done
			}
		}
	done:
	}
}

func Connect(conn *ClientConn, id uint16) *Player{
	p := &Player{}
	p.id = id
	p.conn = conn
	return p

}

func Disconnect(p *Player) {
	p.conn.Close()
}
