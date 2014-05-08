package main

import (
	v "github.com/stojg/vivere/vec"
	"math"
)

type Physics struct{}

func (c *Physics) Update(e *Entity, elapsed float64) {

}

type ParticlePhysics struct {
	Physics
	InvMass  float64
	forces   *v.Vec
	Velocity *v.Vec
	damping  float64
}

func NewParticlePhysics() *ParticlePhysics {
	p := &ParticlePhysics{}
	p.forces = &v.Vec{}
	p.Velocity = &v.Vec{}
	p.InvMass = 1 / 1
	p.damping = 0.999
	return p
}

func (c *ParticlePhysics) Update(entity *Entity, elapsed float64) {
	if c.InvMass == 0 {
		return
	}
	entity.Position.AddScaledVector(c.Velocity, elapsed)
	c.Velocity.AddScaledVector(c.forces, elapsed)
	c.Velocity.Scale(math.Pow(c.damping, elapsed))
	c.ClearForces()
}

func (p *ParticlePhysics) AddForce(force *v.Vec) {
	p.forces.Add(force)
}

func (p *ParticlePhysics) ClearForces() {
	p.forces.Clear()
}
