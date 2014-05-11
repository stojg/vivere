package main

import (
	v "github.com/stojg/vivere/vec"
	. "gopkg.in/check.v1"
)

type GeometryTestSuite struct{}

var _ = Suite(&GeometryTestSuite{})

func (s *GeometryTestSuite) TestCircleVsCircleMiss(c *C) {
	circleA := &Circle{Position : &v.Vec{0,0}, Radius: 5}
	circleB := &Circle{Position : &v.Vec{10,0}, Radius: 5}
	pen, normal := circleA.VsCircle(circleB)
	c.Assert(pen, Equals, float64(0))
	c.Assert(normal, DeepEquals, &v.Vec{-1,0})
}

func (s *GeometryTestSuite) TestCircleVsCircleHit(c *C) {
	circleA := &Circle{Position : &v.Vec{0,0}, Radius: 5}
	circleB := &Circle{Position : &v.Vec{9,0}, Radius: 5}
	pen, normal := circleA.VsCircle(circleB)
	c.Assert(pen, Equals, float64(1))
	c.Assert(normal, DeepEquals, &v.Vec{-1,0})
}
