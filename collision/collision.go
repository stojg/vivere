package collision

import (
	v "github.com/stojg/vivere/vec"
	"math"
)

func (c *Circle) Intersect(point *v.Vec) bool {
	a := v.Vec{c.X, c.Y}
return a.Sub(&).Length() < c.R
}

func (r1 *Rectangle) Intersects(r2 *Rectangle) bool {
	return (r1.left() < r2.right() && r1.right() > r2.left() && r1.top() < r2.bottom() && r1.bottom() > r2.top())
}
