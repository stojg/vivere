package physics

import (
	"math"
	v "github.com/stojg/vivere/vec"
)

/* Point */
type Point struct {
	X, Y float64
}

/* Circle */
type Circle struct {
	X, Y, R float64
}

func (c *Circle) Area() float64 {
	return c.R * c.R * math.Pi
}

func (c *Circle) Intersect(p *Point) bool {
	a := v.Vec{c.X, c.Y}
	b := v.Vec{p.X, p.Y}
	return a.Sub(&b).Length() < c.R
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

func (r1 *Rectangle) Intersects(r2 *Rectangle) bool {
	return (r1.left() < r2.right() && r1.right() > r2.left() && r1.top() < r2.bottom() && r1.bottom() > r2.top())
}
