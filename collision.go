package main

import (
	"math"
	//	"fmt"
)

type CollisionDetector struct{}

func (c *CollisionDetector) Detect(a *Entity, b *Entity) (cp *Collision, hit bool) {

	cp = &Collision{}
	cp.a = a
	cp.b = b

	switch a.geometry.(type) {
	case *Circle:
		switch b.geometry.(type) {
		case *Circle:
			c.CircleVsCircle(cp)
		case *Rectangle:
			c.CircleVsRectangle(cp)
		}
	case *Rectangle:
		switch b.geometry.(type) {
		case *Rectangle:
			c.RectangleVsRectangle(cp)
		case *Circle:
			c.RectangleVsCircle(cp)
		}
	default:
		panic("unknown collision geometry")
	}
	cp.restitution = 0.5
	hit = cp.IsIntersecting
	return
}

func (c *CollisionDetector) CircleVsCircle(col *Collision) {
	a := col.a
	b := col.b
	cA := col.a.geometry.(*Circle)
	cB := col.b.geometry.(*Circle)
	distanceVec := a.Position.Clone().Sub(b.Position)
	distance := distanceVec.Length()
	col.penetration = cA.Radius + cB.Radius - distance
	if col.penetration > 0 {
		col.IsIntersecting = true
	}
	col.normal = distanceVec.Normalize()
}

func (c *CollisionDetector) CircleVsRectangle(contact *Collision) {
	contact.a, contact.b = contact.b, contact.a
	c.RectangleVsCircle(contact)
}

func (c *CollisionDetector) RectangleVsCircle(contact *Collision) {

	rectA := contact.a.geometry.(*Rectangle)
	rectA.ToWorld(contact.a.Position)

	contact.normal = &Vector3{}

	circleB := contact.b.geometry.(*Circle)
	p := contact.b.Position

	closestPoint := &Vector3{}
	for i := 0; i < 3; i++ {
		v := p[i]
		if v < rectA.MinPoint[i] {
			v = rectA.MinPoint[i]
		}
		if v > rectA.MaxPoint[i] {
			v = rectA.MaxPoint[i]
		}
		closestPoint[i] = v
	}

	distVec := closestPoint.Sub(p)
	if distVec.SquareLength() > circleB.Radius*circleB.Radius {
		return
	}
	contact.normal = distVec.Normalize()
	contact.penetration = distVec.Length() - circleB.Radius
	contact.IsIntersecting = true
}

func (c *CollisionDetector) RectangleVsRectangle(contact *Collision) {
	rectA := contact.a.geometry.(*Rectangle)
	rectB := contact.b.geometry.(*Rectangle)

	rectA.ToWorld(contact.a.Position)
	rectB.ToWorld(contact.b.Position)

	// [Minimum Translation Vector]
	mtvDistance := math.MaxFloat32 // Set current minimum distance (max float value so next value is always less)
	mtvAxis := &Vector3{}          // Axis along which to travel with the minimum distance

	// [Axes of potential separation]
	// • Each shape must be projected on these axes to test for intersection:
	// (1, 0, 0)                    A0 (= B0) [X Axis]
	// (0, 1, 0)                    A1 (= B1) [Y Axis]
	// (0, 0, 1)                    A1 (= B2) [Z Axis]
	contact.normal = &Vector3{}
	// [X Axis]
	if !c.TestAxisStatic(UnitX, rectA.MinPoint[0], rectA.MaxPoint[0], rectB.MinPoint[0], rectB.MaxPoint[0], mtvAxis, &mtvDistance) {
		return
	}

	// [Y Axis]
	if !c.TestAxisStatic(UnitY, rectA.MinPoint[1], rectA.MaxPoint[1], rectB.MinPoint[1], rectB.MaxPoint[1], mtvAxis, &mtvDistance) {
		return
	}

	// [Z Axis]
	if !c.TestAxisStatic(UnitZ, rectA.MinPoint[2], rectA.MaxPoint[2], rectB.MinPoint[2], rectB.MaxPoint[2], mtvAxis, &mtvDistance) {
		return
	}

	// We got a hit
	contact.IsIntersecting = true

	// Calculate Minimum Translation Vector (MTV) [normal * penetration]
	contact.normal = mtvAxis.Normalize()

	// Multiply the penetration depth by itself plus a small increment
	// When the penetration is resolved using the MTV, it will no longer intersect
	//contact.pen = float64(math.Sqrt(mtvDistance)) * 1.001;
	contact.penetration = mtvDistance
}

func (c *CollisionDetector) TestAxisStatic(axis Vector3, minA, maxA, minB, maxB float64, mtvAxis *Vector3, mtvDistance *float64) bool {

	// [Separating Axis Theorem]
	// • Two convex shapes only overlap if they overlap on all axes of separation
	// • In order to create accurate responses we need to find the collision vector (Minimum Translation Vector)
	// • Find if the two boxes intersect along a single axis
	// • Compute the intersection interval for that axis
	// • Keep the smallest intersection/penetration value
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

	// If that vector is smaller than our computed Minimum Translation Distance use that vector as our
	// current MTV distance
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

func (c *Collision) CalculateSeparatingVelocity() float64 {
	relativeVel := Vector3{}
	relativeVel.Copy(c.a.Velocity)
	if c.b != nil {
		relativeVel.Sub(c.b.Velocity)
	}
	return relativeVel.Dot(c.normal)
}

func (c *Collision) Resolve(duration float64) {
	c.resolveVelocity(duration)
	c.resolveInterpenetration()
}

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

func (c *Collision) resolveVelocity(duration float64) {
	// Find the velocity in the direction of the contact normal
	separatingVelocity := c.CalculateSeparatingVelocity()

	// The objects are already separating, NOP
	if separatingVelocity > 0 {
		return
	}

	// Calculate the new separating velocity
	newSepVelocity := -separatingVelocity * c.restitution

	// Check the velocity build up due to acceleration only
	accCausedVelocity := &Vector3{}
	accCausedVelocity.Copy(c.a.physics.(*ParticlePhysics).forces)
	if c.b != nil {
		accCausedVelocity.Sub(c.b.physics.(*ParticlePhysics).forces)
	}
	accCausedSepVelocity := accCausedVelocity.Dot(c.normal) * duration

	// if we have closing velocity due to acceleration buildup,
	// remove it from the new separating velocity
	if accCausedSepVelocity < 0 {
		newSepVelocity += c.restitution * accCausedSepVelocity
		// make sure that we haven't removed more than was
		// there to begin with
		if newSepVelocity < 0 {
			newSepVelocity = 0
		}
	}

	deltaVelocity := newSepVelocity - separatingVelocity

	totalInvMass := c.a.physics.(*ParticlePhysics).InvMass
	if c.b != nil {
		totalInvMass += c.b.physics.(*ParticlePhysics).InvMass
	}

	// Both objects have infinite mass, so no velocity change
	if totalInvMass == 0 {
		return
	}

	var impulse float64
	impulse = deltaVelocity / totalInvMass

	var impulsePerIMass *Vector3
	impulsePerIMass = c.normal.Clone().Scale(impulse)

	temp := impulsePerIMass.Clone().Scale(c.a.physics.(*ParticlePhysics).InvMass)
	c.a.Velocity.Add(temp)
	if c.b != nil {
		temp = impulsePerIMass.Clone().Scale(-c.b.physics.(*ParticlePhysics).InvMass)
		c.b.Velocity.Add(temp)
	}
}
