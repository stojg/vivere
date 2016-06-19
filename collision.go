package main

import (
	"github.com/volkerp/goquadtree/quadtree"
	"math"
)

type CollisionDetector struct{
	list *EntityList
	tree *quadtree.QuadTree
	entities map[uint16]*Entity
}

func (c *CollisionDetector) updateCollisionGeometry() {
	tree := quadtree.NewQuadTree(quadtree.NewBoundingBox(-world.sizeX/2, world.sizeX/2, -world.sizeY/2, world.sizeY/2))
	for _, b := range c.list.GetAll() {
		tree.Add(b)
	}
	c.tree = &tree
	c.entities = c.list.GetAll()
}

func (c *CollisionDetector) Collisions() []*Collision {
	collisions := make([]*Collision, 0)
	checked := make(map[string]bool, 0)

	for _, a := range c.entities {
		if !a.Body.isAwake {
			continue
		}

		t := c.tree.Query(a.BoundingBox())
		for _, b := range t {
			if a == b {
				continue
			}

			hashA := string(a.ID) + ":" + string(b.(*Entity).ID)
			hashB := string(b.(*Entity).ID) + ":" + string(a.ID)
			if checked[hashA] || checked[hashB] {
				continue
			}
			checked[hashA], checked[hashB] = true, true
			collision, hit := c.Detect(a, b.(*Entity))
			if hit {
				collisions = append(collisions, collision)
			}
		}
	}
	return collisions
}

func (c *CollisionDetector) isColliding(a *Entity) bool{
	checked := make(map[string]bool, 0)

	for _, b := range c.tree.Query(a.BoundingBox()) {
		if a == b {
			continue
		}

		hashA := string(a.ID) + ":" + string(b.(*Entity).ID)
		hashB := string(b.(*Entity).ID) + ":" + string(a.ID)
		if checked[hashA] || checked[hashB] {
			continue
		}
		checked[hashA], checked[hashB] = true, true
		_, hit := c.Detect(a, b.(*Entity))
		if hit {
			return true
		}
	}
	return false
}


func (c *CollisionDetector) Detect(a *Entity, b *Entity) (collision *Collision, hit bool) {

	// @todo hardcoded restitution
	collision = &Collision{
		a:           a,
		b:           b,
		restitution: 0.1,
		normal:      &Vector3{},
	}

	switch a.Geometry.(type) {
	case *Circle:
		switch b.Geometry.(type) {
		case *Circle:
			c.CircleVsCircle(collision)
		case *Rectangle:
			c.CircleVsRectangle(collision)
		}
	case *Rectangle:
		switch b.Geometry.(type) {
		case *Rectangle:
			c.RectangleVsRectangle(collision)
		case *Circle:
			c.RectangleVsCircle(collision)
		}
	default:
		panic("unknown collision geometry")
	}
	hit = collision.IsIntersecting
	return
}

func (colDec *CollisionDetector) CircleVsCircle(contact *Collision) {
	cA := contact.a.Geometry.(*Circle)
	cB := contact.b.Geometry.(*Circle)

	var c [3]float64
	for i := range c {
		c[i] = contact.a.Position[i] - contact.b.Position[i]
	}

	sqrLength := c[0]*c[0] + c[1]*c[1] + c[2]*c[2]
	if sqrLength < real_epsilon {
		return
	}

	// Early out to avoid expensive sqrt
	if sqrLength > (cA.Radius+cB.Radius)*(cA.Radius+cB.Radius) {
		return
	}

	length := math.Sqrt(sqrLength)

	for i := range c {
		c[i] *= 1 / length
	}

	contact.penetration = cA.Radius + cB.Radius - length
	contact.normal = &Vector3{c[0], c[1], c[2]}
	contact.IsIntersecting = true
}

func (c *CollisionDetector) CircleVsRectangle(collision *Collision) {
	collision.a, collision.b = collision.b, collision.a
	c.RectangleVsCircle(collision)
}

func (colDetector *CollisionDetector) RectangleVsCircle(contact *Collision) {
	rA := contact.a.Geometry.(*Rectangle)
	rA.ToWorld(contact.a.Position)

	cB := contact.b.Geometry.(*Circle)
	contact.normal = &Vector3{}

	closestPoint := &Vector3{}
	for i := 0; i < 3; i++ {
		closestPoint[i] = contact.b.Position[i]
		if closestPoint[i] < rA.MinPoint[i] {
			closestPoint[i] = rA.MinPoint[i]
		} else if closestPoint[i] > rA.MaxPoint[i] {
			closestPoint[i] = rA.MaxPoint[i]
		}
	}

	var c [3]float64
	for i := range c {
		c[i] = closestPoint[i] - contact.b.Position[i]
	}

	sqrLength := c[0]*c[0] + c[1]*c[1] + c[2]*c[2]

	if sqrLength < 1.0e-8 {
		return
	}

	// Early out to avoid expensive sqrt
	if sqrLength > cB.Radius*cB.Radius {
		return
	}

	length := math.Sqrt(sqrLength)
	for i := range c {
		c[i] *= 1 / length
	}

	contact.penetration = length - cB.Radius
	contact.normal = &Vector3{c[0], c[1], c[2]}
	contact.IsIntersecting = true
}

