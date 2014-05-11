package main

import (
	v "github.com/stojg/vivere/vec"
)

type Collision struct{}

func (c *Collision) Detect(a *Entity, b *Entity) (cp *CollisionPair, hit bool) {
	cp = &CollisionPair{}
	cp.a = a
	cp.b = b

	cp.pen, cp.normal = a.geometry.Collision(b.geometry)
	cp.restitution = 0.5
	if cp.pen > 0 {
		hit = true
	}
	return
}

func (c *Collision) CircleCircle(a *Circle, b *Circle) (pen float64, normal *v.Vec) {
	distanceVec := a.Position.NewSub(b.Position)
	distance := distanceVec.Length()
	pen = a.Radius + b.Radius - distance
	normal = distanceVec.Normalize()
	return pen, normal
}

func (c *Collision) RectangleRectangle(a *Entity, b *Entity)  (pen float64, normal *v.Vec) {




	return
}

func (c *Collision) CircleLine(a *Entity, b *Entity)  (pen float64, normal *v.Vec) {

	return
}

func (c *Collision) CircleLineSegment(a *Entity, b *Entity) (pen float64, normal *v.Vec) {

	return
}

type CollisionPair struct {
	a           *Entity
	b           *Entity
	restitution float64
	pen         float64
	normal      *v.Vec
}

func (c *CollisionPair) CalculateSeparatingVelocity() float64 {
	relativeVel := v.Vec{}
	relativeVel.Copy(c.a.physics.(*ParticlePhysics).Velocity)
	if c.b != nil {
		relativeVel.Sub(c.b.physics.(*ParticlePhysics).Velocity)
	}
	return relativeVel.Dot(c.normal)
}

func (c *CollisionPair) Resolve(duration float64) {
	c.resolveVelocity(duration)
	c.resolveInterpenetration()
}

func (c *CollisionPair) resolveInterpenetration() {

	if c.pen <= 0 {
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

	movePerIMass := c.normal.NewScale(c.pen / totalInvMass)

	c.a.Position.Add(movePerIMass.NewScale(c.a.physics.(*ParticlePhysics).InvMass))
	if c.b != nil {
		c.b.Position.Add(movePerIMass.NewScale(-c.b.physics.(*ParticlePhysics).InvMass))
	}
}

func (c *CollisionPair) resolveVelocity(duration float64) {
	// Find the velocity in the direction of the contact normal
	separatingVelocity := c.CalculateSeparatingVelocity()

	// The objects are already separating, NOP
	if separatingVelocity > 0 {
		return
	}

	// Calculate the new separating velocity
	newSepVelocity := -separatingVelocity * c.restitution

	// Check the velocity build up due to acceleration only
	accCausedVelocity := &v.Vec{}
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
	// Both objects have infinite mass, so no velocity
	if totalInvMass == 0 {
		return
	}

	var impulse float64
	impulse = deltaVelocity / totalInvMass

	var impulsePerIMass *v.Vec
	impulsePerIMass = c.normal.NewScale(impulse)

	temp := impulsePerIMass.NewScale(c.a.physics.(*ParticlePhysics).InvMass)
	c.a.physics.(*ParticlePhysics).Velocity.Add(temp)
	if c.b != nil {
		temp = impulsePerIMass.NewScale(-c.b.physics.(*ParticlePhysics).InvMass)
		c.b.physics.(*ParticlePhysics).Velocity.Add(temp)
	}
}
