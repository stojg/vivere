package main

import (
	"math/rand"
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
	state *PrayIdleState
	stuck bool
	pos [2]int
}

func (ai *SimpleAI) Update(me *Entity, elapsed float64) {
	if ai.stuck {
		return
	}
	if ai.state == nil {
		ai.pos = [2]int{99,99}
		ai.state = &PrayIdleState{
			world: ai.world,
		}
		list, _ := PathFinder(ai.state.world.graph, ai.world.toTilePosition(me.Position), ai.pos)
		ai.state.Enter(me, list)
	}

	if me.Position.NewSub(ai.world.toPosition(ai.pos)).SquareLength() < 2500 {
		ai.pos = [2]int{rand.Intn(99), rand.Intn(99)}
		list, _ := PathFinder(ai.state.world.graph, ai.world.toTilePosition(me.Position), ai.pos)

		rerolls := 10
		for len(list) == 0 && rerolls > 0{
			rerolls -=1
			ai.pos = [2]int{rand.Intn(99), rand.Intn(99)}
			if world.graph.inGrid(ai.pos[0], ai.pos[1]) {
				list, _ = PathFinder(ai.state.world.graph, ai.world.toTilePosition(me.Position), ai.pos )
			}
		}

		if rerolls < 1 {
			fmt.Println("stuck!", len(list), me.Position)
			ai.stuck = true

		}
		ai.state.Enter(me, list)
	}

	steer := ai.state.Update(elapsed)
	ai.steer(me, steer)
}

type PrayIdleState struct {
	world *World
	steering Steering
	me       *Entity
}

func (state *PrayIdleState) handleInputs(w *World) Stater {
	return nil
}

func (state *PrayIdleState) Exit(me *Entity) {
	state.me = me
}

func (state *PrayIdleState) Update(elapsed float64) Steering {
	return state.steering
}

func (state *PrayIdleState) Enter(me *Entity, list [][2]int) {
	state.me = me

	// convert back to
	var points []*Vector3
	for i := range list {
		pos  := state.world.toPosition(list[i])
		points = append(points, &Vector3{pos[0], 6, pos[2]})
	}

	path := &Path{ points: points, }
	state.steering = NewFollowPath(state.me, path)

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
		me.Body.AddForce(steering.linear)
		l := NewLookWhereYoureGoing(me)
		me.Body.AddTorque(l.GetSteering().angular)
		//me.Body.AddTorque(steering.angular)

		//st := l.GetSteering()

	}
}
