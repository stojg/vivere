package main

import (
	. "gopkg.in/check.v1"
)

type GeometryTestSuite struct{}

var _ = Suite(&GeometryTestSuite{})

func (s *GeometryTestSuite) TestCircleVsCircleMiss(c *C) {
	circleA := &Circle{Position : &Vector3{0,0}, Radius: 5}
	circleB := &Circle{Position : &Vector3{10,0}, Radius: 5}
	pen, normal := circleA.VsCircle(circleB)
	c.Assert(pen, Equals, float64(0))
	c.Assert(normal, DeepEquals, &Vector3{-1,0,0})
}

func (s *GeometryTestSuite) TestCircleVsCircleHit(c *C) {
	circleA := &Circle{Position : &Vector3{0,0}, Radius: 5}
	circleB := &Circle{Position : &Vector3{9,0}, Radius: 5}
	pen, normal := circleA.VsCircle(circleB)
	c.Assert(pen, Equals, float64(1))
	c.Assert(normal, DeepEquals, &Vector3{-1,0,0})
}
