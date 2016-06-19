package main

import (
	. "gopkg.in/check.v1"
	"testing"
)

type CollisionTestSuite struct{}

var _ = Suite(&CollisionTestSuite{})

func (s *CollisionTestSuite) TestCircleVsCircleMiss(c *C) {
	collision := &Collision{}
	collision.a = NewEntity()
	collision.a.Position = &Vector3{0, 0}
	collision.a.Geometry = &Circle{Radius: 4}
	collision.b = NewEntity()
	collision.b.Position = &Vector3{10, 0}
	collision.b.Geometry = &Circle{Radius: 5}
	collision.normal = &Vector3{}
	collider := &CollisionDetector{}
	collider.CircleVsCircle(collision)
	c.Assert(collision.IsIntersecting, Equals, false)
	c.Assert(collision.penetration, Equals, float64(0))
	c.Assert(collision.normal, DeepEquals, &Vector3{})
}

func BenchmarkCircleVsCircleMiss(testing *testing.B) {
	collision := &Collision{}
	collision.a = NewEntity()
	collision.a.Position = &Vector3{0, 0}
	collision.a.Geometry = &Circle{Radius: 4}
	collision.b = NewEntity()
	collision.b.Position = &Vector3{10, 0}
	collision.b.Geometry = &Circle{Radius: 5}
	collision.normal = &Vector3{}
	collider := &CollisionDetector{}

	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.CircleVsCircle(collision)
	}
}

func (s *CollisionTestSuite) TestCircleVsCircleHit(c *C) {
	collision := &Collision{}
	collision.a = NewEntity()
	collision.a.Position = &Vector3{0, 0}
	collision.a.Geometry = &Circle{Radius: 5}
	collision.b = NewEntity()
	collision.b.Position = &Vector3{9, 0}
	collision.b.Geometry = &Circle{Radius: 5}
	collider := &CollisionDetector{}
	collider.CircleVsCircle(collision)
	c.Assert(collision.IsIntersecting, Equals, true)
	c.Assert(collision.penetration, Equals, float64(1))
	c.Assert(collision.normal, DeepEquals, &Vector3{-1, 0, 0})
}

func BenchmarkCircleVsCircleHit(testing *testing.B) {
	collision := &Collision{}
	collision.a = NewEntity()
	collision.a.Position = &Vector3{0, 0}
	collision.a.Geometry = &Circle{Radius: 5}
	collision.b = NewEntity()
	collision.b.Position = &Vector3{9, 0}
	collision.b.Geometry = &Circle{Radius: 5}
	collider := &CollisionDetector{}

	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.CircleVsCircle(collision)
	}
}

func BenchmarkAabbVsCircleMiss(testing *testing.B) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 4}
	b := &Entity{}
	b.Position = &Vector3{10, 0, 0}
	b.Geometry = &Rectangle{HalfSize: Vector3{2, 2, 2}}
	collision := &Collision{a: b, b: a, restitution: 0.5, normal: &Vector3{}}
	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.RectangleVsCircle(collision)
	}
}

func (s *CollisionTestSuite) TestCircleVsAABBMiss(c *C) {
	collision := &Collision{}
	collision.a = NewEntity()
	collision.a.Position = &Vector3{0, 0}
	collision.a.Geometry = &Circle{Radius: 5}
	collision.b = NewEntity()
	collision.b.Position = &Vector3{9, 0}
	collision.b.Geometry = &Rectangle{HalfSize: Vector3{2, 2, 2}}
	collider := &CollisionDetector{}
	collider.CircleVsRectangle(collision)
	c.Assert(collision.IsIntersecting, Equals, false)
	c.Assert(collision.penetration, Equals, float64(0))
	c.Assert(collision.normal, DeepEquals, &Vector3{0, 0, 0})
}

func BenchmarkCircleVsAABBMiss(testing *testing.B) {
	collision := &Collision{}
	collision.a = NewEntity()
	collision.a.Position = &Vector3{0, 0}
	collision.a.Geometry = &Circle{Radius: 5}
	collision.b = NewEntity()
	collision.b.Position = &Vector3{9, 0}
	collision.b.Geometry = &Rectangle{HalfSize: Vector3{2, 2, 2}}

	collider := &CollisionDetector{}

	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.CircleVsRectangle(collision)
		collision.a, collision.b = collision.b, collision.a
	}
}

