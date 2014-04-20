package collision

import (
	"math"
	"github.com/ungerik/go3d/vec2"
)

/* Point */
type Point struct {
	X, Y float32
}

type Shape interface {
	Area() float32
}

/* Circle */
type Circle struct {
	X, Y, R float32
}

func (c *Circle) Area() float32 {
	return float32(c.R*c.R) * math.Pi
}

func (c *Circle) Intersect(p *Point) bool {
	a := vec2.T{c.X, c.Y}
	b := vec2.T{p.X, p.Y}
	// return a.Sub(&b).Length() < c1.r + c2.r
	return a.Sub(&b).Length() < c.R
}

/* Rectangle */
type Rectangle struct {
	x, y, h, w float32
}

func (r *Rectangle) Area() float32 {
	return r.h * r.w
}

func (r *Rectangle) left() float32 {
	return r.x - (r.w / 2)
}

func (r *Rectangle) right() float32 {
	return r.x + (r.w / 2)
}

func (r *Rectangle) bottom() float32 {
	return r.y + (r.h / 2)
}

func (r *Rectangle) top() float32 {
	return r.y - (r.h / 2)
}

func (r1 *Rectangle) Intersects(r2 *Rectangle) bool {
	return (r1.left() < r2.right() && r1.right() > r2.left() && r1.top() < r2.bottom() && r1.bottom() > r2.top())
}
