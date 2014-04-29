package main

import (
	"time"
)

type Controller interface {
	GetAction(e *Entity) *Input
}

type Input struct {
	action       Action
	acceleration *Vec
	velocity     *Vec
	rotation     float32
}

func NewInput() *Input {
	input := &Input{}
	input.action = ACTION_NONE
	input.acceleration = &Vec{0, 0}
	input.velocity = &Vec{0, 0}
	input.rotation = 0
	return input
}

type NPController struct {
	perception *Perception
	timer      time.Time
	lastAction Action
}

// All good
func NewNPController(p *Perception) *NPController {
	c := &NPController{}
	c.perception = (p)
	c.timer = time.Now()
	c.lastAction = ACTION_NONE
	return c
}

func (c *NPController) GetAction(e *Entity) *Input {
	worldSize := c.perception.WorldDimension()
	input := NewInput()
	if e.position[0] > worldSize[0]/2 {
		input.action = ACTION_LEFT
	} else if e.position[0] < worldSize[0]/2 {
		input.action = ACTION_RIGHT
	}

	if input.action == ACTION_RIGHT {
		input.acceleration[0] = 20
	} else {
		input.acceleration[0] = -20
	}
	return input
}

type PController struct {
	player *Player
}

// GetAction
func (c *PController) GetAction(e *Entity) (input *Input) {
	defer ClearCommand(c.player)

	input = NewInput()

	if !c.player.conn.open {
		input.action = ACTION_DIE
		return
	}

	cmd := c.player.conn.currentCmd
	if cmd.Actions == 0 {
		input.action = ACTION_NONE
		return
	}

	// max velocity
	if cmd.Actions&(1<<ACTION_UP) > 0 {
		input.action = ACTION_UP
		input.acceleration[1] = -200
	}

	if cmd.Actions&(1<<ACTION_DOWN) > 0 {
		input.action = ACTION_DOWN
		input.acceleration[1] = 200
	}

	if cmd.Actions&(1<<ACTION_LEFT) > 0 {
		input.action = ACTION_LEFT
		input.acceleration[0] = -200
	}

	if cmd.Actions&(1<<ACTION_RIGHT) > 0 {
		input.acceleration[0] = 200
		input.action = ACTION_RIGHT
	}
	return
}
