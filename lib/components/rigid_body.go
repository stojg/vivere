package components

import (
	. "github.com/stojg/vector"
	"sync"
)

func NewRigidBodyManager() *RigidBodyList {
	return &RigidBodyList{
		entity: make(map[*Entity]*RigidBody),
	}
}

type RigidBodyList struct {
	sync.Mutex
	entity map[*Entity]*RigidBody
}

func (b *RigidBodyList) New(toEntity *Entity, invMass float64) *RigidBody {
	b.Lock()
	b.entity[toEntity] = newRidigBody(invMass)
	b.Unlock()
	return b.entity[toEntity]
}

func (b *RigidBodyList) All() map[*Entity]*RigidBody {
	result := make(map[*Entity]*RigidBody, len(b.entity))
	b.Lock()
	for k,v := range b.entity {
		result[k] = v
	}
	b.Unlock()
	return result
}

func (b *RigidBodyList) Get(fromEntity *Entity) *RigidBody {
	b.Lock()
	e := b.entity[fromEntity]
	b.Unlock()
	return e
}

func newRidigBody(invMass float64) *RigidBody {
	body := &RigidBody{
		Velocity:                  &Vector3{},
		Rotation:                  &Vector3{},
		Forces:                    &Vector3{},
		transformMatrix:           &Matrix4{},
		InverseInertiaTensor:      &Matrix3{},
		InverseInertiaTensorWorld: &Matrix3{},
		ForceAccum:                &Vector3{},
		TorqueAccum:               &Vector3{},
		MaxAcceleration:           &Vector3{},
		Acceleration:              &Vector3{},
		LinearDamping:             0.99,
		AngularDamping:            0.99,
		MaxRotation:               3.14/10,
		InvMass:                   invMass,
		CanSleep:                  true,
		isAwake:                   true,
		SleepEpsilon:              0.00001,
	}

	it := &Matrix3{}
	it.SetBlockInertiaTensor(&Vector3{1, 1, 1}, invMass)
	body.SetInertiaTensor(it)
	return body
}

type RigidBody struct {
	sync.Mutex
	// Holds the linear velocity of the rigid body in world space.
	Velocity *Vector3
	// Holds the angular velocity, or rotation for the rigid body in world space.
	Rotation *Vector3

	// Holds the inverse of the mass of the rigid body. It is more useful to hold the inverse mass
	// because integration is simpler, and because in real time simulation it is more useful to have
	// bodies with infinite mass (immovable) than zero mass (completely unstable in numerical
	// simulation).
	InvMass float64
	// Holds the inverse of the body's inertia tensor. The inertia tensor provided must not be
	// degenerate (that would mean the body had zero inertia for spinning along one axis). As long
	// as the tensor is finite, it will be invertible. The inverse tensor is used for similar
	// reasons to the use of inverse mass.
	// The inertia tensor, unlike the other variables that define a rigid body, is given in body
	// space.
	InverseInertiaTensor *Matrix3
	// Holds the amount of damping applied to linear motion.  Damping is required to remove energy
	// added through numerical instability in the integrator.
	LinearDamping float64
	// Holds the amount of damping applied to angular motion.  Damping is required to remove energy
	// added through numerical instability in the integrator.
	AngularDamping float64

	/**
	 * Derived Data
	 *
	 * These data members hold information that is derived from the other data in the class.
	 */

	// Holds the inverse inertia tensor of the body in world space. The inverse inertia tensor
	// member is specified in the body's local space. @see inverseInertiaTensor
	InverseInertiaTensorWorld *Matrix3
	// Holds the amount of motion of the body. This is a recency weighted mean that can be used to
	// put a body to sleap.
	Motion                    float64
	// A body can be put to sleep to avoid it being updated by the integration functions or affected
	// by collisions with the world.
	isAwake                   bool
	// Some bodies may never be allowed to fall asleep. User controlled bodies, for example, should
	// be always awake.
	CanSleep                  bool
	// Holds a transform matrix for converting body space into world space and vice versa. This can
	// be achieved by calling the getPointIn*Space functions.
	transformMatrix           *Matrix4

	/**
	 * Force and Torque Accumulators
	 *
	 * These data members store the current force, torque and acceleration of the rigid body. Forces
	 * can be added to the rigid body in any order, and the class decomposes them into their
	 * constituents, accumulating them for the next simulation step. At the simulation step, the
	 * accelerations are calculated and stored to be applied to the rigid body.
	 */

	// Holds the accumulated force to be applied at the next integration step.
	ForceAccum *Vector3

	// Holds the accumulated torque to be applied at the next integration step.
	TorqueAccum *Vector3

	// Holds the acceleration of the rigid body.  This value can be used to set acceleration due to
	// gravity (its primary use), or any other constant acceleration.
	Acceleration *Vector3

	MaxAcceleration *Vector3

	// limits the linear acceleration
	MaxAngularAcceleration *Vector3
	// limits the angular velocity
	MaxRotation            float64

	// Holds the linear acceleration of the rigid body, for the previous frame.
	LastFrameAcceleration *Vector3

	SleepEpsilon float64
	Forces       *Vector3
}

