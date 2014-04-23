package main

import (
	"log"
	"math/rand"
)

// List all the available actions here (bitwise position)
const (
	ACTION_UP    Action = 0
	ACTION_DOWN  Action = 1
	ACTION_RIGHT Action = 2
	ACTION_LEFT  Action = 3
	ACTION_NONE  Action = 32
)

// An integer representing the Player ID
type Action uint32

type Player struct {
	id     Id
	entity *Entity
	conn   *ClientConn
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

	state.AddEntity(ent)

	p.entity = ent

	state.AddPlayer(p)

	log.Printf("[+] Player %d logged in\n", p.id)
}

func disconnect(p *Player) {
	state.RemoveEntity(p.entity)
	state.RemovePlayer(p)
	log.Printf("[-] Player %d was disconnected \n", p.id)
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
