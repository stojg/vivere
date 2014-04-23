package main

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"io"
	"log"
)

type GameState struct {
	entities     *list.List
	players      []*Player
	tick         uint32
	nextPlayerId Id
}

var state, stateOld *GameState

func NewGameState() *GameState {
	st := &GameState{}
	st.entities = list.New()
	st.players = make([]*Player, 0)
	st.tick = 0
	st.nextPlayerId = 0
	return st
}

func (gs *GameState) NextPlayerId() Id {
	gs.nextPlayerId += 1
	return gs.nextPlayerId
}

func init() {
	state = NewGameState()
	stateOld = NewGameState()
	copyState()
}

// Copy all existing entities to the previous state
func copyState() {
	stateOld.players = state.players
	for e := state.entities.Front(); e != nil; e = e.Next() {
		e.Value.(*Entity).UpdatePrev()
	}
}

func (gs *GameState) RemoveEntity(e *Entity) {
	log.Printf("Scheduling #%v for deletion", e.id);
	e.model = ENTITY_DELETE
}

// Serialize the game state
func (gs *GameState) Serialize(buf io.Writer, serAll bool) {
	bufTemp := &bytes.Buffer{}
	var updated uint16

	for e := state.entities.Front(); e != nil; e = e.Next() {
		if e.Value.(*Entity).Serialize(bufTemp, true) {
//			if e.Value.(*Entity).model == ENTITY_DELETE {
//				log.Printf("Deleting %v", e.Value.(*Entity).id)
//				gs.entities.Remove(e.Value.(*Entity).element)
//			}
			updated++
		}
	}

	if updated > 0 {
		binary.Write(buf, binary.LittleEndian, gs.tick)
		binary.Write(buf, binary.LittleEndian, updated)
		buf.Write(bufTemp.Bytes())
	}
}

// Tick increases the internal game state tick counter
func (gs *GameState) Tick() {
	gs.tick += 1
}
