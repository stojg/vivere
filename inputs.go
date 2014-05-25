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
	steering          Steering
}

func (ai *BunnyAI) Update(entity *Entity, elapsed float64) {

	if ai.steering == nil {
		ai.Wander(entity)
		//ai.Seek(entity)
	}
	steering := ai.steering.GetSteering()
	entity.physics.(*ParticlePhysics).AddForce(steering.linear)

	if steering.angular == 0 {
		a := LookWhereYoureGoing{}
		a.character = entity
		look := a.GetSteering()
		entity.physics.(*ParticlePhysics).AddRotation(look.angular)
	} else {
		entity.physics.(*ParticlePhysics).AddRotation(steering.angular)
	}

}

func (ai *BunnyAI) Wander(ent *Entity) {
	ai.steering = NewWander(ent, 200, 100, 0.1)
}

func (ai *BunnyAI) Seek(ent *Entity) {
	target := NewEntity()
	target.Position = &Vector3{500, -300}
	ai.steering = NewSeek(ent, target)
}

func NewBunnyAI(physics interface{}) *BunnyAI {
	b := &BunnyAI{}
	b.physics = physics.(*ParticlePhysics)
	return b
}
