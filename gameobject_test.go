package main

import (
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestGetNewAndIdGeneration(c *C) {
	gol := NewGameObjectList()
	ent := gol.NewGameObject()
	c.Assert(ent.ID(), Equals, uint16(1))
	ent = gol.NewGameObject()
	c.Assert(ent.ID(), Equals, uint16(2))
	ent = gol.NewGameObject()
	c.Assert(ent.ID(), Equals, uint16(3))
}

func (s *TestSuite) TestList(c *C) {
	gol := NewGameObjectList()
	g := gol.NewGameObject()
	c.Assert(gol.Get(g.ID()), Equals, g)
	c.Assert(gol.Length(), Equals, 1)
	gol.Remove(g.ID())
	c.Assert(gol.Length(), Equals, 0)
}
