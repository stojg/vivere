package physics

import (
	v "github.com/stojg/vivere/vec"
	"math"
)

type Kinematic interface {
	InvMass() float64
	Position() *v.Vec
	Velocity() *v.Vec
	AddForce(v *v.Vec)
	Forces() *v.Vec
	ClearForces()
	Damping() float64
}

type ForceKinematicMap struct {
	kinematic Kinematic
	generator ForceGenerator
}

type EntityProvider interface {
	Entities() []Kinematic
}

/**
Physics simulator
*/
type Simulator struct {
	registry []ForceKinematicMap
}

func NewSimulator() *Simulator {
	s := &Simulator{}
	s.registry = make([]ForceKinematicMap, 0)
	return s
}

func (s *Simulator) Update(state EntityProvider, duration float64) {
	if duration == 0 {
		return
	}
	s.UpdateForces(duration)

	for _, entity := range state.Entities() {
		if entity.InvMass() == 0 {
			continue
		}
		entity.Position().AddScaledVector(entity.Velocity(), duration)
		entity.Velocity().AddScaledVector(entity.Forces(), duration)
		entity.Velocity().Scale(math.Pow(entity.Damping(), duration))
		entity.ClearForces()
	}
}

func (s *Simulator) Add(e Kinematic, fg ForceGenerator) {
	if s.registry == nil {
		s.registry = make([]ForceKinematicMap, 0)
	}
	s.registry = append(s.registry, ForceKinematicMap{e, fg})
}

func (s *Simulator) Remove(e Kinematic) {
	for index, fg := range s.registry {
		if fg.kinematic == e {
			s.registry[index] = s.registry[len(s.registry)-1]
			s.registry = s.registry[:len(s.registry)-1]
		}
	}
}

func (s *Simulator) Clear() {
	s.registry = make([]ForceKinematicMap, 0)
}

func (s *Simulator) UpdateForces(duration float64) {
	for _, fg := range s.registry {
		fg.generator.UpdateForce(fg.kinematic, duration)
	}
}
