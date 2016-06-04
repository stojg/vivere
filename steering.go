package main

import (
	"math"
	"math/rand"
)

// SteeringOutput describes wished changes in velocity (linear) and rotation (angular)
type SteeringOutput struct {
	linear  *Vector3
	angular float64
}

// Steering is the interface for all steering behaviour
type Steering interface {
	GetSteering() *SteeringOutput
}

// NewSteeringOutput returns a new zero initialized SteeringOutput
func NewSteeringOutput() *SteeringOutput {
	so := &SteeringOutput{}
	so.linear = &Vector3{}
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
	steering.angular = 0
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
	steering.angular = 0
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

// Align ensures that the character have the same orientation as the target
type Align struct {
	character    *Entity
	target       *Entity
	targetRadius float64 // 0.02
	slowRadius   float64 // 0.1
	timeToTarget float64 // 0.1
}

// GetSteering returns the angular steering to mimic the targets orientation
func (s *Align) GetSteering() *SteeringOutput {
	// Get a new steering output
	steering := NewSteeringOutput()
	// Get the naive direction to the target
	rotation := s.target.Orientation - s.character.Orientation
	// Map the result to (-pi, pi)
	rotation = s.MapToRange(rotation)
	rotationSize := math.Abs(rotation)

	// Check if we are there, return no steering
	if rotationSize < s.targetRadius {
		return steering
	}

	// We are outside the slow radius, so full rotation
	var targetRotation float64
	if rotationSize > s.slowRadius {
		targetRotation = s.character.MaxRotation
	} else {
		targetRotation = s.character.MaxRotation * rotationSize / s.slowRadius
	}

	// The final rotation combines speed (already in the variable and direction
	targetRotation *= rotation / rotationSize

	// Acceleration tries to get to the target rotation
	steering.angular = targetRotation - s.character.Rotation
	if s.timeToTarget == 0 {
		panic("timeToTarget cannot be zero")
	}
	steering.angular /= s.timeToTarget
	return steering
}

func (s *Align) MapToRange(rotation float64) float64 {
	for rotation < -math.Pi {
		rotation += math.Pi * 2
	}
	for rotation > math.Pi {
		rotation -= math.Pi * 2
	}
	return rotation
}

// Face turns the character so it 'looks' at the target
type Face struct {
	Align
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
	// Put the target together
	s.Align.target = NewEntity()
	s.Align.target.Orientation = math.Atan2(direction[0], direction[2])
	return s.Align.GetSteering()
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
	target.Orientation = math.Atan2(s.character.Velocity[0], s.character.Velocity[2])
	align := Align{}
	align.targetRadius = 0.01
	align.slowRadius = 0.04
	align.timeToTarget = 0.1
	align.character = s.character
	align.target = target
	return align.GetSteering()
}

// Wander lets the character wander around
type Wander struct {
	Face
	WanderOffset      float64 // forward offset of the wander circle
	WanderRadius      float64 // radius of the wander circle
	WanderRate        float64 // holds the max rate at which  the wander orientation can change
	WanderOrientation float64 // Holds the current orientation of the wander target
}

// NewWander returns a new Wander behaviour
func NewWander(character *Entity, offset, radius, rate float64) *Wander {
	w := &Wander{}
	w.Align.character = character
	w.WanderOffset = offset
	w.WanderRadius = radius
	w.WanderRate = rate
	w.WanderOrientation = character.Orientation
	return w
}

// GetSteering returns a new linear and angular steering for wander
func (s *Wander) GetSteering() *SteeringOutput {
	// 1. Calculate the center of the wander circle
	target := NewEntity()
	target.Position = s.character.Position.Clone()

	// Offset the character with the offset in the direction of the character orientation
	currentHeading := OrientationAsVector(s.character.Orientation)
	targetCenter := currentHeading.Scale(s.WanderOffset)
	target.Position.Add(targetCenter)

	// Update the wander orientation with a small random value
	s.WanderOrientation += s.randomBinomial() * s.WanderRate

	// From the center of the circle draw a vector in the direction of the current wanderOrientation
	offset := OrientationAsVector(s.WanderOrientation).Scale(s.WanderRadius)
	target.Position.Add(offset)

	// Delegate to face
	s.Face.target = target
	s.Face.timeToTarget = 0.1
	s.Face.targetRadius = 0.05
	s.Face.slowRadius = 0.2
	s.Face.character = s.character

	// Get the new orientation
	steering := s.Face.GetSteering()
	// Set the linear output to the current facing direction of the character
	steering.linear = OrientationAsVector(s.character.Orientation).Scale(s.character.MaxAcceleration)
	return steering
}

// randomBinomial get a random number between -1 and + 1
func (s *Wander) randomBinomial() float64 {
	return rand.Float64() - rand.Float64()
}
