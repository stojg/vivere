package main

import (
	v "github.com/stojg/vivere/vec"
)
type Geometry interface {
	Collision(b Geometry) (penetration float64, normal *v.Vec)
}

type Circle struct {
	Radius float64
	Position *v.Vec
}

func (c *Circle) Collision(b Geometry) (penetration float64, normal *v.Vec) {
	switch b.(type) {
	case *Circle:
		return c.VsCircle(b.(*Circle))
	default:
		panic("unknown collision geometry")
	}
	return
}

func (a *Circle) VsCircle(b *Circle) (penetration float64, normal *v.Vec) {
	distanceVec := a.Position.NewSub(b.Position)
	distance := distanceVec.Length()
	penetration = a.Radius + b.Radius - distance
	normal = distanceVec.Normalize()
	return penetration, normal
}

