package main

import (
	"math"
)

type CollisionDetector struct{}

func (c *CollisionDetector) Detect(a *Entity, b *Entity) (collision *Collision, hit bool) {

	// @todo hardcoded restitution
	collision = &Collision{a: a, b: b, restitution: 0.5, normal: &Vector3{}}

	switch a.geometry.(type) {
	case *Circle:
		switch b.geometry.(type) {
		case *Circle:
			c.CircleVsCircle(collision)
		case *Rectangle:
			c.CircleVsRectangle(collision)
		}
	case *Rectangle:
		switch b.geometry.(type) {
		case *Rectangle:
			c.RectangleVsRectangle(collision)
		case *Circle:
			c.RectangleVsCircle(collision)
		}
	default:
		panic("unknown collision geometry")
	}
	hit = collision.IsIntersecting
	return
}

func (c *CollisionDetector) CircleVsCircle(contact *Collision) {
	cA := contact.a.geometry.(*Circle)
	cB := contact.b.geometry.(*Circle)
	distanceVec := contact.a.Position.Clone().Sub(contact.b.Position)

	// Early out to avoid expensive sqrt
	if distanceVec.SquareLength() > (cA.Radius+cB.Radius)*(cA.Radius+cB.Radius) {
		return
	}
	contact.penetration = cA.Radius + cB.Radius - distanceVec.Length()
	contact.normal = distanceVec.Normalize()
	contact.IsIntersecting = true
}

func (c *CollisionDetector) CircleVsRectangle(contact *Collision) {
	contact.a, contact.b = contact.b, contact.a
	c.RectangleVsCircle(contact)
}

func (c *CollisionDetector) RectangleVsCircle(contact *Collision) {
	rA := contact.a.geometry.(*Rectangle)
	rA.ToWorld(contact.a.Position)
	cB := contact.b.geometry.(*Circle)
	contact.normal = &Vector3{}

	closestPoint := &Vector3{}
	for i := 0; i < 3; i++ {
		closestPoint[i] = contact.b.Position[i]
		if closestPoint[i] < rA.MinPoint[i] {
			closestPoint[i] = rA.MinPoint[i]
		} else if closestPoint[i] > rA.MaxPoint[i] {
			closestPoint[i] = rA.MaxPoint[i]
		}
	}

	distanceVec := closestPoint.Sub(contact.b.Position)
	// Early out to avoid expensive sqrt
	if distanceVec.SquareLength() > cB.Radius*cB.Radius {
		return
	}
	contact.penetration = distanceVec.Length() - cB.Radius
	contact.normal = distanceVec.Normalize()
	contact.IsIntersecting = true
}

func (c *CollisionDetector) RectangleVsRectangle(contact *Collision) {
	rA := contact.a.geometry.(*Rectangle)
	rB := contact.b.geometry.(*Rectangle)

	rA.ToWorld(contact.a.Position)
	rB.ToWorld(contact.b.Position)

	// [Minimum Translation Vector]
	mtvDistance := math.MaxFloat32 // Set current minimum distance (max float value so next value is always less)
	mtvAxis := &Vector3{}          // Axis along which to travel with the minimum distance

	// [Axes of potential separation]
	// [X Axis]
	if !c.testAxisSeparation(UnitX, rA.MinPoint[0], rA.MaxPoint[0], rB.MinPoint[0], rB.MaxPoint[0], mtvAxis, &mtvDistance) {
		return
	}

	// [Y Axis]
	if !c.testAxisSeparation(UnitY, rA.MinPoint[1], rA.MaxPoint[1], rB.MinPoint[1], rB.MaxPoint[1], mtvAxis, &mtvDistance) {
		return
	}

	// [Z Axis]
	if !c.testAxisSeparation(UnitZ, rA.MinPoint[2], rA.MaxPoint[2], rB.MinPoint[2], rB.MaxPoint[2], mtvAxis, &mtvDistance) {
		return
	}

	contact.penetration = mtvDistance * 1.001
	contact.normal = mtvAxis.Normalize()
	contact.IsIntersecting = true
}

