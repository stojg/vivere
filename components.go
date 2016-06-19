package main

type EntityType uint16

const (
	ENTITY_NONE EntityType = iota
	ENTITY_BLOCK
	ENTITY_PRAY
	ENTITY_HUNTER
	ENTITY_SCARED
	ENTITY_CAMO
)

type Component interface{}

func NewBodyComponent(x, y, z float64, w, h, d float64) *BodyComponent {
	return &BodyComponent{
		Position:    &Vector3{x, y, z},
		Orientation: &Quaternion{1, 0, 0, 0},
		Model:       ENTITY_PRAY,
		Scale:       &Vector3{w, h, d},
	}

}

type BodyComponent struct {
	Position    *Vector3    // Holds the linear position of the rigid body in world space.
	Orientation *Quaternion // Holds the angular orientation of the rigid body in world space.
	Scale       *Vector3    // the size of this entity
	Model       EntityType
}

func NewMoveComponent(invMass float64) *MoveComponent {
	return &MoveComponent{
		Velocity:                  &Vector3{},
		Rotation:                  &Vector3{},
		forces:                    &Vector3{},
		transformMatrix:           &Matrix4{},
		InverseInertiaTensor:      &Matrix3{},
		inverseInertiaTensorWorld: &Matrix3{},
		forceAccum:                &Vector3{},
		torqueAccum:               &Vector3{},
		acceleration:              &Vector3{},
		LinearDamping:             0.99,
		AngularDamping:            0.99,
		InvMass:                   invMass,
		canSleep:                  true,
		isAwake:                   true,
		sleepEpsilon:              0.01,
	}
}

type MoveComponent struct {
	// Holds the linear velocity of the rigid body in world space.
	Velocity *Vector3
	// Holds the angular velocity, or rotation for the rigid body in world space.
	Rotation *Vector3

	// Holds the inverse of the mass of the rigid body. It
	// is more useful to hold the inverse mass because
	// integration is simpler, and because in real time
	// simulation it is more useful to have bodies with
	// infinite mass (immovable) than zero mass
	// (completely unstable in numerical simulation).
	InvMass float64
	// Holds the inverse of the body's inertia tensor. The
	// inertia tensor provided must not be degenerate
	// (that would mean the body had zero inertia for
	// spinning along one axis). As long as the tensor is
	// finite, it will be invertible. The inverse tensor
	// is used for similar reasons to the use of inverse
	// mass.
	//
	// The inertia tensor, unlike the other variables that
	// define a rigid body, is given in body space.
	InverseInertiaTensor *Matrix3
	// Holds the amount of damping applied to linear
	// motion.  Damping is required to remove energy added
	// through numerical instability in the integrator.
	LinearDamping float64
	// Holds the amount of damping applied to angular
	// motion.  Damping is required to remove energy added
	// through numerical instability in the integrator.
	AngularDamping float64

	/**
	 * Derived Data
	 *
	 * These data members hold information that is derived from
	 * the other data in the class.
	 */

	// Holds the inverse inertia tensor of the body in world
	// space. The inverse inertia tensor member is specified in
	// the body's local space.
	//  @see inverseInertiaTensor
	inverseInertiaTensorWorld *Matrix3
	// Holds the amount of motion of the body. This is a recency
	// weighted mean that can be used to put a body to sleap.
	motion float64
	// A body can be put to sleep to avoid it being updated
	// by the integration functions or affected by collisions
	// with the world.
	isAwake bool
	// Some bodies may never be allowed to fall asleep.
	// User controlled bodies, for example, should be
	// always awake.
	canSleep bool
	// Holds a transform matrix for converting body space into
	// world space and vice versa. This can be achieved by calling
	// the getPointIn*Space functions.
	transformMatrix *Matrix4

	/**
	 * Force and Torque Accumulators
	 *
	 * These data members store the current force, torque and
	 * acceleration of the rigid body. Forces can be added to the
	 * rigid body in any order, and the class decomposes them into
	 * their constituents, accumulating them for the next
	 * simulation step. At the simulation step, the accelerations
	 * are calculated and stored to be applied to the rigid body.
	 */

	// Holds the accumulated force to be applied at the next
	// integration step.
	forceAccum *Vector3

	// Holds the accumulated torque to be applied at the next
	// integration step.
	torqueAccum *Vector3

	// Holds the acceleration of the rigid body.  This value
	// can be used to set acceleration due to gravity (its primary
	// use), or any other constant acceleration.
	acceleration *Vector3

	// Holds the linear acceleration of the rigid body, for the
	// previous frame.
	lastFrameAcceleration *Vector3

	sleepEpsilon float64
	forces       *Vector3
}

func (rb *MoveComponent) Mass() float64 {
	return 1 / rb.InvMass
}

func (rb *MoveComponent) SetInertiaTensor(inertiaTensor *Matrix3) {
	rb.InverseInertiaTensor.SetInverse(inertiaTensor)
}

func (rb *MoveComponent) AddForce(force *Vector3) {
	rb.forceAccum.Add(force)
	rb.isAwake = true
}

func (rb *MoveComponent) AddForceAtBodyPoint(ent *BodyComponent, force, point *Vector3) {
	// convert to coordinates relative to center of mass
	pt := rb.getPointInWorldSpace(point)
	rb.AddForceAtPoint(ent, force, pt)
	rb.isAwake = true
}

func (rb *MoveComponent) AddForceAtPoint(body *BodyComponent, force, point *Vector3) {
	// convert to coordinates relative to center of mass
	pt := point.NewSub(body.Position)
	rb.forceAccum.Add(force)
	rb.torqueAccum.Add(pt.NewCross(force))
	rb.isAwake = true
}

func (rb *MoveComponent) AddTorque(torque *Vector3) {
	rb.torqueAccum.Add(torque)
	rb.isAwake = true
}

func (rb *MoveComponent) ClearAccumulators() {
	rb.forces.Clear()
	rb.forceAccum.Clear()
	rb.torqueAccum.Clear()
}

func (rb *MoveComponent) getPointInWorldSpace(point *Vector3) *Vector3 {
	return rb.transformMatrix.TransformVector3(point)
}

func (rb *MoveComponent) getTransform() *Matrix4 {
	return rb.transformMatrix
}

/**
 * Internal function to do an intertia tensor transform by a quaternion.
 * Note that the implementation of this function was created by an
 * automated code-generator and optimizer.
 */
func (rb *MoveComponent) transformInertiaTensor(iitWorld *Matrix3, q *Quaternion, iitBody *Matrix3, rotmat *Matrix4) {
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
func (rb *MoveComponent) calculateTransformMatrix(transformMatrix *Matrix4, position *Vector3, orientation *Quaternion) {

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
	transformMatrix[11] = position[2]
}

func (rb *MoveComponent) calculateDerivedData(body *BodyComponent) {
	body.Orientation.Normalize()
	rb.calculateTransformMatrix(rb.transformMatrix, body.Position, body.Orientation)
	rb.transformInertiaTensor(rb.inverseInertiaTensorWorld, body.Orientation, rb.InverseInertiaTensor, rb.transformMatrix)
}
