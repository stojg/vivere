package collision

import (
	v "github.com/stojg/vivere/vec"
	"math"
)

type Circle interface {
	Position() *v.Vec
	Radius() float64
}

type Rectangle interface {
	Position() *v.Vec
	Width() float64
	Height() float64
}

func CirclePoint(c Circle, point *v.Vec) bool {
	return c.Position().Sub(point).Length() <= c.Radius()
}

func RectangleRectangle(a Rectangle, b Rectangle) bool {
	return (math.Abs(a.Position()[0]-b.Position()[0])*2 <= (a.Width() + b.Width())) &&
		(math.Abs(a.Position()[1]-b.Position()[1])*2 <= (a.Height() + b.Height()))
}