func BenchmarkAabbVsCircleHit(testing *testing.B) {
	collision := &Collision{}
	collision.b = NewEntity()
	collision.b.Position = &Vector3{0, 0}
	collision.b.Geometry = &Circle{Radius: 5}
	collision.a = NewEntity()
	collision.a.Position = &Vector3{5, 0}
	collision.a.Geometry = &Rectangle{HalfSize: Vector3{2, 2, 2}}
	collider := &CollisionDetector{}
	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.RectangleVsCircle(collision)
	}
}

func BenchmarkAabbVsAabbMiss(testing *testing.B) {
	collision := &Collision{}

	collision.a = &Entity{}
	collision.a.Position = &Vector3{0, 0, 0}
	collision.a.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}

	collision.b = &Entity{}
	collision.b.Position = &Vector3{5, 0, 0}
	collision.b.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}

	collider := &CollisionDetector{}
	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.RectangleVsRectangle(collision)
	}
}

func BenchmarkAabbVsAabbHit(testing *testing.B) {
	collision := &Collision{}

	collision.a = &Entity{}
	collision.a.Position = &Vector3{0, 0, 0}
	collision.a.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}

	collision.b = &Entity{}
	collision.b.Position = &Vector3{3, 0, 0}
	collision.b.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}

	collider := &CollisionDetector{}
	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.RectangleVsRectangle(collision)
	}
}

func (s *CollisionTestSuite) TestDetectCircleVsCircleMiss(c *C) {
	collider := &CollisionDetector{}
	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 4}

	b := &Entity{}
	b.Position = &Vector3{10, 0, 0}
	b.Geometry = &Circle{Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, false)
	c.Assert(pair.penetration, Equals, float64(0))
	c.Assert(pair.normal, DeepEquals, &Vector3{})
}

func BenchmarkDetectCircleVsCircleMiss(testing *testing.B) {
	collider := &CollisionDetector{}
	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 4}

	b := &Entity{}
	b.Position = &Vector3{10, 0, 0}
	b.Geometry = &Circle{Radius: 5}

	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.Detect(a, b)
	}
}

func (s *CollisionTestSuite) TestDetectCircleVsCircleHit(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 5}

	b := &Entity{}
	b.Position = &Vector3{9, 0, 0}
	b.Geometry = &Circle{Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)
	c.Assert(pair.penetration, Equals, float64(1))
	c.Assert(pair.normal, DeepEquals, &Vector3{-1, 0, 0})
}

func BenchmarkDetectCircleVsCircleHit(testing *testing.B) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 5}

	b := &Entity{}
	b.Position = &Vector3{9, 0, 0}
	b.Geometry = &Circle{Radius: 5}

	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.Detect(a, b)
	}
}

func (s *CollisionTestSuite) TestDetectAabbVsAabbHit(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}
	a.Geometry.(*Rectangle).ToWorld(a.Position)

	b := &Entity{}
	b.Position = &Vector3{5, 0, 0}
	b.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}
	b.Geometry.(*Rectangle).ToWorld(b.Position)

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)
	c.Assert(pair.penetration, Equals, float64(3.0029999999999997))
	c.Assert(pair.normal, DeepEquals, &Vector3{-1, 0, 0})
}

func BenchmarkDetectAabbVsAabbHit(testing *testing.B) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}
	a.Geometry.(*Rectangle).ToWorld(a.Position)

	b := &Entity{}
	b.Position = &Vector3{5, 0, 0}
	b.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}
	b.Geometry.(*Rectangle).ToWorld(b.Position)

	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.Detect(a, b)
	}
}

func (s *CollisionTestSuite) TestDetectAabbVsAabbMiss(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}

	b := &Entity{}
	b.Position = &Vector3{10, 0, 0}
	b.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, false)
	c.Assert(pair.penetration, Equals, float64(0))
	c.Assert(pair.normal, DeepEquals, &Vector3{0, 0, 0})
}

func BenchmarkDetectAabbVsAabbMiss(testing *testing.B) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}

	b := &Entity{}
	b.Position = &Vector3{10, 0, 0}
	b.Geometry = &Rectangle{HalfSize: Vector3{4, 4, 4}}

	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.Detect(a, b)
	}
}

