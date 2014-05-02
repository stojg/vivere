package net

import (
	"bytes"
	"encoding/binary"
	"github.com/stojg/vivere/state"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestDeserialization(c *C) {

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

	c.Assert(nextCmd, Equals, cmd)
}

func (s *TestSuite) TestDeserializationBadPkt(c *C) {

	buf := &bytes.Buffer{}

	type BadUserCommand struct {
		Actions  uint32
		Sequence int
		Msec     uint32
	}

	// Add the user command
	cmd := BadUserCommand{}
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

	c.Assert(nextCmd.Actions, Equals, cmd.Actions)
	c.Assert(nextCmd.Sequence, Equals, cmd.Sequence)
	c.Assert(nextCmd.Msec, Equals, cmd.Msec)
}