func (c *CollisionDetector) RectangleVsRectangle(contact *Collision) {
	rA := contact.a.Geometry.(*Rectangle)
	rB := contact.b.Geometry.(*Rectangle)

	rA.ToWorld(contact.a.Position)
	rB.ToWorld(contact.b.Position)

	// [Minimum Translation Vector]
	mtvDistance := math.MaxFloat32 // Set current minimum distance (max float value so next value is always less)
	mtvAxis := &Vector3{}          // Axis along which to travel with the minimum distance

	// [Axes of potential separation]
	// [X Axis]
	if !c.testAxisSeparation(UnitX, rA.MinPoint[0], rA.MaxPoint[0], rB.MinPoint[0], rB.MaxPoint[0], mtvAxis, &mtvDistance) {
		return
	}

	// [Y Axis]
	if !c.testAxisSeparation(UnitY, rA.MinPoint[1], rA.MaxPoint[1], rB.MinPoint[1], rB.MaxPoint[1], mtvAxis, &mtvDistance) {
		return
	}

	// [Z Axis]
	if !c.testAxisSeparation(UnitZ, rA.MinPoint[2], rA.MaxPoint[2], rB.MinPoint[2], rB.MaxPoint[2], mtvAxis, &mtvDistance) {
		return
	}

	contact.penetration = mtvDistance * 1.001
	contact.normal = mtvAxis.Normalize()
	contact.IsIntersecting = true
}

// TestAxisStatic checks if two axis overlaps and in that case calculates how much
// * Two convex shapes only overlap if they overlap on all axes of separation
// * In order to create accurate responses we need to find the
//    collision vector (Minimum Translation Vector)
// * Find if the two boxes intersect along a single axis
// * Compute the intersection interval for that axis
// * Keep the smallest intersection/penetration value
func (c *CollisionDetector) testAxisSeparation(axis Vector3, minA, maxA, minB, maxB float64, mtvAxis *Vector3, mtvDistance *float64) bool {

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

func (collision *CollisionDetector) raycast(origin, ray *Vector3) *Collision  {
	return nil
}

type Collision struct {
	a              *Entity
	b              *Entity
	restitution    float64
	penetration    float64
	normal         *Vector3
	IsIntersecting bool
}

func (collision *Collision) SeparatingVelocity() float64 {
	relativeVel := collision.a.Velocity.Clone()
	if collision.b != nil {
		relativeVel.Sub(collision.b.Velocity)
	}
	return relativeVel.Dot(collision.normal)
}

func (c *Collision) Resolve(duration float64) {
	c.resolveVelocity(duration)
	c.resolveInterpenetration()
}

// resolveInterpenetration separates two objects that has penetrated
func (c *Collision) resolveInterpenetration() {

	if c.penetration <= 0 {
		return
	}

	totalInvMass := c.a.Body.InvMass
	if c.b != nil {
		totalInvMass += c.b.Body.InvMass
	}
	// Both objects have infinite mass, so no velocity
	if totalInvMass == 0 {
		return
	}

	movePerIMass := c.normal.NewScale(c.penetration / totalInvMass)

	c.a.Position.Add(movePerIMass.NewScale(c.a.Body.InvMass))
	if c.b != nil {
		c.b.Position.Add(movePerIMass.NewScale(-c.b.Body.InvMass))
	}
}

// resolveVelocity calculates the new velocity that is the result of the collision
func (collision *Collision) resolveVelocity(duration float64) {
	// Find the velocity in the direction of the contact normal
	separatingVelocity := collision.SeparatingVelocity()

	// The objects are already separating, NOP
	if separatingVelocity > 0 {
		return
	}

	// Calculate the new separating velocity
	newSepVelocity := -separatingVelocity * collision.restitution

	// Check the velocity build up due to acceleration only
	accCausedVelocity := collision.a.Body.forces.Clone()
	if collision.b != nil {
		accCausedVelocity.Sub(collision.b.Body.forces)
	}

	// If we have closing velocity due to acceleration buildup,
	// remove it from the new separating velocity
	accCausedSepVelocity := accCausedVelocity.Dot(collision.normal) * duration
	if accCausedSepVelocity < 0 {
		newSepVelocity += collision.restitution * accCausedSepVelocity
		// make sure that we haven't removed more than was there to begin with
		if newSepVelocity < 0 {
			newSepVelocity = 0
		}
	}

	deltaVelocity := newSepVelocity - separatingVelocity

	totalInvMass := collision.a.Body.InvMass
	if collision.b != nil {
		totalInvMass += collision.b.Body.InvMass
	}

	// Both objects have infinite mass, so they can't actually move
	if totalInvMass == 0 {
		return
	}

	impulsePerIMass := collision.normal.NewScale(deltaVelocity / totalInvMass)

	velocityChangeA := impulsePerIMass.NewScale(collision.a.Body.InvMass)
	collision.a.Velocity.Add(velocityChangeA)
	if collision.b != nil {
		velocityChangeB := impulsePerIMass.NewScale(-collision.b.Body.InvMass)
		collision.b.Velocity.Add(velocityChangeB)
	}
}
