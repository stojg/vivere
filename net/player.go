package net

import (
	"github.com/stojg/vivere/physics"
	v "github.com/stojg/vivere/vec"
)

type Action uint16

const (
	ACTION_UP    Action = 0  // 1
	ACTION_DOWN  Action = 1  // 2
	ACTION_RIGHT Action = 2  // 4
	ACTION_LEFT  Action = 3  // 8
	ACTION_DIE   Action = 4  // 16
	ACTION_NONE  Action = 5  // 32
	ACTION_3     Action = 6  // 64
	ACTION_4     Action = 7  // 128
	ACTION_5     Action = 8  // 256
	ACTION_6     Action = 9  // 512
	ACTION_7     Action = 10 // 1024
	ACTION_8     Action = 11 // 2048
	ACTION_9     Action = 12 // 4096
	ACTION_10    Action = 13 // 8192
	ACTION_11    Action = 14 // 8192
	ACTION_12    Action = 15 // 8192
)

type Player struct {
	id uint16
	//	entity *Entity
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
