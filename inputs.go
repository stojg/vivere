package main

import (
	v "github.com/stojg/vivere/vec"
	_ "log"
)

type Inputs struct{}

func (c *Inputs) Update(e *Entity, elapsed float64) {

}

type BunnyAI struct {
	physics *ParticlePhysics
}

func (ai *BunnyAI) Update(e *Entity, elapsed float64) {
	t := &v.Vec{1, 0}
	ai.physics.AddForce(t)
}

func NewBunnyAI(physics interface{}) *BunnyAI {
	b := &BunnyAI{}
	b.physics = physics.(*ParticlePhysics)
	return b
}
