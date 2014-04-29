package main

import (
	"log"
	"math"
)

type Simulator struct{}

func (s *Simulator) Update(duration float64) {

	for i := 0; i < len(state.entities); i++ {

		e := state.entities[i]

		if e.InvMass() == 0 {
			continue
		}

		if duration == 0 {
			log.Println("Elapsed time is zero")
		}

		input := e.controller.GetAction(e)

		e.acceleration = *input.acceleration

		e.position.AddScaledVector(&e.velocity, duration)
		e.velocity.AddScaledVector(&e.acceleration, duration)

		e.velocity.Scale(math.Pow(e.damping, duration))
	}
}
