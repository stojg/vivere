package main

import (
	"math"
	"math/rand"
)

// SteeringOutput describes wished changes in velocity (linear) and rotation (angular)
type SteeringOutput struct {
	linear  *Vector3
	angular *Vector3
}

// Steering is the interface for all steering behaviour
type Steering interface {
	GetSteering() *SteeringOutput
}

// NewSteeringOutput returns a new zero initialized SteeringOutput
func NewSteeringOutput() *SteeringOutput {
	so := &SteeringOutput{}
	so.linear = &Vector3{}
	so.angular = &Vector3{}
	return so
}

// Seek makes the character to go full speed against the target
type Seek struct {
	character *Entity
	target    *Entity
}

func NewSeek(character, target *Entity) *Seek {
	s := &Seek{}
	s.character = character
	s.target = target
	return s
}

// GetSteering returns a linear steering
func (s *Seek) GetSteering() *SteeringOutput {
	steering := NewSteeringOutput()
	// Get the direction to the target
	steering.linear = s.target.Position.NewSub(s.character.Position)
	// Go full speed ahead
	steering.linear.Normalize()
	steering.linear.ComponentProduct(s.character.MaxAcceleration)
	steering.angular = &Vector3{}
	return steering
}

func NewFlee(character, target *Entity) *Flee {
	return &Flee{
		character: character,
		target:    target,
	}
}

// Flee makes the character to flee from the target
type Flee struct {
	character *Entity
	target    *Entity
}

// GetSteering returns a linear steering
func (s *Flee) GetSteering() *SteeringOutput {
	steering := &SteeringOutput{}
	steering.linear = s.character.Position.NewSub(s.target.Position)
	steering.linear.Normalize()
	steering.linear.ComponentProduct(s.character.MaxAcceleration)
	steering.angular = &Vector3{}
	return steering
}

// Arrive tries to get the character to arrive slowly at a target
type Arrive struct {
	character    *Entity
	target       *Entity
	targetRadius float64
	slowRadius   float64
	timeToTarget float64
}

// GetSteering returns a linear steering
func (s *Arrive) GetSteering() *SteeringOutput {
	// Get a new steering output
	steering := NewSteeringOutput()
	// Get the direction to the target
	direction := s.target.Position.NewSub(s.character.Position)
	distance := direction.Length()
	// We have arrived, no output
	if distance < s.targetRadius {
		return steering
	}
	// We are outside the slow radius, so full speed ahead
	var targetSpeed float64
	if distance > s.slowRadius {
		targetSpeed = s.character.MaxSpeed
	} else {
		targetSpeed = s.character.MaxSpeed * distance / s.slowRadius
	}
	// The target velocity combines speed and direction
	targetVelocity := direction
	targetVelocity.Normalize()
	targetVelocity.Scale(targetSpeed)
	// Acceleration tries to get to the target velocity
	steering.linear = targetVelocity.NewSub(s.character.Velocity)
	steering.linear.Scale(1 / s.timeToTarget)
	return steering
}

func NewAlign(c, t *Entity, slowRadius, targetRadius, timeToTarget float64) *Align {
	return &Align{
		character:    c,
		target:       t,
		targetRadius: targetRadius,
		slowRadius:   slowRadius,
		timeToTarget: timeToTarget,
	}
}

// Align ensures that the character have the same orientation as the target
type Align struct {
	character    *Entity
	target       *Entity
	targetRadius float64 // 0.02
	slowRadius   float64 // 0.1
	timeToTarget float64 // 0.1
}

// GetSteering returns the angular steering to mimic the targets orientation
func (align *Align) GetSteering() *SteeringOutput {

	steering := NewSteeringOutput()

	final := align.target.Orientation.Clone()
	invInitial := &Quaternion{
		r: align.character.Orientation.r,
		i: -align.character.Orientation.i,
		j: -align.character.Orientation.j,
		k: -align.character.Orientation.k,
	}

	q := final.Multiply(invInitial)
	// protect the ArcCos from numerical instabilities
	if q.r > 1.0 {
		q.r = 1.0
	} else if q.r < -1.0 {
		q.r = -1.0
	}

	theta := 2 * math.Acos(q.r)

	sin := 1 / (math.Sin(theta / 2))
	axis := &Vector3{
		sin * q.i,
		sin * q.j,
		sin * q.k,
	}

	theta = align.mapToRange(theta)
	thetaNoSign := math.Abs(theta)
	// Check if we are there, return no steering
	if (thetaNoSign) < align.targetRadius {
		return steering
	}

	var targetRotation float64
	if thetaNoSign > align.slowRadius {
		targetRotation = align.character.MaxRotation
	} else {
		targetRotation = align.character.MaxRotation * (thetaNoSign / align.slowRadius)
	}

	targetRotation *= theta / thetaNoSign

	axis.Normalize()
	axis.Scale(targetRotation)
	axis.Sub(align.character.Rotation)
	axis.Scale(1 / align.timeToTarget)

	steering.angular = axis
	return steering

}

