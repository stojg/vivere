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
	position *v.Vec
	r        float64
}

func NewCircle(pos *v.Vec, r float64) *Circle {
	return &Circle{pos, r}
}

func (c *Circle) Area() float64 {
	return c.r * c.r * math.Pi
}

func (c *Circle) Radius() float64 {
	return c.r
}

func (c *Circle) Position() *v.Vec {
	return c.position
}

/* Rectangle */
type Rectangle struct {
	position *v.Vec
	w        float64
	h        float64
}

func NewRectangle(pos *v.Vec, w, h float64) *Rectangle {
	return &Rectangle{pos, w, h}
}

func (r *Rectangle) Position() *v.Vec {
	return r.position
}

func (r *Rectangle) Width() float64 {
	return r.w
}

func (r *Rectangle) Height() float64 {
	return r.h
}

func (r *Rectangle) Area() float64 {
	return r.w * r.h
}

func (r *Rectangle) Size() *v.Vec {
	return &v.Vec{r.w, r.h}
}

func (r *Rectangle) left() float64 {
	return r.position[0] - (r.w / 2)
}

func (r *Rectangle) right() float64 {
	return r.position[0] + (r.w / 2)
}

func (r *Rectangle) bottom() float64 {
	return r.position[1] + (r.h / 2)
}

func (r *Rectangle) top() float64 {
	return r.position[1] - (r.h / 2)
}
