package main

import (
	//	"log"
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

func (p *NPController) GetAction(e *Entity) *Input {

	elapsed := time.Now().Sub(p.timer)

	if elapsed > (time.Second * 10) {
		p.timer = time.Now()
		if p.lastAction == ACTION_RIGHT {
			p.lastAction = ACTION_LEFT
		} else {
			p.lastAction = ACTION_RIGHT
		}
	}

	input := &Input{}
	input.force = &Vec{0, 0}

	if p.lastAction == ACTION_RIGHT {
		e.angularVel = 1.0
		input.force.Add(&Vec{0.5, 0})
		input.action = ACTION_RIGHT
	} else if p.lastAction == ACTION_LEFT {
		e.angularVel = -1.0
		input.action = ACTION_LEFT
		input.force.Add(&Vec{-0.5, 0})
	}
	return input
}

type PController struct {
	player *Player
}

// GetAction
func (c *PController) GetAction(e *Entity) *Input {

	if !c.player.conn.open {
		return &Input{ACTION_DIE, &Vec{0, 0}, 0}
	}

	cmd := c.player.conn.currentCmd
	if cmd.Actions == 0 {
		return &Input{ACTION_NONE, &Vec{0, 0}, 0}
	}

	input := &Input{}
	input.force = &Vec{0, 0}

	// max velocity
	if cmd.Actions&(1<<ACTION_UP) > 0 {
		input.action = ACTION_UP
		input.force.Add(&Vec{0, -1})
	}

	if cmd.Actions&(1<<ACTION_DOWN) > 0 {
		input.action = ACTION_DOWN
		input.force.Add(&Vec{0, 1})
	}

	if cmd.Actions&(1<<ACTION_LEFT) > 0 {
		input.action = ACTION_LEFT
		input.force.Add(&Vec{-1, 0})
	}

	if cmd.Actions&(1<<ACTION_RIGHT) > 0 {
		input.action = ACTION_RIGHT
		input.force.Add(&Vec{1, 0})
	}
	ClearCommand(c.player)
	return input
}
