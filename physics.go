package main

import (
	"math"
)

type Physics struct{}

func (c *Physics) Update(e *Entity, elapsed float64) {

}

type ParticlePhysics struct {
	Physics
	InvMass  float64
	forces   *Vector3
	Velocity *Vector3
	Damping  float64
}

func NewParticlePhysics() *ParticlePhysics {
	p := &ParticlePhysics{}
	p.forces = &Vector3{}
	p.Velocity = &Vector3{}
	p.InvMass = 1 / 1
	p.Damping = 0.999
	return p
}

func (c *ParticlePhysics) Update(entity *Entity, elapsed float64) {
	if c.InvMass == 0 {
		return
	}
	entity.Position.AddScaledVector(c.Velocity, elapsed)
	c.Velocity.AddScaledVector(c.forces, elapsed)
	c.Velocity.Scale(math.Pow(c.Damping, elapsed))

	// clamp velocity
	if c.Velocity.Length() > 160 {
		c.Velocity.Normalize().Scale(160)
	}
}

func (p *ParticlePhysics) AddForce(force *Vector3) {
	p.forces.Add(force)
}

func (p *ParticlePhysics) ClearForces() {
	p.forces.Clear()
}
