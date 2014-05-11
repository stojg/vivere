package main

import (
	v "github.com/stojg/vivere/vec"
	. "gopkg.in/check.v1"
)

type CollisionTestSuite struct{}

var _ = Suite(&CollisionTestSuite{})

func (s *CollisionTestSuite) TestNoCircleCollision(c *C) {
	collider := &Collision{}

	a := &Entity{}
	a.physics = &ParticlePhysics{}
	a.Position = &v.Vec{0, 0}
	a.geometry = &Circle{Position: a.Position, Radius: 5}

	b := &Entity{}
	b.physics = &ParticlePhysics{}
	b.Position = &v.Vec{10, 0}
	b.geometry = &Circle{Position: b.Position, Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, false)
	c.Assert(pair.pen, Equals, float64(0))
	c.Assert(pair.normal, DeepEquals, &v.Vec{-1, 0})
}

func (s *CollisionTestSuite) TestCircleCollision(c *C) {
	collider := &Collision{}

	a := &Entity{}
	a.physics = &ParticlePhysics{}
	a.Position = &v.Vec{0, 0}
	a.geometry = &Circle{Position: a.Position, Radius: 5}

	b := &Entity{}
	b.physics = &ParticlePhysics{}
	b.Position = &v.Vec{9, 0}
	b.geometry = &Circle{Position: b.Position, Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)
	c.Assert(pair.pen, Equals, float64(1))
	c.Assert(pair.normal, DeepEquals, &v.Vec{-1, 0})
}

func (s *CollisionTestSuite) TestCollisionResolve(c *C) {
	collider := &Collision{}

	a := &Entity{}
	a.physics = &ParticlePhysics{Velocity: &v.Vec{10, 0}, InvMass: 1, forces: &v.Vec{}}
	a.Position = &v.Vec{0, 0}
	a.geometry = &Circle{Position: a.Position, Radius: 5}

	b := &Entity{}
	b.physics = &ParticlePhysics{Velocity: &v.Vec{0, 0}, InvMass: 1, forces: &v.Vec{}}
	b.Position = &v.Vec{9, 0}
	b.geometry = &Circle{Position: b.Position, Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)

	pair.restitution = 1
	pair.Resolve(1)

	c.Assert(a.physics.(*ParticlePhysics).Velocity, DeepEquals, &v.Vec{0, 0})
	c.Assert(a.Position, DeepEquals, &v.Vec{-0.5, 0})
	c.Assert(b.physics.(*ParticlePhysics).Velocity, DeepEquals, &v.Vec{10, 0})
	c.Assert(b.Position, DeepEquals, &v.Vec{9.5, 0})
}

func (s *CollisionTestSuite) TestCollisionResolveOpposite(c *C) {
	collider := &Collision{}

	a := &Entity{}
	a.physics = &ParticlePhysics{Velocity: &v.Vec{5, 0}, InvMass: 1, forces: &v.Vec{}}
	a.Position = &v.Vec{0, 0}
	a.geometry = &Circle{Position: a.Position, Radius: 5}

	b := &Entity{}
	b.physics = &ParticlePhysics{Velocity: &v.Vec{-5, 0}, InvMass: 1, forces: &v.Vec{}}
	b.Position = &v.Vec{7, 0}
	b.geometry = &Circle{Position: b.Position, Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)

	pair.restitution = 1
	pair.Resolve(1)

	c.Assert(a.physics.(*ParticlePhysics).Velocity, DeepEquals, &v.Vec{-5, 0})
	c.Assert(a.Position, DeepEquals, &v.Vec{-1.5, 0})
	c.Assert(b.physics.(*ParticlePhysics).Velocity, DeepEquals, &v.Vec{5, 0})
	c.Assert(b.Position, DeepEquals, &v.Vec{8.5, 0})
}
