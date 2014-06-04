package main

import (
	. "gopkg.in/check.v1"
)

type PhysicsTestSuite struct{}

var _ = Suite(&PhysicsTestSuite{})

func (s *PhysicsTestSuite) TestNoForcesDontMove(c *C) {
	el := NewEntityList()
	ent := el.NewEntity()
	ent.physics = NewParticlePhysics(1)
	ent.Update(1)
	ent.Update(1)
	c.Assert(ent.Position, DeepEquals, &Vector3{0, 0})
}

func (s *PhysicsTestSuite) TestNoForcesMove(c *C) {
	el := NewEntityList()
	ent := el.NewEntity()
	ent.physics = NewParticlePhysics(1)
	ent.physics.(*ParticlePhysics).Damping = 1
	ent.physics.(*ParticlePhysics).AddForce(&Vector3{1, 0})
	ent.Update(1)
	ent.Update(1)
	c.Assert(ent.Position, DeepEquals, &Vector3{1, 0})
}

func (s *PhysicsTestSuite) TestToHeavyToMove(c *C) {
	el := NewEntityList()
	ent := el.NewEntity()
	ent.physics = NewParticlePhysics(0)
	ent.physics.(*ParticlePhysics).Damping = 1
	ent.physics.(*ParticlePhysics).AddForce(&Vector3{1, 0})
	ent.Update(1)
	ent.Update(1)
	c.Assert(ent.Position, DeepEquals, &Vector3{0, 0})
}
