package main

import (
	"math"
)

/* Point */
type Point struct {
	X, Y float64
}

type Shape interface {
	Area() float64
	Size() *Vec
}

/* Circle */
type Circle struct {
	X, Y, R float64
}

func (c *Circle) Area() float64 {
	return c.R * c.R * math.Pi
}

func (c *Circle) Intersect(p *Point) bool {
	a := Vec{c.X, c.Y}
	b := Vec{p.X, p.Y}
	return a.Sub(&b).Length() < c.R
}

/* Rectangle */
type Rectangle struct {
	x, y, h, w float64
}

func (r *Rectangle) Area() float64 {
	return r.h * r.w
}

func (r *Rectangle) Size() *Vec {
	return &Vec{r.w, r.h}
}

func (r *Rectangle) left() float64 {
	return r.x - (r.w / 2)
}

func (r *Rectangle) right() float64 {
	return r.x + (r.w / 2)
}

func (r *Rectangle) bottom() float64 {
	return r.y + (r.h / 2)
}

func (r *Rectangle) top() float64 {
	return r.y - (r.h / 2)
}

func (r1 *Rectangle) Intersects(r2 *Rectangle) bool {
	return (r1.left() < r2.right() && r1.right() > r2.left() && r1.top() < r2.bottom() && r1.bottom() > r2.top())
}
