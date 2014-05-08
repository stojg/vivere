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
	center := &v.Vec{500, 300}
	center.Sub(e.Position)
	center.Normalize().Scale(10)
	ai.physics.AddForce(center)
}

func NewBunnyAI(physics interface{}) *BunnyAI {
	b := &BunnyAI{}
	b.physics = physics.(*ParticlePhysics)
	return b
}
