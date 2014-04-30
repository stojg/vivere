package ai

import (
	v "github.com/stojg/vivere/vec"
)
type Perception struct{}

func (p *Perception) WorldDimension() *v.Vec {
	return &v.Vec{1000, 600}
}
