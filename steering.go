package main

import (
	"github.com/stojg/vivere/lib/components"
	. "github.com/stojg/vivere/lib/vector"
	"math"
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
	return &SteeringOutput{
		linear:  &Vector3{},
		angular: &Vector3{},
	}
}

// Seek makes the character to go full speed against the target
type Seek struct {
	character *components.Model
	rigidbody *components.RigidBody
	target    *components.Model
}

func NewSeek(model *components.Model, rigid *components.RigidBody, target *components.Model) *Seek {
	s := &Seek{}
	s.character = model
	s.rigidbody = rigid
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
	steering.linear.HadamardProduct(s.rigidbody.MaxAcceleration)
	steering.angular = &Vector3{}
	return steering
}

func NewLookWhereYoureGoing(character *components.RigidBody, cbody *components.Model) *LookWhereYoureGoing {
	return &LookWhereYoureGoing{
		model: character,
		cbody: cbody,
	}
}

// LookWhereYoureGoing turns the character so it faces the direction the character is moving
type LookWhereYoureGoing struct {
	model *components.RigidBody
	cbody *components.Model
}

// GetSteering returns a angular steering
func (s *LookWhereYoureGoing) GetSteering() *SteeringOutput {
	if s.model.Velocity.Length() == 0 {
		return NewSteeringOutput()
	}
	target := &components.Model{}
	target.Position = s.model.Velocity.NewAdd(s.cbody.Position)

	face := NewFace(s.cbody, target, s.model)
	return face.GetSteering()
}

//
//func NewFlee(character, target *BodyComponent) *Flee {
//	return &Flee{
//		character: character,
//		target:    target,
//	}
//}
//
//// Flee makes the character to flee from the target
//type Flee struct {
//	character *BodyComponent
//	target *BodyComponent
//}
//
//// GetSteering returns a linear steering
//func (s *Flee) GetSteering() *SteeringOutput {
//	steering := &SteeringOutput{}
//	steering.linear = s.character.Position.NewSub(s.target.Position)
//	steering.linear.Normalize()
//	steering.linear.NewHadamardProduct(s.character.MaxAcceleration)
//	steering.angular = &Vector3{}
//	return steering
//}
//
//func NewArrive(character, target *BodyComponent) *Arrive {
//	return &Arrive{
//		character:    character,
//		target:       target,
//		targetRadius: 2,
//		slowRadius:   50,
//		timeToTarget: 0.1,
//	}
//}
//
//// Arrive tries to get the character to arrive slowly at a target
//type Arrive struct {
//	character    *BodyComponent
//	target       *BodyComponent
//	targetRadius float64
//	slowRadius   float64
//	timeToTarget float64
//}
//
//// GetSteering returns a linear steering
//func (s *Arrive) GetSteering() *SteeringOutput {
//	// Get a new steering output
//	steering := NewSteeringOutput()
//	// Get the direction to the target
//	direction := s.target.Position.NewSub(s.character.Position)
//	distance := direction.Length()
//	// We have arrived, no output
//	if distance < s.targetRadius {
//		return steering
//	}
//	// We are outside the slow radius, so full speed ahead
//	var targetSpeed float64
//	if distance > s.slowRadius {
//		targetSpeed = s.character.MaxSpeed
//	} else {
//		targetSpeed = s.character.MaxSpeed * distance / s.slowRadius
//	}
//
//	// The target velocity combines speed and direction
//	targetVelocity := direction
//	targetVelocity.Normalize()
//	targetVelocity.Scale(targetSpeed)
//	// Acceleration tries to get to the target velocity
//	steering.linear = targetVelocity.NewSub(s.character.Velocity)
//	steering.linear.Scale(1 / s.timeToTarget)
//	return steering
//}
//
func NewAlign(c, t *components.Model, cbody *components.RigidBody, slowRadius, targetRadius, timeToTarget float64) *Align {
	return &Align{
		character: c,
		cbody:     cbody,

		target: t,

		targetRadius: targetRadius,
		slowRadius:   slowRadius,
		timeToTarget: timeToTarget,
	}
}

// Align ensures that the character have the same orientation as the target
type Align struct {
	character    *components.Model
	cbody        *components.RigidBody
	target       *components.Model
	targetRadius float64 // 0.02
	slowRadius   float64 // 0.1
	timeToTarget float64 // 0.1
}

