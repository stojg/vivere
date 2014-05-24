package main

import (
//	"math"
//	"log"
)

type Inputs struct{}

func (c *Inputs) Update(e *Entity, elapsed float64) {

}

type BunnyAI struct {
	physics           *ParticlePhysics
	wanderOffset      float64
	wanderRadius      float64
	wanderRate        float64
	wanderOrientation float64
	maxAcceleration   float64
}

func (ai *BunnyAI) Update(entity *Entity, elapsed float64) {
	target := NewEntity()
	target.Position = &Vector3{500, -300}
	s := Seek{
		character: entity,
		target:    target,
	}
	a := LookWhereYoureGoing{}
	a.character = entity
	look := a.GetSteering()
	steering := s.GetSteering()
	entity.physics.(*ParticlePhysics).AddForce(steering.linear)
	entity.physics.(*ParticlePhysics).AddRotation(look.angular)
}

func NewBunnyAI(physics interface{}) *BunnyAI {
	b := &BunnyAI{}
	b.physics = physics.(*ParticlePhysics)
	return b
}
