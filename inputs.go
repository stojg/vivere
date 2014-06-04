package main

type SimpleAI struct {
	steering Steering
}

func (ai *SimpleAI) Update(entity *Entity, elapsed float64) {
	if ai.steering == nil {
		ai.Wander(entity)
	}
	steering := ai.steering.GetSteering()
	entity.physics.(*ParticlePhysics).AddForce(steering.linear)

	if steering.angular != 0 {
		entity.physics.(*ParticlePhysics).AddRotation(steering.angular)
		return
	}
	a := LookWhereYoureGoing{}
	a.character = entity
	look := a.GetSteering()
	entity.physics.(*ParticlePhysics).AddRotation(look.angular)
}

func (ai *SimpleAI) Wander(ent *Entity) {
	ai.steering = NewWander(ent, 200, 100, 0.1)
}

func (ai *SimpleAI) Seek(ent *Entity) {
	target := NewEntity()
	target.Position = &Vector3{500, -300}
	ai.steering = NewSeek(ent, target)
}

func NewSimpleAI(physics interface{}) *SimpleAI {
	return &SimpleAI{}
}