func (s *CollisionTestSuite) TestDetectAABBvsCircleMiss(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 4}

	b := &Entity{}
	b.Position = &Vector3{10, 0, 0}
	b.Geometry = &Rectangle{HalfSize: Vector3{2, 2, 2}}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, false)
	c.Assert(pair.penetration, Equals, float64(0))
	c.Assert(pair.normal, DeepEquals, &Vector3{0, 0, 0})
}

func BenchmarkDetectAABBvsCircleMiss(testing *testing.B) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 4}

	b := &Entity{}
	b.Position = &Vector3{10, 0, 0}
	b.Geometry = &Rectangle{HalfSize: Vector3{2, 2, 2}}

	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.Detect(a, b)
	}
}

func (s *CollisionTestSuite) TestDetectAABBvsCircleHit(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 4}

	b := &Entity{}
	b.Position = &Vector3{4, 0, 0}
	b.Geometry = &Rectangle{HalfSize: Vector3{2, 2, 2}}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)
	c.Assert(pair.penetration, Equals, float64(-2))
	c.Assert(pair.normal, DeepEquals, &Vector3{1, 0, 0})
}

func BenchmarkDetectAABBvsCircleHit(testing *testing.B) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 4}

	b := &Entity{}
	b.Position = &Vector3{4, 0, 0}
	b.Geometry = &Rectangle{HalfSize: Vector3{2, 2, 2}}

	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		collider.Detect(a, b)
	}
}

func (s *CollisionTestSuite) TestCollisionResolve(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Velocity = &Vector3{10, 0, 0}
	a.Body = &RigidBody{InvMass: 1, forces: &Vector3{}}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 5}

	b := &Entity{}
	b.Velocity = &Vector3{0, 0, 0}
	b.Body = &RigidBody{InvMass: 1, forces: &Vector3{}}
	b.Position = &Vector3{9, 0, 0}
	b.Geometry = &Circle{Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)

	pair.restitution = 1
	pair.Resolve(1)

	c.Assert(a.Velocity, DeepEquals, &Vector3{0, 0, 0})
	c.Assert(a.Position, DeepEquals, &Vector3{-0.5, 0})
	c.Assert(b.Velocity, DeepEquals, &Vector3{10, 0, 0})
	c.Assert(b.Position, DeepEquals, &Vector3{9.5, 0})
}

func BenchmarkCollisionResolve(testing *testing.B) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Velocity = &Vector3{10, 0, 0}
	a.Body = &RigidBody{InvMass: 1, forces: &Vector3{}}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 5}

	b := &Entity{}
	b.Velocity = &Vector3{0, 0, 0}
	b.Body = &RigidBody{InvMass: 1, forces: &Vector3{}}
	b.Position = &Vector3{9, 0, 0}
	b.Geometry = &Circle{Radius: 5}

	pair, _ := collider.Detect(a, b)
	pair.restitution = 1
	testing.ResetTimer()
	for i := 0; i < testing.N; i++ {
		pair.Resolve(1)
	}
}

func (s *CollisionTestSuite) TestCollisionResolveOpposite(c *C) {
	collider := &CollisionDetector{}

	a := &Entity{}
	a.Velocity = &Vector3{5, 0, 0}
	a.Body = &RigidBody{InvMass: 1, forces: &Vector3{}}
	a.Position = &Vector3{0, 0, 0}
	a.Geometry = &Circle{Radius: 5}

	b := &Entity{}
	b.Velocity = &Vector3{-5, 0, 0}
	b.Body = &RigidBody{InvMass: 1, forces: &Vector3{}}
	b.Position = &Vector3{7, 0, 0}
	b.Geometry = &Circle{Radius: 5}

	pair, hit := collider.Detect(a, b)
	c.Assert(hit, Equals, true)

	pair.restitution = 1
	pair.Resolve(1)

	c.Assert(a.Velocity, DeepEquals, &Vector3{-5, 0, 0})
	c.Assert(a.Position, DeepEquals, &Vector3{-1.5, 0, 0})
	c.Assert(b.Velocity, DeepEquals, &Vector3{5, 0, 0})
	c.Assert(b.Position, DeepEquals, &Vector3{8.5, 0, 0})
}

func TestRaygun (t *testing.T) {
	collider := &CollisionDetector{}

	origin := &Vector3{}
	ray := &Vector3{1,0,0}

	collider.raycast(origin, ray)
}
