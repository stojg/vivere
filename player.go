package main

import (
	"log"
)

// List all the available actions here
const (
	ACTION_UP   Action = 0
	ACTION_DOWN Action = 1
)

// An integer representing the Player ID
type PlayerId uint32

// An integer representation of user actions
type Action uint32

// UserCommand represent a recieved command (Action) from the user
type UserCommand struct {
	Actions uint32
}

// a list of users
//var players = make([]PlayerId, 0)

var maxId = PlayerId(0)

func newId() PlayerId {
	maxId++
	return maxId
}

func login(id PlayerId) {
	log.Printf("Player %d logged in\n", id)
	state.players = append(state.players, id)
}

func disconnect(id PlayerId) {
	indexPosition := -1
	for index, playerid := range state.players {
		if playerid == id {
			indexPosition = index
			break
		}
	}
	if indexPosition != -1 {
		// Copy the last entry to the PlayerID position
		state.players[indexPosition] = state.players[len(state.players)-1]
		// Shrink the list
		state.players = state.players[:len(state.players)-1]
	}
	log.Printf("Player %d was disconnected \n", id)
}

// Get all the messages from the client and push the latest one to the
// clientConnection.currentCMD
func getClientInputs() {
	for _, cl := range clients {
		for {
			select {
			case cmd := <-cl.cmdBuf:
				cl.currentCmd = cmd
			default:
				goto done
			}
		}
	done:
	}
}

// Check if this player have sent a command
func active(id PlayerId, action Action) bool {
	// This user has a command that is anything else than 0
	if (clients[id].currentCmd.Actions & (1 << action)) > 0 {
		return true
	}
	return false
}
