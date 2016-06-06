package main

import (
	"fmt"
)

type State struct {
}

func NewSimpleAI(world *World) *SimpleAI {
	ai := &SimpleAI{
		world: world,
	}
	return ai
}

type SimpleAI struct {
	SteeringAI
	world  *World
	enemy  *Entity
	states []Stater
}

func (ai *SimpleAI) Update(me *Entity, elapsed float64) {
	if len(ai.states) == 0 {
		ai.states = append(ai.states, &PrayIdleState{})
		ai.states[len(ai.states)-1].Enter(me)
	}

	state := ai.states[len(ai.states)-1].handleInputs(world)
	if state != nil {
		ai.states[len(ai.states)-1] = state
		ai.states[len(ai.states)-1].Enter(me)
	}

	steer := ai.states[len(ai.states)-1].Update(elapsed)

	ai.steer(me, steer)
}

type PrayIdleState struct {
	steering Steering
	me       *Entity
}

func (state *PrayIdleState) handleInputs(w *World) Stater {
	target, distance := world.findClosest(state.me, ENTITY_HUNTER)
	if target != nil && distance < 200 {
		return &PrayFleeState{
			target: target,
		}
	}
	return nil
}

func (state *PrayIdleState) Enter(me *Entity) {
	state.me = me
	state.steering = NewWander(state.me, 200, 100, 0.1)
}

func (state *PrayIdleState) Exit(me *Entity) {
	state.me = me
}

func (state *PrayIdleState) Update(elapsed float64) Steering {
	return state.steering
}

type PrayFleeState struct {
	steering Steering
	me       *Entity
	target   *Entity
}

func (state *PrayFleeState) handleInputs(w *World) Stater {
	distance := state.target.Position.NewSub(state.me.Position).Length()
	if distance > 300 {
		return &PrayIdleState{}
	}
	return nil
}

func (state *PrayFleeState) Enter(me *Entity) {
	state.me = me
	state.steering = NewFlee(me, state.target)
}
func (state *PrayFleeState) Exit(me *Entity) {
	state.me = me
}

func (state *PrayFleeState) Update(elapsed float64) Steering {
	return state.steering
}

func NewHunterAI(world *World) *HunterAI {
	ai := &HunterAI{world: world}
	return ai
}

type HunterAI struct {
	SteeringAI
	states []Stater
	world  *World
}

func (ai *HunterAI) Update(me *Entity, elapsed float64) {
	if len(ai.states) == 0 {
		ai.states = append(ai.states, &HunterIdleState{})
		ai.states[len(ai.states)-1].Enter(me)
	}

	state := ai.states[len(ai.states)-1].handleInputs(world)
	if state != nil {
		ai.states[len(ai.states)-1] = state
		ai.states[len(ai.states)-1].Enter(me)
	}

	steer := ai.states[len(ai.states)-1].Update(elapsed)

	ai.steer(me, steer)
}

type HunterIdleState struct {
	steering Steering
	me       *Entity
	energy   float64
}

func (state *HunterIdleState) handleInputs(w *World) Stater {
	if state.energy > 200 {
		target, distance := world.findClosest(state.me, ENTITY_PRAY)
		if target != nil && distance < 350 {
			world.Log(fmt.Sprintf("will hunt for %d", target.ID))
			return &HuntState{
				energy: state.energy,
				target: target,
			}
		}
	}
	return nil
}

func (state *HunterIdleState) Enter(me *Entity) {
	world.Log(fmt.Sprintf("Enter idle state"))
	me.Type = ENTITY_CAMO
	state.me = me
	state.steering = NewWander(state.me, 200, 100, 0.1)
}

func (state *HunterIdleState) Exit(me *Entity) {
	state.me = me
}

func (state *HunterIdleState) Update(elapsed float64) Steering {
	state.energy += elapsed * 10
	return state.steering
}

type HuntState struct {
	HunterIdleState
	energy float64
	target *Entity
}

func (state *HuntState) handleInputs(w *World) Stater {
	if state.energy < 0 {
		return &HunterIdleState{
			energy: state.energy,
		}
	}

	distance := state.target.Position.NewSub(state.me.Position).Length()
	width := state.target.Scale.NewSub(state.me.Scale).Length() + 5
	if distance < width {
		world.Log(fmt.Sprintf("%d caught %d ", state.me.ID, state.target.ID))
		state.target.Type = ENTITY_SCARED
		return &HunterIdleState{
			energy: state.energy + 50,
		}
	}
	return nil
}

func (state *HuntState) Enter(me *Entity) {
	world.Log(fmt.Sprintf("Enter hunting state"))
	me.Type = ENTITY_HUNTER
	state.me = me
	state.steering = NewSeek(me, state.target)
}

func (state *HuntState) Update(elapsed float64) Steering {
	state.energy -= elapsed * 20
	return state.steering
}

type Stater interface {
	handleInputs(w *World) Stater
	Update(elapsed float64) Steering
	Enter(me *Entity)
	Exit(me *Entity)
}

type SteeringAI struct{}

func (ai *SteeringAI) steer(me *Entity, steer Steering) {
	if steer != nil {
		steering := steer.GetSteering()

		//transform := me.physics.(*RigidBody).getTransform()
		//propulsion := LocalToWorldDirn(VectorForward(), transform)
		//me.physics.(*RigidBody).AddForce(propulsion)
		me.physics.(*RigidBody).AddTorque(steering.angular)
	}
}
