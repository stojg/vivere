package main

import (
	"fmt"
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
	steering.linear.Scale(s.character.MaxAcceleration)
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
	steering.linear.Scale(s.character.MaxAcceleration)
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

func (s *Align) mapToRange(rotation float64) float64 {
	for rotation < -math.Pi {
		rotation += math.Pi * 2
	}
	for rotation > math.Pi {
		rotation -= math.Pi * 2
	}
	return rotation
}

// GetSteering returns the angular steering to mimic the targets orientation
func (align *Align) GetSteering() *SteeringOutput {
	// Get a new steering output
	steering := NewSteeringOutput()

	s := &Quaternion{
		r: align.character.Orientation.r,
		i: -align.character.Orientation.i,
		j: -align.character.Orientation.j,
		k: -align.character.Orientation.k,
	}

	q := s.Multiply(align.target.Orientation)
	q.Normalize()

	theta := 2 * math.Acos(q.r)

	// Map the result to (-pi, pi)
	angle := theta

	sin := 1 / (math.Sin(theta / 2))

	axis := &Vector3{
		sin * q.i,
		sin * q.j,
		sin * q.k,
	}

	// Check if we are there, return no steering
	if (angle) < align.targetRadius {
		return steering
	}

	var targetRotation float64
	if angle > align.slowRadius {
		targetRotation = align.character.MaxRotation
	} else {
		targetRotation = align.character.MaxRotation * (angle / align.slowRadius)
	}

	finalRotation := axis.Scale(targetRotation).Sub(align.character.Rotation)

	// apply acc to target rotation
	steering.angular = finalRotation.Scale(1 / align.timeToTarget)

	// @todo check for max acceleration?
	return steering
}

func NewFace(character, target *Entity) *Face {
	return &Face{
		character: character,
		target:    target,
	}
}

// Face turns the character so it 'looks' at the target
type Face struct {
	character *Entity
	target    *Entity
	// @todo fix
	baseOrientation *Vector3
}

// GetSteering returns a angular steering
func (s *Face) GetSteering() *SteeringOutput {

	// 1. Calculate the target to delegate to align

	// Work out the direction to target
	direction := s.target.Position.NewSub(s.character.Position)

	// Check for zero direction
	if direction.SquareLength() == 0 {
		return NewSteeringOutput()
	}

	target := NewEntity()
	target.Orientation = QuaternionToTarget(s.character.Position, s.target.Position)
	align := NewAlign(s.character, target, 0.5, 0.01, 0.1)

	return align.GetSteering()
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
	Face
	character         *Entity
	WanderOffset      float64 // forward offset of the wander circle
	WanderRadius      float64 // radius of the wander circle
	WanderRate        float64 // holds the max rate at which  the wander orientation can change
	WanderOrientation float64 // Holds the current orientation of the wander target
}

// NewWander returns a new Wander behaviour
func NewWander(character *Entity, offset, radius, rate float64) *Wander {
	w := &Wander{}
	w.character = character
	w.WanderOffset = offset
	w.WanderRadius = radius
	w.WanderRate = rate
	w.WanderOrientation = 0
	return w
}

// GetSteering returns a new linear and angular steering for wander
func (s *Wander) GetSteering() *SteeringOutput {
	// 1. Calculate the center of the wander circle
	target := NewEntity()
	target.Position = s.character.Position.Clone()

	targetCenter := s.character.physics.(*RigidBody).getPointInWorldSpace(VectorForward())
	targetCenter.Scale(s.WanderOffset)

	target.Position.Add(targetCenter)

	s.WanderOrientation += s.randomBinomial() * s.WanderRate
	offset := OrientationAsVector(s.WanderOrientation).Scale(s.WanderRadius)
	target.Position.Add(offset)

	fmt.Println(s.WanderOrientation)

	// Go full speed ahead
	target.Position.Normalize()
	//target.Position.Scale(s.character.MaxAcceleration)

	look := NewLookWhereYoureGoing(s.character)
	// Get the new orientation
	steering := look.GetSteering()

	//steering.linear = s.character.physics.(*RigidBody).getPointInWorldSpace(target.Position)
	//fmt.Println(steering.linear.Normalize())

	// Set the linear output to the current facing direction of the character
	// @todo fix for rigidbody
	//steering.linear = OrientationAsVector(s.character.Orientation).Scale(s.character.MaxAcceleration)
	return steering
}

// randomBinomial get a random number between -1 and + 1
func (s *Wander) randomBinomial() float64 {
	return rand.Float64() - rand.Float64()
}
