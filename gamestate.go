package main

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"io"
)

type GameState struct {
	entities     *list.List
	players      []*Player
	tick         uint32
	nextPlayerId Id
	prevState	*GameState
}

var state *GameState

func NewGameState() *GameState {
	st := &GameState{}
	st.entities = list.New()
	st.players = make([]*Player, 0)
	st.tick = 0
	st.nextPlayerId = 0
	return st
}

// NextPlayerId returns the next id
func (gs *GameState) NextPlayerId() Id {
	gs.nextPlayerId += 1
	return gs.nextPlayerId
}

func init() {
	state = NewGameState()
	state.prevState = NewGameState()
	state.UpdatePrev()
}

func(gs *GameState) AddPlayer(p *Player) {
	gs.players = append(state.players, p)
}

func(gs *GameState) RemovePlayer(p *Player) {
	for index, pInList := range state.players {
		if p.id != pInList.id {
			continue;
		}
		p.conn.Close()
		// Copy the last entry to the PlayerID position
		gs.players[index] = gs.players[len(gs.players)-1]
		// Shrink the list
		gs.players = gs.players[:len(gs.players)-1]
		return
	}
}

// Copy all existing entities to the previous state
func (gs *GameState)UpdatePrev() {
	state.prevState.entities = state.entities
	state.prevState.players = state.players
	state.prevState.tick = state.tick
	state.nextPlayerId = state.nextPlayerId
	for e := state.entities.Front(); e != nil; e = e.Next() {
		e.Value.(*Entity).UpdatePrev()
	}
}

func (gs *GameState) AddEntity(e *Entity) {
	e.element = state.entities.PushBack(e	)
}

func (gs *GameState) RemoveEntity(e *Entity) {
	state.entities.Remove(e.element);
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
