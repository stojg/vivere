package main

type CollisionPair struct {
	a           *Entity
	b           *Entity
	restitution float64
	pen         float64
	normal      *Vector3
}

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

func (c *CollisionPair) CalculateSeparatingVelocity() float64 {
	relativeVel := Vector3{}
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

	movePerIMass := c.normal.Clone().Scale(c.pen / totalInvMass)

	c.a.Position.Add(movePerIMass.Clone().Scale(c.a.physics.(*ParticlePhysics).InvMass))
	if c.b != nil {
		c.b.Position.Add(movePerIMass.Clone().Scale(-c.b.physics.(*ParticlePhysics).InvMass))
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
	// Both objects have infinite mass, so no velocity
	if totalInvMass == 0 {
		return
	}

	var impulse float64
	impulse = deltaVelocity / totalInvMass

	var impulsePerIMass *Vector3
	impulsePerIMass = c.normal.Clone().Scale(impulse)

	temp := impulsePerIMass.Clone().Scale(c.a.physics.(*ParticlePhysics).InvMass)
	c.a.physics.(*ParticlePhysics).Velocity.Add(temp)
	if c.b != nil {
		temp = impulsePerIMass.Clone().Scale(-c.b.physics.(*ParticlePhysics).InvMass)
		c.b.physics.(*ParticlePhysics).Velocity.Add(temp)
	}
}
