package physics

import (
	v "github.com/stojg/vivere/vec"
)

type Shape interface {
	Area() float64
	Size() *v.Vec
}
