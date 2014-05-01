package physics

import (
	v "github.com/stojg/vivere/vec"
	"math"
)

type Shape interface {
	Area() float64
	Size() *v.Vec
}

/* Circle */
type Circle struct {
	X, Y, R float64
}

func (c *Circle) Area() float64 {
	return c.R * c.R * math.Pi
}

/* Rectangle */
type Rectangle struct {
	x, y, H, W float64
}

func (r *Rectangle) Area() float64 {
	return r.H * r.W
}

func (r *Rectangle) Size() *v.Vec {
	return &v.Vec{r.W, r.H}
}

func (r *Rectangle) left() float64 {
	return r.x - (r.W / 2)
}

func (r *Rectangle) right() float64 {
	return r.x + (r.W / 2)
}

func (r *Rectangle) bottom() float64 {
	return r.y + (r.H / 2)
}

func (r *Rectangle) top() float64 {
	return r.y - (r.H / 2)
}
