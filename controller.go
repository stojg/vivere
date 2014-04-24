package main

import (
	"time"
)

type Controller interface {
	GetAction(e *Entity) Action
}

type NPController struct {
	perception *Perception
	timer      time.Time
	lastAction Action
}

func NewNPController(p *Perception) *NPController {
	c := &NPController{}
	c.perception = p
	c.timer = time.Now()
	c.lastAction = ACTION_NONE
	return c
}

func (p *NPController) GetAction(e *Entity) Action {

	elapsed := time.Now().Sub(p.timer)

	if elapsed > (time.Second * 10) {
		p.timer = time.Now()
		if p.lastAction == ACTION_RIGHT {
			p.lastAction = ACTION_LEFT
		} else {
			p.lastAction = ACTION_RIGHT
		}
	}
	if p.lastAction == ACTION_RIGHT {
		e.vel[0] = 5
	} else if p.lastAction == ACTION_LEFT {
		e.vel[0] = -5
	}
	return p.lastAction
}

type PController struct {
	player *Player
}

// GetAction
func (p *PController) GetAction(e *Entity) Action {

	if !p.player.conn.open {
		return ACTION_DIE
	}

	if ActiveCommand(p.player, ACTION_UP) {
		e.vel[1] = -100
	}
	if ActiveCommand(p.player, ACTION_DOWN) {
		e.vel[1] = 100
	}
	if ActiveCommand(p.player, ACTION_LEFT) {
		e.vel[0] = -100
	}
	if ActiveCommand(p.player, ACTION_RIGHT) {
		e.vel[0] = 100
	}
	ClearCommand(p.player)
	return ACTION_NONE
}