func (rb *RigidBody) Mass() float64 {
	return 1 / rb.InvMass
}

func (rb *RigidBody) SetInertiaTensor(inertiaTensor *Matrix3) {
	rb.InverseInertiaTensor.SetInverse(inertiaTensor)
}

func (rb *RigidBody) AddForce(force *Vector3) {
	rb.ForceAccum.Add(force)
	rb.SetAwake(true)
}

func (rb *RigidBody) AddForceAtBodyPoint(ent *Model, force, point *Vector3) {
	// convert to coordinates relative to center of mass
	pt := rb.GetPointInWorldSpace(point)
	rb.AddForceAtPoint(ent, force, pt)
	rb.SetAwake(true)
}

func (rb *RigidBody) AddForceAtPoint(body *Model, force, point *Vector3) {
	// convert to coordinates relative to center of mass
	pt := point.NewSub(body.position)
	rb.ForceAccum.Add(force)
	rb.TorqueAccum.Add(pt.NewCross(force))
	rb.SetAwake(true)
}

func (rb *RigidBody) AddTorque(torque *Vector3) {
	rb.TorqueAccum.Add(torque)
	rb.SetAwake(true)
}

func (rb *RigidBody) ClearAccumulators() {
	rb.Forces.Clear()
	rb.ForceAccum.Clear()
	rb.TorqueAccum.Clear()
}

func (rb *RigidBody) GetPointInWorldSpace(point *Vector3) *Vector3 {
	return rb.transformMatrix.TransformVector3(point)
}

func (rb *RigidBody) getTransform() *Matrix4 {
	return rb.transformMatrix
}

func (rb *RigidBody) CalculateDerivedData(body *Model) {
	body.orientation.Normalize()
	rb.calculateTransformMatrix(rb.transformMatrix, body.position, body.orientation)
	rb.transformInertiaTensor(rb.InverseInertiaTensorWorld, body.orientation, rb.InverseInertiaTensor, rb.transformMatrix)
}

func (rb *RigidBody) Awake() bool {
	rb.Lock()
	defer rb.Unlock()
	return rb.isAwake
}

func (rb *RigidBody) SetAwake(t bool) {
	rb.Lock()
	defer rb.Unlock()
	rb.isAwake = t
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

/**
 * Inline function that creates a transform matrix from a
 * position and orientation.
 */
func (rb *RigidBody) calculateTransformMatrix(transformMatrix *Matrix4, position *Vector3, orientation *Quaternion) {

	transformMatrix[0] = 1 - 2*orientation.J*orientation.J - 2*orientation.K*orientation.K
	transformMatrix[1] = 2*orientation.I*orientation.J - 2*orientation.R*orientation.K
	transformMatrix[2] = 2*orientation.I*orientation.K + 2*orientation.R*orientation.J
	transformMatrix[3] = position[0]

	transformMatrix[4] = 2*orientation.I*orientation.J + 2*orientation.R*orientation.K
	transformMatrix[5] = 1 - 2*orientation.I*orientation.I - 2*orientation.K*orientation.K
	transformMatrix[6] = 2*orientation.J*orientation.K - 2*orientation.R*orientation.I
	transformMatrix[7] = position[1]

	transformMatrix[8] = 2*orientation.I*orientation.K - 2*orientation.R*orientation.J
	transformMatrix[9] = 2*orientation.J*orientation.K + 2*orientation.R*orientation.I
	transformMatrix[10] = 1 - 2*orientation.I*orientation.I - 2*orientation.J*orientation.J
	transformMatrix[11] = position[2]
}
