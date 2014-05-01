package physics

import (
	v "github.com/stojg/vivere/vec"
)

type Body struct {
	shape      Shape
	position   v.Vec
	velocity   v.Vec
	forces     v.Vec
	rotation   float64
	mass       float64
	invMass    float64
	inertia    float64
	invInertia float64
	damping    float64
}

func (body *Body) Position() *v.Vec {
	return &body.position
}

func (body *Body) Velocity() *v.Vec {
	return &body.velocity
}

func (body *Body) Rotation() float64 {
	return body.rotation
}

func (body *Body) SetRotation(r float64) {
	body.rotation = r
}

func (body *Body) InvMass() float64 {
	if body.invMass == 0 {
		body.invMass = 1 / body.mass
	}
	return body.invMass
}

func (body *Body) SetMass(m float64) *Body {
	body.mass = m
	body.invMass = 1 / body.mass
	return body
}

func (body *Body) SetInertia(i float64) *Body {
	body.inertia = i
	body.invInertia = 1 / body.inertia
	return body
}

func (body *Body) Forces() *v.Vec {
	return &body.forces
}

func (body *Body) AddForce(vec *v.Vec) {
	body.forces.Add(vec)
}

func (body *Body) ClearForces() {
	body.forces.Clear()
}

func (body *Body) Damping() float64 {
	return body.damping
}

func (body *Body) SetDamping(damping float64) {
	body.damping = damping
}

func (body *Body) Shape() Shape {
	return body.shape
}

func (body *Body) SetShape(s Shape) {
	body.shape = s
}
