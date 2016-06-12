package main

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
	return nil
}

func (state *PrayIdleState) Exit(me *Entity) {
	state.me = me
}

func (state *PrayIdleState) Update(elapsed float64) Steering {
	return state.steering
}

func (state *PrayIdleState) Enter(me *Entity) {
	state.me = me
	state.steering = NewWander(state.me, 100, 100, 0, 0.01)

	// target := NewEntity()
	// state.steering = NewFlee(state.me, target)
	state.steering = NewWander(state.me, 200, 50, 0, 0.05)
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
		me.Body.AddTorque(steering.angular)

		//l := NewLookWhereYoureGoing(me)
		//st := l.GetSteering()

	}
}
