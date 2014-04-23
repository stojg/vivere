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

	ent.element = state.entities.PushBack(ent)
	p.entity = ent
	state.players = append(state.players, p)
	log.Printf("[+] Player %d logged in\n", p.id)
}

func disconnect(id Id) {
	indexPosition := -1
	for index, player := range state.players {
		if player.id == id {
			indexPosition = index
			//state.RemoveEntity(player.entity);
			break
		}
	}

	if indexPosition != -1 {
		// Copy the last entry to the PlayerID position
		state.players[indexPosition] = state.players[len(state.players)-1]
		// Shrink the list
		state.players = state.players[:len(state.players)-1]
	}
	log.Printf("[-] Player %d was disconnected \n", id)
}

// Get all the messages from the client and push the latest one to the
// clientConnection.currentCMD
func getClientInputs() {
	for _, player := range state.players {
		for {
			select {
			case cmd := <-player.conn.cmdBuf:
				player.conn.currentCmd = cmd
			default:
				goto done
			}
		}
	done:
	}
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
