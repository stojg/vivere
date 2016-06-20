package main

import (
	. "github.com/stojg/vivere/lib/components"
)

func NewAI(ents []*Entity) *AI {
	ai := &AI{
		states: make(map[*Entity]Steering),
	}
	targetEid := entities.Create()
	target := modelList.New(targetEid, 1, 1, 1, ENTITY_CAMO)
	for _, e := range ents {
		ai.states[e] = NewSeek(modelList.Get(e), rigidList.Get(e), target)
	}
	return ai
}

type AI struct {
	states map[*Entity]Steering
}

func (s *AI) Update(elapsed float64) {
	for id, ent := range s.states {
		steering := ent.GetSteering()
		rigidList.Get(id).AddForce(steering.linear)
	}
}
