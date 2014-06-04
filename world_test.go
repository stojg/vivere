package main

import (
	"encoding/binary"
	. "gopkg.in/check.v1"
)

type TestWorldSuite struct{}

var _ = Suite(&TestWorldSuite{})

func (s *TestWorldSuite) TestSerialize(c *C) {
	w := NewWorld(false)
	w.Tick = 5
	buf := w.Serialize(true)
	var expected float32
	binary.Read(buf, binary.LittleEndian, &expected)
	c.Assert(expected, Equals, float32(5))
}
