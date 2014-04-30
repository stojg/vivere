package physics

import (
	"log"
	"math"
	v "github.com/stojg/vivere/vec"
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

type EntityProvider interface {
	Entities() []Kinematic
}

type ForceRegistry struct {
	registry map[Kinematic]ForceGenerator
}

func (f *ForceRegistry) Add(e Kinematic, fg ForceGenerator) *ForceRegistry {
	if f.registry == nil {
		f.registry = make(map[Kinematic]ForceGenerator)
	}
	f.registry[e] = fg
	return f
}

func (f *ForceRegistry) Remove(e Kinematic) *ForceRegistry {
	delete(f.registry, e)
	return f
}

// @todo implement when needed
func (f *ForceRegistry) Clear() *ForceRegistry {
	return f
}

func (f *ForceRegistry) UpdateForces(duration float64) *ForceRegistry {
	for entity, forcegenerator := range f.registry {
		forcegenerator.UpdateForce(entity, duration)
	}
	return f
}

type Simulator struct {
	Forceregistry *ForceRegistry
}

func NewSimulator() *Simulator {
	s := &Simulator{}
	s.Forceregistry = &ForceRegistry{}
	return s
}

func (s *Simulator) Update(state EntityProvider, duration float64) {

	if duration == 0 {
		log.Println("Elapsed time is zero?")
		return
	}

	s.Forceregistry.UpdateForces(duration);

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
