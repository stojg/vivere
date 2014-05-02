package ai

import (
	e "github.com/stojg/vivere/engine"
	v "github.com/stojg/vivere/vec"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestUpdateForce(c *C) {
	obj := &Simple{}
	ent := e.NewEntity(1)
	obj.UpdateForce(ent, 1)
	c.Assert(ent.Forces(), DeepEquals, &v.Vec{20, 0})
}
