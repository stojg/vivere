package ai

import (
	v "github.com/stojg/vivere/vec"
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestWorldDimension(c *C) {
	obj := &Perception{}
	c.Assert(obj.WorldDimension(), DeepEquals, &v.Vec{1000, 600})
}
