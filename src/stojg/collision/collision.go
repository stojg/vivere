package collision

import (
	"math"
	// http://godoc.org/github.com/ungerik/go3d/float64/vec2
	"github.com/ungerik/go3d/vec2"
)

/* Point */
type Point struct {
	x, y float32
}

type Shape interface {
	Area() float32
}

/* Circle */
type Circle struct {
	x, y, r float32
}

func (c *Circle ) Area() float32 {
	return float32(c.r * c.r) * math.Pi
}

func (c *Circle) Intersect(p *Point) bool {
	a := vec2.T{c.x,c.y}
	b := vec2.T{p.x, p.y}
	// return a.Sub(&b).Length() < c1.r + c2.r
	return a.Sub(&b).Length() < c.r
}

/* Rectangle */
type Rectangle struct {
	x, y, h, w float32
}

func (r *Rectangle) Area() float32 {
	return r.h * r.w;
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
