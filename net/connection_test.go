package net

import (
	"bytes"
	"encoding/binary"
	"github.com/stojg/vivere/state"
	"testing"
)

func TestDeSerialization(t *testing.T) {

	buf := &bytes.Buffer{}

	// Add the user command
	cmd := UserCommand{}
	cmd.Actions |= 1 << state.DEAD
	cmd.Sequence = 12
	cmd.Msec = 50
	// Add the game tick
	binary.Write(buf, binary.LittleEndian, uint32(5))
	// Add the sequence number
	binary.Write(buf, binary.LittleEndian, cmd.Sequence)
	// Add the msec the command was run for
	binary.Write(buf, binary.LittleEndian, cmd.Msec)
	// Send the actions
	binary.Write(buf, binary.LittleEndian, cmd.Actions)

	cc := &ClientConn{}
	nextCmd, _ := cc.ReadMessage(buf)

	if nextCmd != cmd {
		t.Errorf("cmd sent and recieved mismatch %v != %v", nextCmd)
	}
}
