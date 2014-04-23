package main

import (
	"log"
	"math/rand"
)

// List all the available actions here (bitwise position)
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

// An integer representing the Player ID
type Action uint16

type Player struct {
	id Id
	//	entity *Entity
	conn *ClientConn
}

// UserCommand represent a recieved command (Action) from the user
type UserCommand struct {
	Actions uint32
}

type Controller interface {
	GetAction(e *Entity) Action
}

type PlayerController struct {
	player *Player
}

// GetAction
func (p *PlayerController) GetAction(e *Entity) Action {

	if !p.player.conn.open {
		log.Println("Action %v", ACTION_DIE)
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

func login(conn *ClientConn) {
	p := &Player{}
	p.id = state.NextPlayerId()
	p.conn = conn

	ent := NewEntity(p.id)
	ent.model = ENTITY_BUNNY
	ent.rotation = 0.0
	ent.angularVel = 0.0
	ent.pos = NewVec(rand.Float64()*1000, rand.Float64()*600)
	ent.size = NewVec(20, 40)
	ent.controller = &PlayerController{player: p}

	state.AddPlayer(p)
	state.AddEntity(ent)
}

func disconnect(p *Player) {
	p.conn.Close()
	state.RemovePlayer(p)
}

// Check if this player have sent a command
func ActiveCommand(p *Player, action Action) bool {
	cmd := p.conn.currentCmd.Actions
	if cmd == 0 {
		return false
	}
	if cmd&(1<<action) > 0 {
		return true
	}
	return false
}

func ClearCommand(p *Player) {
	p.conn.currentCmd = UserCommand{}
}
