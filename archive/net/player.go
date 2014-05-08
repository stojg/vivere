package net

import (
	"github.com/stojg/vivere/physics"
	"github.com/stojg/vivere/state"
	v "github.com/stojg/vivere/vec"
)

type Action uint16

const (
	ACTION_UP    Action = 0 // 1
	ACTION_DOWN  Action = 1 // 2
	ACTION_RIGHT Action = 2 // 4
	ACTION_LEFT  Action = 3 // 8
)

type Player struct {
	id   uint16
	conn *ClientConn
}

func NewPlayer(id uint16, c *ClientConn) *Player {
	p := &Player{}
	p.id = id
	p.conn = c
	return p
}

func (p *Player) Id() uint16 {
	return p.id
}

func (p *Player) Connected() bool {
	return p.conn.open
}

func (p *Player) Disconnect() {
	p.conn.Close()
}

func (p *Player) CurrentCommand() UserCommand {
	return p.conn.currentCmd
}

func (p *Player) ClearCommand() {
	p.conn.currentCmd = UserCommand{}
}

func (player *Player) UpdateForce(entity physics.Kinematic, duration float64) {

	defer player.ClearCommand()

	if !player.Connected() {
		entity.(state.Stater).SetState(state.DEAD)
		return
	}

	cmd := player.CurrentCommand()
	if cmd.Actions == 0 {
		return
	}

	// max velocity
	if cmd.Actions&(1<<ACTION_UP) > 0 {
		entity.AddForce(&v.Vec{0, -200})
	}

	if cmd.Actions&(1<<ACTION_DOWN) > 0 {
		entity.AddForce(&v.Vec{0, 200})
	}

	if cmd.Actions&(1<<ACTION_LEFT) > 0 {
		entity.AddForce(&v.Vec{-200, 0})
	}

	if cmd.Actions&(1<<ACTION_RIGHT) > 0 {
		entity.AddForce(&v.Vec{200, 0})
	}
}
