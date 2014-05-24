package main

import (
	"math"
)

type Physics struct{}

func (c *Physics) Update(e *Entity, elapsed float64) {

}

type ParticlePhysics struct {
	Physics
	InvMass float64
	forces  *Vector3
	Damping float64
}

func NewParticlePhysics() *ParticlePhysics {
	p := &ParticlePhysics{}
	p.forces = &Vector3{}
	p.InvMass = 1 / 1
	p.Damping = 0.999
	return p
}

func (c *ParticlePhysics) Update(entity *Entity, elapsed float64) {
	if c.InvMass == 0 {
		return
	}
	entity.Position.AddScaledVector(entity.Velocity, elapsed)
	entity.Velocity.AddScaledVector(c.forces, elapsed)
	entity.Velocity.Scale(math.Pow(c.Damping, elapsed))

	// clamp velocity
	if entity.Velocity.Length() > 160 {
		entity.Velocity.Normalize().Scale(160)
	}
}

func (p *ParticlePhysics) AddForce(force *Vector3) {
	p.forces.Add(force)
}

func (p *ParticlePhysics) ClearForces() {
	p.forces.Clear()
}

type RigidBody struct {
	ParticlePhysics
	inverseMass               float64
	linearDamping             float64
	Position                  *Vector3
	Orientation               *Quaternion
	velocity                  *Vector3
	rotation                  *Vector3
	transformMatrix           *Matrix4
	inverseInertiaTensor      *Matrix3
	inverseInertiaTensorWorld *Matrix3
	forceAccum                Vector3
	torqueAccum               Vector3
	isAwake                   bool
}

func (rb *RigidBody) SetInertiaTensor(inertiaTensor *Matrix3) {
	rb.inverseInertiaTensor.SetInverse(inertiaTensor)
}

func (rb *RigidBody) calculateDerivedData() {
	rb.Orientation.Normalize()
	rb.calculateTransformMatrix(rb.transformMatrix, rb.Position, rb.Orientation)
	rb.transformInertiaTensor(rb.inverseInertiaTensorWorld, rb.Orientation, rb.inverseInertiaTensor, rb.transformMatrix)
}

/**
 * Inline function that creates a transform matrix from a
 * position and orientation.
 */
func (rb *RigidBody) calculateTransformMatrix(transformMatrix *Matrix4, position *Vector3, orientation *Quaternion) {
	transformMatrix[0] = 1 - 2*orientation.j*orientation.j - 2*orientation.k*orientation.k
	transformMatrix[1] = 2*orientation.i*orientation.j - 2*orientation.r*orientation.k
	transformMatrix[2] = 2*orientation.i*orientation.k + 2*orientation.r*orientation.j
	transformMatrix[3] = position[0]

	transformMatrix[4] = 2*orientation.i*orientation.j + 2*orientation.r*orientation.k
	transformMatrix[5] = 1 - 2*orientation.i*orientation.i - 2*orientation.k*orientation.k
	transformMatrix[6] = 2*orientation.j*orientation.k - 2*orientation.r*orientation.i
	transformMatrix[7] = position[1]

	transformMatrix[8] = 2*orientation.i*orientation.k - 2*orientation.r*orientation.j
	transformMatrix[9] = 2*orientation.j*orientation.k + 2*orientation.r*orientation.i
	transformMatrix[10] = 1 - 2*orientation.i*orientation.i - 2*orientation.j*orientation.j
	transformMatrix[11] = position[1]
}

func (rb *RigidBody) AddForce(force *Vector3) {
	rb.forceAccum.Add(force)
	rb.isAwake = true
}

func (rb *RigidBody) Update(entity *Entity, elapsed float64) {
	rb.ClearAccumulators()
}

func (rb *RigidBody) ClearAccumulators() {
	rb.forces.Clear()
	rb.torqueAccum.Clear()
}

func (rb *RigidBody) AddForceAtBodyPoint(force, point *Vector3) {
	pt := rb.getPointInWorldSpace(point)
	rb.AddForceAtPoint(force, pt)
	rb.isAwake = true
}

func (rb *RigidBody) AddForceAtPoint(force, point *Vector3) {

	pt := point.Clone()
	pt.Sub(rb.Position)

	rb.forceAccum.Add(force)
	rb.torqueAccum.Add(pt.VectorProduct(force))
}

func (rb *RigidBody) getPointInWorldSpace(point *Vector3) *Vector3 {
	return rb.transformMatrix.TransformVector3(point)
}

/**
 * Internal function to do an intertia tensor transform by a quaternion.
 * Note that the implementation of this function was created by an
 * automated code-generator and optimizer.
 */
func (rb *RigidBody) transformInertiaTensor(iitWorld *Matrix3, q *Quaternion, iitBody *Matrix3, rotmat *Matrix4) {
	t4 := rotmat[0]*iitBody[0] + rotmat[1]*iitBody[3] + rotmat[2]*iitBody[6]
	t9 := rotmat[0]*iitBody[1] + rotmat[1]*iitBody[4] + rotmat[2]*iitBody[7]
	t14 := rotmat[0]*iitBody[2] + rotmat[1]*iitBody[5] + rotmat[2]*iitBody[8]
	t28 := rotmat[4]*iitBody[0] + rotmat[5]*iitBody[3] + rotmat[6]*iitBody[6]
	t33 := rotmat[4]*iitBody[1] + rotmat[5]*iitBody[4] + rotmat[6]*iitBody[7]
	t38 := rotmat[4]*iitBody[2] + rotmat[5]*iitBody[5] + rotmat[6]*iitBody[8]
	t52 := rotmat[8]*iitBody[0] + rotmat[9]*iitBody[3] + rotmat[10]*iitBody[6]
	t57 := rotmat[8]*iitBody[1] + rotmat[9]*iitBody[4] + rotmat[10]*iitBody[7]
	t62 := rotmat[8]*iitBody[2] + rotmat[9]*iitBody[5] + rotmat[10]*iitBody[8]

	iitWorld[0] = t4*rotmat[0] + t9*rotmat[1] + t14*rotmat[2]
	iitWorld[1] = t4*rotmat[4] + t9*rotmat[5] + t14*rotmat[6]
	iitWorld[2] = t4*rotmat[8] + t9*rotmat[9] + t14*rotmat[10]
	iitWorld[3] = t28*rotmat[0] + t33*rotmat[1] + t38*rotmat[2]
	iitWorld[4] = t28*rotmat[4] + t33*rotmat[5] + t38*rotmat[6]
	iitWorld[5] = t28*rotmat[8] + t33*rotmat[9] + t38*rotmat[10]
	iitWorld[6] = t52*rotmat[0] + t57*rotmat[1] + t62*rotmat[2]
	iitWorld[7] = t52*rotmat[4] + t57*rotmat[5] + t62*rotmat[6]
	iitWorld[8] = t52*rotmat[8] + t57*rotmat[9] + t62*rotmat[10]
}