func (align *Align) mapToRange(rotation float64) float64 {
	for rotation < -math.Pi {
		rotation += math.Pi * 2
	}
	for rotation > math.Pi {
		rotation -= math.Pi * 2
	}
	return rotation
}

func NewFace(character, target *Entity) *Face {
	return &Face{
		character:       character,
		target:          target,
		baseOrientation: &Quaternion{1, 0, 0, 0},
	}
}

// Face turns the character so it 'looks' at the target
type Face struct {
	character *Entity
	target    *Entity
	// @todo fix
	baseOrientation *Quaternion
}

// GetSteering returns a angular steering
func (face *Face) GetSteering() *SteeringOutput {

	// 1. Calculate the target to delegate to align

	// Work out the direction to target
	direction := face.target.Position.NewSub(face.character.Position)

	// Check for zero direction
	if direction.SquareLength() == 0 {
		return NewSteeringOutput()
	}

	target := NewEntity()
	target.Orientation = face.calculateOrientation(direction)
	align := NewAlign(face.character, target, 0.2, 0.01, 0.1)
	return align.GetSteering()
}

func (face *Face) calculateOrientation(vector *Vector3) *Quaternion {
	vector.Normalize()

	baseZVector := VectorX().Rotate(face.baseOrientation)

	if baseZVector.Equals(vector) {
		return face.baseOrientation.Clone()
	}
	if baseZVector.Equals(vector.NewInverse()) {
		// @todo need to fix this is the base orientation isn't 1,0,0,0?
		return NewQuaternion(0, 0, 1, 0)
	}

	// find the minimal rotation from the base to the target
	angle := math.Acos(baseZVector.Dot(vector))
	axis := baseZVector.Cross(vector).Normalize()

	return QuaternionFromAxisAngle(axis, angle)
}

func NewLookWhereYoureGoing(character *Entity) *LookWhereYoureGoing {
	return &LookWhereYoureGoing{
		character: character,
	}
}

// LookWhereYoureGoing turns the character so it faces the direction the character is moving
type LookWhereYoureGoing struct {
	character *Entity
}

// GetSteering returns a angular steering
func (s *LookWhereYoureGoing) GetSteering() *SteeringOutput {
	if s.character.Velocity.Length() == 0 {
		return NewSteeringOutput()
	}
	target := NewEntity()
	target.Position = s.character.Velocity.Clone().Add(s.character.Position)

	face := NewFace(s.character, target)
	return face.GetSteering()
}

// Wander lets the character wander around
type Wander struct {
	character *Entity
	// Holds the radius and offset of the wander circle. The
	// offset is now a full 3D vector
	offset         *Vector3
	WanderRadiusXZ float64
	WanderRadiusY  float64

	// holds the maximum rate at which the wander orientation
	// can change. Should be strictly less than 1/sqrt(3) = 0.577
	// to avoid the chance of ending up with a zero length wander vector
	rate float64

	// Holds the current offset of the wander target
	Vector *Vector3

	// holds the max acceleration for this character, this
	// again should be a 3D vector, typically with only a
	// non zero z value
	maxAcceleration *Vector3
}

// NewWander returns a new Wander behaviour
func NewWander(character *Entity, offset, radiusXZ, radiusY, rate float64) *Wander {
	w := &Wander{}
	w.character = character
	w.offset = &Vector3{offset, 0, 0}
	w.WanderRadiusXZ = radiusXZ
	w.WanderRadiusY = radiusY
	w.rate = rate

	w.maxAcceleration = &Vector3{1, 0, 0}
	// start by wandering straight forward
	w.Vector = &Vector3{0.92, 0, 0}
	return w
}

// GetSteering returns a new linear and angular steering for wander
func (wander *Wander) GetSteering() *SteeringOutput {

	// 1. Calculate the target to delegate to face
	wander.Vector[0] += wander.randomBinomial() * wander.rate
	//wander.Vector[1] += wander.randomBinomial() * wander.rate
	wander.Vector[2] += wander.randomBinomial() * wander.rate

	wander.Vector.Normalize()

	// 2. Calculate the transformed target direction and scale it
	target := NewEntity()
	target.Position = wander.Vector.NewRotate(wander.character.Orientation)
	target.Position[0] *= wander.WanderRadiusXZ
	target.Position[1] *= wander.WanderRadiusY
	target.Position[2] *= wander.WanderRadiusXZ

	// 3. calculate the target to send to face
	temp := wander.character.Position.Clone()
	temp.Add(wander.offset.NewRotate(wander.character.Orientation))
	target.Position.Add(temp)

	// 4. Delegate to face
	face := NewFace(wander.character, target)

	// 5. Now set the linear acceleration to be at full
	// acceleration in the direction of the orientation
	steering := face.GetSteering()
	//wander.maxAcceleration.NewRotate(wander.character.Orientation)
	steering.linear = wander.maxAcceleration.NewRotate(wander.character.Orientation)

	return steering
}

// randomBinomial get a random number between -1 and + 1
func (s *Wander) randomBinomial() float64 {
	return rand.Float64() - rand.Float64()
}
