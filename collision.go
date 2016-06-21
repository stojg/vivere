package main

import (
	"github.com/stojg/vivere/lib/components"
	. "github.com/stojg/vivere/lib/vector"
	"github.com/volkerp/goquadtree/quadtree"
	"math"
)

type collisonBody struct {
	model    *components.Model
	rigid    *components.RigidBody
	geometry interface{}
}

func (a collisonBody) BoundingBox() quadtree.BoundingBox {
	return quadtree.BoundingBox{
		MinX: a.model.Position[0] - a.model.Scale[0],
		MaxX: a.model.Position[0] + a.model.Scale[0],
		MinY: a.model.Position[1] - a.model.Scale[1],
		MaxY: a.model.Position[1] + a.model.Scale[1],
		MinZ: a.model.Position[2] - a.model.Scale[2],
		MaxZ: a.model.Position[2] + a.model.Scale[2],
	}
}

type CollisionSystem struct{}

func (s *CollisionSystem) Update(elapsed float64) {

	// @todo sort collisions in the order of the most severe
	for j := 0; j < 5; j++ {
		collisions := s.Check()
		if len(collisions) == 0 {
			return
		}
		for i := range collisions {
			collisions[i].Resolve(elapsed)
		}
	}

}

func (s *CollisionSystem) Check() []*Contact {
	collisions := make([]*Contact, 0)

	tree := quadtree.NewQuadTree(quadtree.BoundingBox{-3200 / 2, 3200 / 2, -3200 / 2, 3200 / 2, -3200 / 2, 3200 / 2})

	bodies := make([]collisonBody, 0)
	for aID, a := range collisionList.All() {
		body := collisonBody{
			geometry: a.Geometry,
			model:    modelList.Get(aID),
			rigid:    rigidList.Get(aID),
		}
		bodies = append(bodies, body)
		tree.Add(body)
	}

	for _, a := range bodies {
		if !a.rigid.IsAwake {
			continue
		}

		broadPhase := tree.Query(a.BoundingBox())
		for _, b := range broadPhase {
			if a == b {
				continue
			}

			//hashA := string(a.ID) + ":" + string(b.(*Entity).ID)
			//hashB := string(b.(*Entity).ID) + ":" + string(a.ID)
			//if checked[hashA] || checked[hashB] {
			//	continue
			//}
			//checked[hashA], checked[hashB] = true, true

			collisionPair := &Contact{
				a:           a,
				b:           b.(collisonBody),
				restitution: 0.9,
				normal:      &Vector3{},
			}

			collision, hit := s.Detect(collisionPair)
			if hit {
				collisions = append(collisions, collision)
			}
		}
	}

	return collisions

}

func (c *CollisionSystem) Detect(pair *Contact) (*Contact, bool) {

	switch pair.a.geometry.(type) {
	case *components.Circle:
		switch pair.b.geometry.(type) {
		case *components.Circle:
			CircleVsCircle(pair)
		case *components.Rectangle:
			CircleVsRectangle(pair)
		}
	case *components.Rectangle:
		switch pair.b.geometry.(type) {
		case *components.Rectangle:
			RectangleVsRectangle(pair)
		case *components.Circle:
			RectangleVsCircle(pair)
		}
	default:
		panic("unknown collision geometry")
	}
	return pair, pair.IsIntersecting
}

type Contact struct {
	a              collisonBody
	b              collisonBody
	restitution    float64
	penetration    float64
	normal         *Vector3
	IsIntersecting bool
}

func (pair *Contact) Resolve(duration float64) {
	pair.resolveVelocity(duration)
	pair.resolveInterpenetration()
}

