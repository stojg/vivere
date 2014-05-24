package main

import (
	"encoding/binary"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestGetNewAndIdGeneration(c *C) {
	ol := NewEntityList()
	ent := ol.NewEntity()
	c.Assert(ent.ID(), Equals, uint16(1))
	ent = ol.NewEntity()
	c.Assert(ent.ID(), Equals, uint16(2))
	ent = ol.NewEntity()
	c.Assert(ent.ID(), Equals, uint16(3))
}

func (s *TestSuite) TestSetPosition(c *C) {
	gol := NewEntityList()
	g := gol.NewEntity()
	g.Position.Set(100, 200, 0)
	c.Assert(g.Position[0], Equals, float64(100))
	c.Assert(g.Position[1], Equals, float64(200))
}

func (s *TestSuite) TestSetScale(c *C) {
	gol := NewEntityList()
	g := gol.NewEntity()
	g.scale.Set(1, 2, 0)
	c.Assert(g.scale[0], Equals, float64(1))
	c.Assert(g.scale[1], Equals, float64(2))
}

func (s *TestSuite) TestList(c *C) {
	gol := NewEntityList()
	g := gol.NewEntity()
	c.Assert(gol.Get(g.ID()), Equals, g)
	c.Assert(gol.Length(), Equals, 1)
	gol.Remove(g.ID())
	c.Assert(gol.Length(), Equals, 0)
}

func (s *TestSuite) TestSimpleUpdate(c *C) {
	gol := NewEntityList()
	g := gol.NewEntity()
	c.Assert(g.Changed(), Equals, true)
	//	g.orientation = 2
	//	c.Assert(g.Changed(), Equals, true)
}

func (s *TestSuite) TestSerialize(c *C) {
	gol := NewEntityList()
	g := gol.NewEntity()
	g.Position.Set(10, 20, 0)

	var literal byte
	var expected float32
	buf := g.Serialize()

	// id
	binary.Read(buf, binary.LittleEndian, &literal)
	c.Assert(literal, Equals, byte(INST_ENTITY_ID))
	binary.Read(buf, binary.LittleEndian, &expected)
	c.Assert(float32(1), Equals, expected)

	// get position
	binary.Read(buf, binary.LittleEndian, &literal)
	c.Assert(literal, Equals, byte(INST_SET_POSITION))
	binary.Read(buf, binary.LittleEndian, &expected)
	c.Assert(float32(10), Equals, expected)
	// We expect the server to run with y + 1 is up and the client is rendering it in screen coordinates
	binary.Read(buf, binary.LittleEndian, &expected)
	c.Assert(float32(-20), Equals, expected)
}
