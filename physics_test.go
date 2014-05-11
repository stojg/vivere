package main

import (
	v "github.com/stojg/vivere/vec"
	. "gopkg.in/check.v1"
)

type PhysicsTestSuite struct{}

var _ = Suite(&PhysicsTestSuite{})

func (s *PhysicsTestSuite) TestNoForcesDontMove(c *C) {
	el := NewEntityList()
	ent := el.NewEntity()
	ent.physics = NewParticlePhysics()
	ent.physics.(*ParticlePhysics).InvMass = 1
	ent.Update(1)
	ent.Update(1)
	c.Assert(ent.Position, DeepEquals, &v.Vec{0,0})
}

func (s *PhysicsTestSuite) TestNoForcesMove(c *C) {
	el := NewEntityList()
	ent := el.NewEntity()
	ent.physics = NewParticlePhysics()
	ent.physics.(*ParticlePhysics).InvMass = 1
	ent.physics.(*ParticlePhysics).Damping = 1
	ent.physics.(*ParticlePhysics).AddForce(&v.Vec{1,0})
	ent.Update(1)
	ent.Update(1)
	c.Assert(ent.Position, DeepEquals, &v.Vec{1,0})
}

func (s *PhysicsTestSuite) TestToHeavyToMove(c *C) {
	el := NewEntityList()
	ent := el.NewEntity()
	ent.physics = NewParticlePhysics()
	ent.physics.(*ParticlePhysics).InvMass = 0
	ent.physics.(*ParticlePhysics).Damping = 1
	ent.physics.(*ParticlePhysics).AddForce(&v.Vec{1,0})
	ent.Update(1)
	ent.Update(1)
	c.Assert(ent.Position, DeepEquals, &v.Vec{0,0})
}
