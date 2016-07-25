package main

import (
	. "github.com/stojg/vivere/lib/components"
)

func NewAI(ent *Entity) *AI {
	ai := &AI{
		states: make(map[*Entity]Steering),
	}
	entity := modelList.Rand()
	ai.states[ent] = NewSeek(modelList.Get(ent), rigidList.Get(ent), entity)
	return ai
}

type AI struct {
	states map[*Entity]Steering
}

func (s *AI) Update(elapsed float64) {
	for id, ent := range s.states {
		steering := ent.GetSteering()
		body := rigidList.Get(id)
		body.AddForce(steering.linear)
		model := modelList.Get(id)
		ste := NewLookWhereYoureGoing(body, model).GetSteering()
		body.AddTorque(ste.angular)
	}
}
