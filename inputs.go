package main

import (
	"math"
	"time"
)

type State int

const (
	STATE_IDLE State = iota
	STATE_FLEE
	STATE_HUNT
)

func NewSimpleAI(world *World) *SimpleAI {
	return &SimpleAI{
		world: world,
	}
}

type SimpleAI struct {
	world *World
	enemy *Entity
	StatefulAI
	SteeringAI
}

func (ai *SimpleAI) Update(me *Entity, elapsed float64) {

	if ai.timePassed(100 * time.Millisecond) {
		if ai.enemy != nil {
			distance := ai.enemy.Position.NewSub(me.Position).Length()
			if distance > 400 {
				ai.changeState(STATE_IDLE)
				ai.enemy = nil
				ai.steering = NewWander(me, 200, 100, 0.1)
			}
		} else {
			hunter, dist := ai.findHunter(me)
			if hunter != nil && dist < 400 {
				ai.changeState(STATE_FLEE)
				ai.enemy = hunter
				ai.steering = NewFlee(me, hunter)
			}
		}
		ai.tick()
	}

	if ai.steering == nil {
		ai.steering = NewWander(me, 200, 100, 0.1)
	}

	ai.steer(me)
}

func (ai *SimpleAI) findHunter(me *Entity) (*Entity, float64) {
	set := ai.world.entities.GetAll()
	var closest *Entity
	closestDist := math.Inf(+1)
	for _, ent := range set {
		if ent.Model != 3 {
			continue
		}
		distance := ent.Position.NewSub(me.Position).Length()
		if distance < closestDist {
			closest = ent
			closestDist = distance
		}
	}
	if closest != nil {
		return closest, closestDist
	}
	return nil, 0
}

func NewHunterAI(world *World) *HunterAI {
	return &HunterAI{
		world: world,
	}
}

type HunterAI struct {
	world  *World
	target *Entity
	StatefulAI
	SteeringAI
	Energy float64
}

func (ai *HunterAI) Update(me *Entity, elapsed float64) {

	if ai.timePassed(3 * time.Second) {

		if ai.Energy > 100 {
			pray, dist := ai.findPray(me)
			if pray != nil && dist < 350 {
				ai.changeState(STATE_HUNT)
				ai.target = pray
				ai.steering = NewSeek(me, ai.target)
				me.Model = 3
			}
		}

		if ai.Energy < 1 {
			ai.changeState(STATE_IDLE)
			me.Model = 5
			ai.steering = NewWander(me, 200, 100, 0.1)
		}

		if ai.steering == nil {
			ai.changeState(STATE_IDLE)
			me.Model = 5
			ai.steering = NewWander(me, 200, 100, 0.1)
		}

		ai.tick()
	}

	switch ai.state {
	case STATE_IDLE:
		ai.Energy += elapsed * 10
	default:
		ai.Energy -= elapsed * 20
	}

	if ai.target != nil {
		distance := ai.target.Position.NewSub(me.Position).Length()
		if distance < ai.target.Scale.NewSub(me.Scale).Length()+5 {
			ai.target.Model = 4
			ai.target = nil
			ai.stateTime = time.Time{}
			ai.Energy += 50
		}
	}

	ai.steer(me)

}

func (ai *HunterAI) findPray(me *Entity) (*Entity, float64) {
	set := ai.world.entities.GetAll()
	var closest *Entity
	closestDist := math.Inf(+1)
	for _, ent := range set {
		if ent.Model != 2 {
			continue
		}
		distance := ent.Position.NewSub(me.Position).Length()
		if distance < closestDist {
			closest = ent
			closestDist = distance
		}
	}
	if closest != nil {
		return closest, closestDist
	}
	return nil, 0
}

type StatefulAI struct {
	state     State
	stateTime time.Time
}

func (ai *StatefulAI) clear() {
	ai.stateTime = time.Time{}
	ai.state = STATE_IDLE
}

func (ai *StatefulAI) timePassed(t time.Duration) bool {
	if ai.stateTime.IsZero() {
		ai.stateTime = time.Now()
		return true
	}
	return time.Now().Sub(ai.stateTime) > t
}

func (ai *StatefulAI) changeState(s State) {
	ai.state = s
	ai.tick()
}

func (ai *StatefulAI) tick() {
	ai.stateTime = time.Now()
}

type SteeringAI struct {
	steering Steering
}

func (ai *SteeringAI) steer(me *Entity) {
	if ai.steering != nil {
		steering := ai.steering.GetSteering()
		me.physics.(*ParticlePhysics).AddForce(steering.linear)

		if steering.angular != 0 {
			me.physics.(*ParticlePhysics).AddRotation(steering.angular)
			return
		}
	}
	a := LookWhereYoureGoing{}
	a.character = me
	look := a.GetSteering()
	me.physics.(*ParticlePhysics).AddRotation(look.angular)
}