// resolveVelocity calculates the new velocity that is the result of the collision
func (pair *Contact) resolveVelocity(duration float64) {
	// Find the velocity in the direction of the contact normal
	separatingVelocity := pair.SeparatingVelocity()

	// The objects are already separating, NOP
	if separatingVelocity > 0 {
		return
	}

	// Calculate the new separating velocity
	newSepVelocity := -separatingVelocity * pair.restitution

	// Check the velocity build up due to acceleration only
	accCausedVelocity := pair.a.rigid.Forces.Clone()
	if pair.b.rigid != nil {
		accCausedVelocity.Sub(pair.b.rigid.Forces)
	}

	// If we have closing velocity due to acceleration buildup,
	// remove it from the new separating velocity
	accCausedSepVelocity := accCausedVelocity.Dot(pair.normal) * duration
	if accCausedSepVelocity < 0 {
		newSepVelocity += pair.restitution * accCausedSepVelocity
		// make sure that we haven't removed more than was there to begin with
		if newSepVelocity < 0 {
			newSepVelocity = 0
		}
	}

	deltaVelocity := newSepVelocity - separatingVelocity

	totalInvMass := pair.a.rigid.InvMass
	if pair.b.rigid != nil {
		totalInvMass += pair.b.rigid.InvMass
	}

	// Both objects have infinite mass, so they can't actually move
	if totalInvMass == 0 {
		return
	}

	impulsePerIMass := pair.normal.NewScale(deltaVelocity / totalInvMass)

	velocityChangeA := impulsePerIMass.NewScale(pair.a.rigid.InvMass)
	pair.a.rigid.Velocity.Add(velocityChangeA)
	if pair.b.rigid != nil {
		velocityChangeB := impulsePerIMass.NewScale(-pair.b.rigid.InvMass)
		pair.b.rigid.Velocity.Add(velocityChangeB)
	}
}

func (c *Contact) SeparatingVelocity() float64 {
	relativeVel := c.a.rigid.Velocity.Clone()
	if c.b.rigid != nil {
		relativeVel.Sub(c.b.rigid.Velocity)
	}
	return relativeVel.Dot(c.normal)
}

// resolveInterpenetration separates two objects that has penetrated
func (c *Contact) resolveInterpenetration() {

	if c.penetration <= 0 {
		return
	}

	totalInvMass := c.a.rigid.InvMass
	if c.b.rigid != nil {
		totalInvMass += c.b.rigid.InvMass
	}
	// Both objects have infinite mass, so no velocity
	if totalInvMass == 0 {
		return
	}

	movePerIMass := c.normal.NewScale(c.penetration / totalInvMass)

	c.a.model.Position.Add(movePerIMass.NewScale(c.a.rigid.InvMass))
	if c.b.rigid != nil {
		c.b.model.Position.Add(movePerIMass.NewScale(-c.b.rigid.InvMass))
	}
}

func CircleVsCircle(c *Contact) {
	cA := c.a.geometry.(*components.Circle)
	cB := c.b.geometry.(*components.Circle)

	var d [3]float64
	for i := range d {
		d[i] = c.a.model.Position[i] - c.b.model.Position[i]
	}

	sqrLength := d[0]*d[0] + d[1]*d[1] + d[2]*d[2]
	if sqrLength < RealEpsilon {
		return
	}

	// Early out to avoid expensive sqrt
	if sqrLength > (cA.Radius+cB.Radius)*(cA.Radius+cB.Radius) {
		return
	}

	length := math.Sqrt(sqrLength)

	for i := range d {
		d[i] *= 1 / length
	}

	c.penetration = cA.Radius + cB.Radius - length
	c.normal = &Vector3{d[0], d[1], d[2]}
	c.IsIntersecting = true
}

func CircleVsRectangle(c *Contact) {
	c.a, c.b = c.b, c.a
	RectangleVsCircle(c)
}

