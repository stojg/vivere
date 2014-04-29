package main

import (
	"time"
)

type Controller interface {
	GetAction(e *Entity) *Input
}

type Input struct {
	action   Action
	force    *Vec
	rotation float32
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

	input := &Input{}
	input.force = &Vec{0, 0}
	input.rotation = 0

	if e.tx.position[0] > worldSize[0]/2 {
		input.action = ACTION_LEFT
	} else if e.tx.position[0] < worldSize[0]/2 {
		input.action = ACTION_RIGHT
	}

	if input.action == ACTION_RIGHT {
		input.force[0] = 0.01
	} else {
		input.force[0] = -0.01
	}
	return input
}

type PController struct {
	player *Player
}

// GetAction
func (c *PController) GetAction(e *Entity) (input *Input) {
	defer ClearCommand(c.player)

	input = &Input{}
	input.force = &Vec{0, 0}

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
		input.force[1] = -0.1
	}

	if cmd.Actions&(1<<ACTION_DOWN) > 0 {
		input.action = ACTION_DOWN
		input.force[1] = 0.1
	}

	if cmd.Actions&(1<<ACTION_LEFT) > 0 {
		input.action = ACTION_LEFT
		input.force[0] = -0.1
	}

	if cmd.Actions&(1<<ACTION_RIGHT) > 0 {
		input.force[0] = 0.1
		input.action = ACTION_RIGHT
	}
	return
}
