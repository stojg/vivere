package main

import (
	"bytes"
	"encoding/binary"
	"testing"
	//"time"
)

func TestDeSerialization(t *testing.T) {

	// Add the tick
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, uint32(5))

	// Add the user command
	cmd := UserCommand{}
	cmd.Actions |= 1 << ACTION_DIE
	cmd.Sequence = 12
	// Add the command number
	binary.Write(buf, binary.LittleEndian, cmd.Sequence)
	binary.Write(buf, binary.LittleEndian, cmd.Actions)

	cc := &ClientConn{}
	nextCmd, _ := cc.ReadMessage(buf)

	if nextCmd != cmd {
		t.Errorf("cmd sent and recieved mismatch %v != %v", nextCmd )
	}
}
