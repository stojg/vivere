package main

import (
	"math"
	"math/rand"
)

type SteeringOutput struct {
	linear  *Vector3
	angular float64
}

type Steering interface {
	GetSteering() *SteeringOutput
}

func NewSteeringOutput() *SteeringOutput {
	so := &SteeringOutput{}
	so.linear = &Vector3{}
	return so
}

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

type Flee struct {
	character *Entity
	target    *Entity
}

func (s *Flee) GetSteering() *SteeringOutput {
	steering := &SteeringOutput{}
	steering.linear = s.character.Position.NewSub(s.target.Position)
	steering.linear.Normalize()
	steering.linear.Scale(s.character.MaxAcceleration)
	steering.angular = 0
	return steering
}

type Arrive struct {
	character    *Entity
	target       *Entity
	targetRadius float64
	slowRadius   float64
	timeToTarget float64
}

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

type Align struct {
	character    *Entity
	target       *Entity
	targetRadius float64 // 0.02
	slowRadius   float64 // 0.1
	timeToTarget float64 // 0.1
}

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

type Face struct {
	Align
}

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
	s.Align.target.Orientation = math.Atan2(direction[0], direction[1])

	return s.Align.GetSteering()
}

type LookWhereYoureGoing struct {
	character *Entity
}

func (s *LookWhereYoureGoing) GetSteering() *SteeringOutput {
	if s.character.Velocity.Length() == 0 {
		return NewSteeringOutput()
	}
	target := NewEntity()
	target.Orientation = math.Atan2(s.character.Velocity[0], s.character.Velocity[1])
	align := Align{}
	align.targetRadius = 0.01
	align.slowRadius = 0.04
	align.timeToTarget = 0.1
	align.character = s.character
	align.target = target
	return align.GetSteering()
}

type Wander struct {
	Face
	WanderOffset      float64 // forward offset of the wander circle
	WanderRadius      float64 // radius of the wander circle
	WanderRate        float64 // holds the max rate at which  the wander orientation can change
	WanderOrientation float64 // Holds the current orientation of the wander target
}

func NewWander(character *Entity, offset, radius, rate float64) *Wander {
	w := &Wander{}
	w.Align.character = character
	w.WanderOffset = offset
	w.WanderRadius = radius
	w.WanderRate = rate
	w.WanderOrientation = character.Orientation
	return w
}

func (s *Wander) GetSteering() *SteeringOutput {
	// Calculate the center of the wander circle
	target := NewEntity()
	target.Position = s.character.Position.Clone()

	// Offset the character with the offset in the direction of the character orientation
	currentHeading := OrientationAsVector(s.character.Orientation)

	targetCenter := currentHeading.Scale(s.WanderOffset)
	target.Position.Add(targetCenter)

	// Update the wander orientation
	s.WanderOrientation += s.randomBinomial() * s.WanderRate

	// From the center draw a vector in the direction of the current wanderOrientation
	offset := OrientationAsVector(s.WanderOrientation).Scale(s.WanderRadius)

	target.Position.Add(offset)

	s.Face.target = target
	s.Face.timeToTarget = 0.1
	s.Face.targetRadius = 0.1
	s.Face.slowRadius = 0.3
	s.Face.character = s.character

	steering := s.Face.GetSteering()

	steering.linear = OrientationAsVector(s.character.Orientation).Scale(s.character.MaxAcceleration)
	return steering
}

func (s *Wander) randomBinomial() float64 {
	return rand.Float64() - rand.Float64()
}
