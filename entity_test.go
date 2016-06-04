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
	g.Position.Set(100, 200, 300)
	c.Assert(g.Position[0], Equals, float64(100))
	c.Assert(g.Position[1], Equals, float64(200))
	c.Assert(g.Position[2], Equals, float64(300))
}

func (s *TestSuite) TestSetScale(c *C) {
	gol := NewEntityList()
	g := gol.NewEntity()
	g.Scale.Set(1, 2, 0)
	c.Assert(g.Scale[0], Equals, float64(1))
	c.Assert(g.Scale[1], Equals, float64(2))
}

func (s *TestSuite) TestList(c *C) {
	gol := NewEntityList()
	g := gol.NewEntity()
	c.Assert(gol.Get(g.ID()), Equals, g)
	c.Assert(gol.Length(), Equals, 1)
	gol.Remove(g.ID())
	c.Assert(gol.Length(), Equals, 0)
}

func TestEntitySerialize(t *testing.T) {
	gol := NewEntityList()
	g := gol.NewEntity()
	g.Position.Set(10, 20, 10)

	var command byte
	var value float32
	buf := g.Serialize()

	// id
	binary.Read(buf, binary.LittleEndian, &command)
	if command != byte(INST_ENTITY_ID) {
		t.Errorf("INST_ENTITY_ID wasn't serialised")
	}
	binary.Read(buf, binary.LittleEndian, &value)
	if g.id != uint16(value) {
		t.Errorf("ID wasn't serialised properly")
	}

	// get position
	binary.Read(buf, binary.LittleEndian, &command)
	if command != byte(INST_SET_POSITION) {
		t.Errorf("INST_SET_POSITION wasn't serialised")
	}
	binary.Read(buf, binary.LittleEndian, &value)
	if float32(g.Position[0]) != value {
		t.Errorf("Position.X wasn't serialised properly")
	}
	binary.Read(buf, binary.LittleEndian, &value)
	if float32(g.Position[1]) != value {
		t.Errorf("Position.Y wasn't serialised properly, got %v", value)
	}
	binary.Read(buf, binary.LittleEndian, &value)
	if float32(g.Position[2]) != value {
		t.Errorf("Position.Z wasn't serialised properly, got %v", value)
	}

	// get INST_SET_ORIENTATION
	binary.Read(buf, binary.LittleEndian, &command)
	if command != byte(INST_SET_ORIENTATION) {
		t.Errorf("INST_SET_ORIENTATION wasn't serialised")
	}
	binary.Read(buf, binary.LittleEndian, &value)
	if float32(g.Orientation) != value {
		t.Errorf("Orientation wasn't serialised properly, got %v", value)
	}

	// INST_SET_TYPE
	binary.Read(buf, binary.LittleEndian, &command)
	if command != byte(INST_SET_TYPE) {
		t.Errorf("INST_SET_TYPE wasn't serialised")
	}
	binary.Read(buf, binary.LittleEndian, &value)
	if float32(g.Model) != value {
		t.Errorf("Model wasn't serialised properly, got %v", value)
	}

	// INST_SET_SCALE
	binary.Read(buf, binary.LittleEndian, &command)
	if command != byte(INST_SET_SCALE) {
		t.Errorf("INST_SET_TYPE wasn't serialised")
	}
	binary.Read(buf, binary.LittleEndian, &value)
	if float32(g.Scale[0]) != value {
		t.Errorf("Scale[0] wasn't serialised properly, got %v", value)
	}
	binary.Read(buf, binary.LittleEndian, &value)
	if float32(g.Scale[1]) != value {
		t.Errorf("Scale[1] wasn't serialised properly, got %v", value)
	}
	binary.Read(buf, binary.LittleEndian, &value)
	if float32(g.Scale[2]) != value {
		t.Errorf("Scale[2] wasn't serialised properly, got %v", value)
	}

}
