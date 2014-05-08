package physics

import (
	v "github.com/stojg/vivere/vec"
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestNewSimulator(c *C) {
	obj := NewSimulator()
	c.Assert(obj, Equals, obj)
}

func (s *TestSuite) TestUpdateVelocityWithDuration(c *C) {
	sim := NewSimulator()

	e := NewTestKinematic(0, 0, 2, 0)
	testList := &EntityList{}
	testList.entitities = append(testList.entitities, e)

	sim.Update(testList, 1)
	c.Assert(e.Position(), DeepEquals, &v.Vec{2, 0})
	sim.Update(testList, 1)
	c.Assert(e.Position(), DeepEquals, &v.Vec{4, 0})
	sim.Update(testList, 2)
	c.Assert(e.Position(), DeepEquals, &v.Vec{8, 0})

}

func (s *TestSuite) TestUpdateAcceleration(c *C) {
	sim := NewSimulator()
	e := NewTestKinematic(0, 0, 0, 0)

	testList := &EntityList{}
	testList.entitities = append(testList.entitities, e)

	e.AddForce(&v.Vec{2, 0})
	sim.Update(testList, 1)
	c.Assert(e.Velocity(), DeepEquals, &v.Vec{2, 0}, Commentf("Velocity should have changed"))
	c.Assert(e.Position(), DeepEquals, &v.Vec{0, 0}, Commentf("Position shouldnt have changed"))

	sim.Update(testList, 1)
	c.Assert(e.Velocity(), DeepEquals, &v.Vec{2, 0}, Commentf("Velocity shouldnt have changed"))
	c.Assert(e.Position(), DeepEquals, &v.Vec{2, 0}, Commentf("Position should have changed"))

	e.AddForce(&v.Vec{-2, 0})
	sim.Update(testList, 1)
	c.Assert(e.Velocity(), DeepEquals, &v.Vec{0, 0}, Commentf("Velocity should have changed"))
	c.Assert(e.Position(), DeepEquals, &v.Vec{4, 0}, Commentf("Position should have changed"))

	sim.Update(testList, 1)
	c.Assert(e.Velocity(), DeepEquals, &v.Vec{0, 0}, Commentf("Velocity shouldnt have changed"))
	c.Assert(e.Position(), DeepEquals, &v.Vec{4, 0}, Commentf("Position shouldnt have changed"))

}