// TestAxisStatic checks if two axis overlaps and in that case calculates how much
// * Two convex shapes only overlap if they overlap on all axes of separation
// * In order to create accurate responses we need to find the
//    collision vector (Minimum Translation Vector)
// * Find if the two boxes intersect along a single axis
// * Compute the intersection interval for that axis
// * Keep the smallest intersection/penetration value
func (c *CollisionDetector) testAxisSeparation(axis Vector3, minA, maxA, minB, maxB float64, mtvAxis *Vector3, mtvDistance *float64) bool {

	axisLengthSquared := axis.Dot(&axis)

	// If the axis is degenerate then ignore
	if axisLengthSquared < 1.0e-8 {
		return true
	}

	// Calculate the two possible overlap ranges
	// Either we overlap on the left or the right sides
	d0 := (maxB - minA) // 'Left' side
	d1 := (maxA - minB) // 'Right' side

	// Intervals do not overlap, so no intersection
	if d0 <= 0.0 || d1 <= 0.0 {
		return false
	}

	var overlap float64
	// Find out if we overlap on the 'right' or 'left' of the object.
	if d0 < d1 {
		overlap = d0
	} else {
		overlap = -d1
	}

	// The mtd vector for that axis
	sep := axis.Scale(overlap / axisLengthSquared)

	// The mtd vector length squared
	sepLengthSquared := sep.Dot(sep)

	// If that vector is smaller than our computed Minimum Translation
	// Distance use that vector as our current MTV distance
	if sepLengthSquared < *mtvDistance {
		*mtvDistance = math.Sqrt(sepLengthSquared)
		mtvAxis.Copy(sep)
	}
	return true
}

type Collision struct {
	a              *Entity
	b              *Entity
	restitution    float64
	penetration    float64
	normal         *Vector3
	IsIntersecting bool
}

func (collision *Collision) SeparatingVelocity() float64 {
	relativeVel := collision.a.Velocity.Clone()
	if collision.b != nil {
		relativeVel.Sub(collision.b.Velocity)
	}
	return relativeVel.Dot(collision.normal)
}

func (c *Collision) Resolve(duration float64) {
	c.resolveVelocity(duration)
	c.resolveInterpenetration()
}

// resolveInterpenetration separates two objects that has penetrated
func (c *Collision) resolveInterpenetration() {

	if c.penetration <= 0 {
		return
	}

	totalInvMass := c.a.physics.(*ParticlePhysics).InvMass
	if c.b != nil {
		totalInvMass += c.b.physics.(*ParticlePhysics).InvMass
	}
	// Both objects have infinite mass, so no velocity
	if totalInvMass == 0 {
		return
	}

	movePerIMass := c.normal.Clone().Scale(c.penetration / totalInvMass)

	c.a.Position.Add(movePerIMass.Clone().Scale(c.a.physics.(*ParticlePhysics).InvMass))
	if c.b != nil {
		c.b.Position.Add(movePerIMass.Clone().Scale(-c.b.physics.(*ParticlePhysics).InvMass))
	}
}

// resolveVelocity calculates the new velocity that is the result of the collision
func (collision *Collision) resolveVelocity(duration float64) {
	// Find the velocity in the direction of the contact normal
	separatingVelocity := collision.SeparatingVelocity()

	// The objects are already separating, NOP
	if separatingVelocity > 0 {
		return
	}

	// Calculate the new separating velocity
	newSepVelocity := -separatingVelocity * collision.restitution

	// Check the velocity build up due to acceleration only
	accCausedVelocity := collision.a.physics.(*ParticlePhysics).forces.Clone()
	if collision.b != nil {
		accCausedVelocity.Sub(collision.b.physics.(*ParticlePhysics).forces)
	}

	// If we have closing velocity due to acceleration buildup,
	// remove it from the new separating velocity
	accCausedSepVelocity := accCausedVelocity.Dot(collision.normal) * duration
	if accCausedSepVelocity < 0 {
		newSepVelocity += collision.restitution * accCausedSepVelocity
		// make sure that we haven't removed more than was there to begin with
		if newSepVelocity < 0 {
			newSepVelocity = 0
		}
	}

	deltaVelocity := newSepVelocity - separatingVelocity

	totalInvMass := collision.a.physics.(*ParticlePhysics).InvMass
	if collision.b != nil {
		totalInvMass += collision.b.physics.(*ParticlePhysics).InvMass
	}

	// Both objects have infinite mass, so they can't actually move
	if totalInvMass == 0 {
		return
	}

	impulsePerIMass := collision.normal.Clone().Scale(deltaVelocity / totalInvMass)

	velocityChangeA := impulsePerIMass.Clone().Scale(collision.a.physics.(*ParticlePhysics).InvMass)
	collision.a.Velocity.Add(velocityChangeA)
	if collision.b != nil {
		velocityChangeB := impulsePerIMass.Clone().Scale(-collision.b.physics.(*ParticlePhysics).InvMass)
		collision.b.Velocity.Add(velocityChangeB)
	}
}
