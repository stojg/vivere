package ai

import (
	p "github.com/stojg/vivere/physics"
	v "github.com/stojg/vivere/vec"
)

type Simple struct {
	perception *Perception
	right      bool
}

func (ai *Simple) UpdateForce(e p.Kinematic, duration float64) {
	worldSize := ai.perception.WorldDimension()

	if e.Position()[0] > worldSize[0]/2 {
		ai.right = false
	} else if e.Position()[0] < worldSize[0]/2 {
		ai.right = true
	}

	force := &v.Vec{0, 0}

	if ai.right {
		force.Set(20, 0)
	} else {
		force.Set(-20, 0)
	}

	e.AddForce(force)

}
