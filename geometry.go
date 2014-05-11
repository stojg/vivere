package main

import (
	v "github.com/stojg/vivere/vec"
)

type Geometry interface {
	Collision(b Geometry) (penetration float64, normal *v.Vec)
}

type Circle struct {
	Position *v.Vec
	Radius float64
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

type Rectangle struct {
	Position *v.Vec
	Height float64
	Width float64
}