func RectangleVsCircle(c *Contact) {
	rA := c.a.geometry.(*components.Rectangle)
	rA.ToWorld(c.a.model.Position)

	cB := c.b.geometry.(*components.Circle)
	c.normal = &Vector3{}

	closestPoint := &Vector3{}
	for i := 0; i < 3; i++ {
		closestPoint[i] = c.b.model.Position[i]
		if closestPoint[i] < rA.MinPoint[i] {
			closestPoint[i] = rA.MinPoint[i]
		} else if closestPoint[i] > rA.MaxPoint[i] {
			closestPoint[i] = rA.MaxPoint[i]
		}
	}

	var d [3]float64
	for i := range d {
		d[i] = closestPoint[i] - c.b.model.Position[i]
	}

	sqrLength := d[0]*d[0] + d[1]*d[1] + d[2]*d[2]

	if sqrLength < 1.0e-8 {
		return
	}

	// Early out to avoid expensive sqrt
	if sqrLength > cB.Radius*cB.Radius {
		return
	}

	length := math.Sqrt(sqrLength)
	for i := range d {
		d[i] *= 1 / length
	}

	c.penetration = length - cB.Radius
	c.normal = &Vector3{d[0], d[1], d[2]}
	c.IsIntersecting = true
}

func RectangleVsRectangle(c *Contact) {
	rA := c.a.geometry.(*components.Rectangle)
	rB := c.b.geometry.(*components.Rectangle)

	rA.ToWorld(c.a.model.Position)
	rB.ToWorld(c.b.model.Position)

	// [Minimum Translation Vector]
	mtvDistance := math.MaxFloat32 // Set current minimum distance (max float value so next value is always less)
	mtvAxis := &Vector3{}          // Axis along which to travel with the minimum distance

	// [Axes of potential separation]
	// [X Axis]
	if !testAxisSeparation(UnitX, rA.MinPoint[0], rA.MaxPoint[0], rB.MinPoint[0], rB.MaxPoint[0], mtvAxis, &mtvDistance) {
		return
	}

	// [Y Axis]
	if !testAxisSeparation(UnitY, rA.MinPoint[1], rA.MaxPoint[1], rB.MinPoint[1], rB.MaxPoint[1], mtvAxis, &mtvDistance) {
		return
	}

	// [Z Axis]
	if !testAxisSeparation(UnitZ, rA.MinPoint[2], rA.MaxPoint[2], rB.MinPoint[2], rB.MaxPoint[2], mtvAxis, &mtvDistance) {
		return
	}

	c.penetration = mtvDistance * 1.001
	c.normal = mtvAxis.Normalize()
	c.IsIntersecting = true
}

// TestAxisStatic checks if two axis overlaps and in that case calculates how much
// * Two convex shapes only overlap if they overlap on all axes of separation
// * In order to create accurate responses we need to find the
//    collision vector (Minimum Translation Vector)
// * Find if the two boxes intersect along a single axis
// * Compute the intersection interval for that axis
// * Keep the smallest intersection/penetration value
func testAxisSeparation(axis Vector3, minA, maxA, minB, maxB float64, mtvAxis *Vector3, mtvDistance *float64) bool {

	//	axisLengthSquared := axis.Dot(&axis)
	axisLengthSquared := axis[0]*axis[0] + axis[1]*axis[1] + axis[2]*axis[2]

	// If the axis is degenerate then ignore
	if axisLengthSquared < 1.0e-8 {
		return false
	}

	// Calculate the two possible overlap ranges
	// Either we overlap on the left or the right sides
	d0 := maxB - minA // 'Left' side
	d1 := maxA - minB // 'Right' side

	// Intervals do not overlap, so no intersection
	if d0 <= 0.0 || d1 <= 0.0 {
		return false
	}

	var overlap float64
	// Find out if we overlap on the 'right' or 'left' of the object.
	if d0 < d1 {
		overlap = d0
	} else {
		overlap = -d1
	}

	// The mtd vector for that axis
	var sep [3]float64
	sep[0] = axis[0] * (overlap / axisLengthSquared)
	sep[1] = axis[1] * (overlap / axisLengthSquared)
	sep[2] = axis[2] * (overlap / axisLengthSquared)

	// The mtd vector length squared
	sepLengthSquared := sep[0]*sep[0] + sep[1]*sep[1] + sep[2]*sep[2]

	// If that vector is smaller than our computed Minimum Translation
	// Distance use that vector as our current MTV distance
	if sepLengthSquared < *mtvDistance {
		*mtvDistance = math.Sqrt(sepLengthSquared)
		mtvAxis.Set(sep[0], sep[1], sep[2])
	}
	return true
}
