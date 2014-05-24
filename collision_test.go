package main

import (
	. "gopkg.in/check.v1"
)

type CollisionTestSuite struct{}

var _ = Suite(&CollisionTestSuite{})

func (s *CollisionTestSuite) TestCircleVsCircleMiss(c *C) {
	contact := &Collision{}
	contact.a = NewEntity()
	contact.a.Position = &Vector3{0, 0}
	contact.a.geometry = &Circle{Radius: 5}
	contact.b = NewEntity()
	contact.b.Position = &Vector3{10, 0}
	contact.b.geometry = &Circle{Radius: 5}
	collider := &CollisionDetector{}
	collider.CircleVsCircle(contact)
	c.Assert(contact.IsIntersecting, Equals, false)
	c.Assert(contact.penetration, Equals, float64(0))
	c.Assert(contact.normal, DeepEquals, &Vector3{-1, 0, 0})
}

func (s *CollisionTestSuite) TestCircleVsCircleHit(c *C) {
	contact := &Collision{}
	contact.a = NewEntity()
	contact.a.Position = &Vector3{0, 0}
	contact.a.geometry = &Circle{Radius: 5}
	contact.b = NewEntity()
	contact.b.Position = &Vector3{9, 0}
	contact.b.geometry = &Circle{Radius: 5}
	collider := &CollisionDetector{}
	collider.CircleVsCircle(contact)
	c.Assert(contact.IsIntersecting, Equals, true)
	c.Assert(contact.penetration, Equals, float64(1))
	c.Assert(contact.normal, DeepEquals, &Vector3{-1, 0, 0})
}

func (s *CollisionTestSuite) TestNoCircleCollision(c *C) {
	collider := &CollisionDetector{}
	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.geometry = &Circle{Radius: 5}

	b := &Entity{}
	b.Position = &Vector3{10, 0, 0}
	b.geometry = &Circle{Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, false)
	c.Assert(pair.penetration, Equals, float64(0))
	c.Assert(pair.normal, DeepEquals, &Vector3{-1, 0, 0})
}

func (s *CollisionTestSuite) TestCircleCollision(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.geometry = &Circle{Radius: 5}

	b := &Entity{}
	b.Position = &Vector3{9, 0, 0}
	b.geometry = &Circle{Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)
	c.Assert(pair.penetration, Equals, float64(1))
	c.Assert(pair.normal, DeepEquals, &Vector3{-1, 0, 0})
}

func (s *CollisionTestSuite) TestAABBHit(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.geometry = &Rectangle{SizeX: 4, SizeY: 4, SizeZ: 4}
	a.geometry.(*Rectangle).ToWorld(a.Position)

	b := &Entity{}
	b.Position = &Vector3{1, 0, 0}
	b.geometry = &Rectangle{SizeX: 4, SizeY: 4, SizeZ: 4}
	b.geometry.(*Rectangle).ToWorld(b.Position)

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)
	c.Assert(pair.penetration, Equals, float64(3))
	c.Assert(pair.normal, DeepEquals, &Vector3{-1, 0, 0})
}

func (s *CollisionTestSuite) TestAABBNoHit(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.geometry = &Rectangle{SizeX: 4, SizeY: 4, SizeZ: 4}

	b := &Entity{}
	b.Position = &Vector3{10, 0, 0}
	b.geometry = &Rectangle{SizeX: 4, SizeY: 4, SizeZ: 4}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, false)
	c.Assert(pair.penetration, Equals, float64(0))
	c.Assert(pair.normal, DeepEquals, &Vector3{0, 0, 0})
}

func (s *CollisionTestSuite) TestCollisionResolve(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Velocity = &Vector3{10, 0, 0}
	a.physics = &ParticlePhysics{InvMass: 1, forces: &Vector3{}}
	a.Position = &Vector3{0, 0, 0}
	a.geometry = &Circle{Radius: 5}

	b := &Entity{}
	b.Velocity = &Vector3{0, 0, 0}
	b.physics = &ParticlePhysics{InvMass: 1, forces: &Vector3{}}
	b.Position = &Vector3{9, 0, 0}
	b.geometry = &Circle{Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)

	pair.restitution = 1
	pair.Resolve(1)

	c.Assert(a.Velocity, DeepEquals, &Vector3{0, 0, 0})
	c.Assert(a.Position, DeepEquals, &Vector3{-0.5, 0})
	c.Assert(b.Velocity, DeepEquals, &Vector3{10, 0, 0})
	c.Assert(b.Position, DeepEquals, &Vector3{9.5, 0})
}

func (s *CollisionTestSuite) TestCollisionResolveOpposite(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Velocity = &Vector3{5, 0, 0}
	a.physics = &ParticlePhysics{InvMass: 1, forces: &Vector3{}}
	a.Position = &Vector3{0, 0, 0}
	a.geometry = &Circle{Radius: 5}

	b := &Entity{}
	b.Velocity = &Vector3{-5, 0, 0}
	b.physics = &ParticlePhysics{InvMass: 1, forces: &Vector3{}}
	b.Position = &Vector3{7, 0, 0}
	b.geometry = &Circle{Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)

	pair.restitution = 1
	pair.Resolve(1)

	c.Assert(a.Velocity, DeepEquals, &Vector3{-5, 0, 0})
	c.Assert(a.Position, DeepEquals, &Vector3{-1.5, 0, 0})
	c.Assert(b.Velocity, DeepEquals, &Vector3{5, 0, 0})
	c.Assert(b.Position, DeepEquals, &Vector3{8.5, 0, 0})
}