// GetSteering returns the angular steering to mimic the targets orientation
func (align *Align) GetSteering() *SteeringOutput {

	steering := NewSteeringOutput()

	invInitial := &Quaternion{
		R: align.character.Orientation.R,
		I: -align.character.Orientation.I,
		J: -align.character.Orientation.J,
		K: -align.character.Orientation.K,
	}

	q := align.target.Orientation.NewMultiply(invInitial)
	// protect the ArcCos from numerical instabilities
	if q.R > 1.0 {
		q.R = 1.0
	} else if q.R < -1.0 {
		q.R = -1.0
	}

	theta := 2 * math.Acos(q.R)

	sin := 1 / (math.Sin(theta / 2))
	axis := &Vector3{
		sin * q.I,
		sin * q.J,
		sin * q.K,
	}

	theta = align.mapToRange(theta)
	thetaNoSign := math.Abs(theta)
	// Check if we are there, return no steering
	if (thetaNoSign) < align.targetRadius {
		return steering
	}

	var targetRotation float64
	if thetaNoSign > align.slowRadius {
		targetRotation = align.cbody.MaxRotation
	} else {
		targetRotation = align.cbody.MaxRotation * (thetaNoSign / align.slowRadius)
	}

	targetRotation *= theta / thetaNoSign

	axis.Normalize()
	axis.Scale(targetRotation)
	axis.Sub(align.cbody.Rotation)
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

func NewFace(character, target *components.Model, cbody *components.RigidBody) *Face {
	return &Face{
		character:       character,
		cbody:           cbody,
		target:          target,
		baseOrientation: &Quaternion{1, 0, 0, 0},
	}
}

// Face turns the character so it 'looks' at the target
type Face struct {
	character *components.Model
	cbody     *components.RigidBody
	target    *components.Model
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

	target := &components.Model{}
	target.Orientation = face.calculateOrientation(direction)
	align := NewAlign(face.character, target, face.cbody, 0.2, 0.01, 0.1)
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
	axis := baseZVector.NewCross(vector).Normalize()

	return QuaternionFromAxisAngle(axis, angle)
}

//// Wander lets the character wander around
//type Wander struct {
//	character *BodyComponent
//	// Holds the radius and offset of the wander circle. The
//	// offset is now a full 3D vector
//	offset         *Vector3
//	WanderRadiusXZ float64
//	WanderRadiusY  float64
//
//	// holds the maximum rate at which the wander orientation
//	// can change. Should be strictly less than 1/sqrt(3) = 0.577
//	// to avoid the chance of ending up with a zero length wander vector
//	rate float64
//
//	// Holds the current offset of the wander target
//	Vector *Vector3
//
//	// holds the max acceleration for this character, this
//	// again should be a 3D vector, typically with only a
//	// non zero z value
//	maxAcceleration *Vector3
//}
//
//// NewWander returns a new Wander behaviour
//func NewWander(character *BodyComponent, offset, radiusXZ, radiusY, rate float64) *Wander {
//	w := &Wander{}
//	w.character = character
//	w.offset = &Vector3{offset, 0, 0}
//	w.WanderRadiusXZ = radiusXZ
//	w.WanderRadiusY = radiusY
//	w.rate = rate
//
//	w.maxAcceleration = &Vector3{1, 0, 0}
//	// start by wandering straight forward
//	w.Vector = &Vector3{1, 0, 0}
//
//	return w
//}
//
//// GetSteering returns a new linear and angular steering for wander
//func (wander *Wander) GetSteering() *SteeringOutput {
//
//	// 1. make a target that looks ahead
//	charOffset := wander.character.Position.NewAdd(wander.offset.NewRotate(wander.character.Orientation))
//	target := NewEntity()
//	target.Position.Add(charOffset)
//
//	// 2. randomise the wander vector a bit, this represents the "small" sphere at the center of the
//	// target
//	wander.Vector[0] += (wander.randomBinomial() * wander.rate)
//	wander.Vector[1] += (wander.randomBinomial() * wander.rate)
//	wander.Vector[2] += (wander.randomBinomial() * wander.rate)
//	wander.Vector.Normalize()
//
//	// 3. offset the target with the scaled "small" sphere
//	target.Position[0] += wander.Vector[0] * wander.WanderRadiusXZ
//	target.Position[1] += wander.Vector[1] * wander.WanderRadiusY
//	target.Position[2] += wander.Vector[2] * wander.WanderRadiusXZ
//
//	// 4. Delegate to face
//	face := NewFace(wander.character, target)
//
//	// 5. Now set the linear acceleration to be at full
//	// acceleration in the direction of the orientation
//	steering := face.GetSteering()
//
//	steering.linear = wander.maxAcceleration.NewRotate(wander.character.Orientation)
//
//	return steering
//}
//
//// randomBinomial get a random number between -1 and + 1
//func (s *Wander) randomBinomial() float64 {
//	return rand.Float64() - rand.Float64()
//}
//
//func NewFollowPath(character *BodyComponent, path *Path) *FollowPath {
//	return &FollowPath{
//		character:    character,
//		path:         path,
//		pathOffset:   1,
//		currentParam: 0,
//	}
//}
//
//type FollowPath struct {
//	character    *BodyComponent
//	path         *Path
//	pathOffset   int
//	currentParam int
//}
//
//func (follow *FollowPath) GetSteering() *SteeringOutput {
//
//	// find the current position on the path
//	follow.currentParam = follow.path.getParam(follow.character.Position, follow.currentParam)
//
//	// offset it
//	targetParam := follow.currentParam + follow.pathOffset
//
//	target := NewEntity()
//	target.Position = follow.path.getPosition(targetParam)
//
//	seek := NewSeek(follow.character, target)
//	return seek.GetSteering()
//}
//
//type Path struct {
//	points []*Vector3
//}
//
//func (p *Path) getParam(position *Vector3, lastparam int) int {
//	closest := 0
//	distance := math.MaxFloat64
//	for i := range p.points {
//		sqrDist := position.NewSub(p.points[i]).SquareLength()
//		if sqrDist < distance {
//			closest = i
//			distance = sqrDist
//		}
//	}
//	return closest
//}
//
//func (p *Path) getPosition(param int) *Vector3 {
//	if param > len(p.points)-1 {
//		param = len(p.points) - 1
//	}
//	if param < 0 {
//		param = 0
//	}
//	if len(p.points) == 0 {
//		Println("Getting a request for a Path.getPosition when Path.points is empty")
//		return &Vector3{0, 0, 0}
//	}
//
//	//fmt.Println(param, len(p.points))
//
//	return p.points[param]
//}
