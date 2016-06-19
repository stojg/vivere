package main

import (
	"math/rand"
)

type State struct {
}

func NewSimpleAI() *SimpleAI {
	return &SimpleAI{}
}

type SimpleAI struct {
	SteeringAI
	enemy *Entity
	state *RandomSearching
	stuck bool
	pos   [2]int
}

func (ai *SimpleAI) Update(me *Entity, elapsed float64) {
	if ai.stuck {
		return
	}
	if ai.state == nil {
		ai.state = &RandomSearching{}
		list := ai.getRandomPath(me)
		ai.state.Enter(me, list)
	}

	if me.Position.NewSub(world.toPosition(ai.pos)).SquareLength() < 2500 {
		ai.pos = [2]int{rand.Intn(99), rand.Intn(99)}
		ai.state = &RandomSearching{}
		list := ai.getRandomPath(me)
		ai.state.Enter(me, list)
	}

	steer := ai.state.Update(elapsed)
	ai.steer(me, steer)
}

func (ai *SimpleAI) getRandomPath(me *Entity) [][2]int {
	ai.pos = [2]int{rand.Intn(99), rand.Intn(99)}
	list, _ := PathFinder(world.graph, world.toTilePosition(me.Position), ai.pos)

	maxReRolls := 10
	for len(list) == 0 && maxReRolls > 0 {
		maxReRolls -= 1
		ai.pos = [2]int{rand.Intn(99), rand.Intn(99)}
		if world.graph.inGrid(ai.pos[0], ai.pos[1]) {
			list, _ = PathFinder(world.graph, world.toTilePosition(me.Position), ai.pos)
		}
	}

	if maxReRolls < 1 {
		Printf("entity %d could not find a new target, stuck at %v", me.ID, me.Position)
		ai.stuck = true
	}

	return list
}

type RandomSearching struct {
	steering Steering
	me       *Entity
}

func (state *RandomSearching) Update(elapsed float64) Steering {
	return state.steering
}

func (state *RandomSearching) Enter(me *Entity, list [][2]int) {
	state.me = me
	if len(list) < 1 {
		iPrintf("Entity %d is will arrive at 0,0,0", me.ID)
		target := NewEntity()
		target.Position.Set(0, 0, 0)
		state.steering = NewArrive(state.me, target)
		return
	}
	var points []*Vector3
	for i := range list {
		pos := world.toPosition(list[i])
		points = append(points, &Vector3{pos[0], 6, pos[2]})
	}
	path := &Path{points: points}
	iPrintf("Entity %d is following a %d step path to %v", me.ID, len(points), points[len(points)-1])
	state.steering = NewFollowPath(state.me, path)
}

type SteeringAI struct{}

func (ai *SteeringAI) steer(me *Entity, steer Steering) {
	if steer != nil {
		steering := steer.GetSteering()
		me.Body.AddForce(steering.linear)
		l := NewLookWhereYoureGoing(me)
		me.Body.AddTorque(l.GetSteering().angular)

	}
}
