package main

type Geometry interface {
	Collision(b Geometry) (penetration float64, normal *Vector3)
}

type Circle struct {
	Position *Vector3
	Radius float64
}

func (c *Circle) Collision(b Geometry) (penetration float64, normal *Vector3) {
	switch b.(type) {
	case *Circle:
		return c.VsCircle(b.(*Circle))
	default:
		panic("unknown collision geometry")
	}
	return
}

func (a *Circle) VsCircle(b *Circle) (penetration float64, normal *Vector3) {
	distanceVec := a.Position.Clone().Sub(b.Position)
	distance := distanceVec.Length()
	penetration = a.Radius + b.Radius - distance
	normal = distanceVec.Normalize()
	return penetration, normal
}

type Rectangle struct {
	Position *Vector3
	Height float64
	Width float64
}
