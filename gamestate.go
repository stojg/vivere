package main

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"io"
	"math/rand"
)

type GameState struct {
	entities *list.List
	players  []PlayerId
	tick uint64
}

var state *GameState
var stateOld *GameState

func NewGameState() *GameState {
	st := &GameState{}
	st.entities = list.New()
	st.players = make([]PlayerId, 0)
	st.tick = 0
	return st
}

func init() {
	state = NewGameState()
	stateOld = NewGameState()
	for i := 0; i < 10; i++ {
		ent := NewEntity(Id(i + 3))
		ent.model = ENTITY_BUNNY
		ent.rotation = 0.0
		ent.angularVel = (rand.Float32() - 0.5) * 12.56;
		ent.pos = NewVec(rand.Float64()*1000, rand.Float64()*600)
		ent.size = NewVec(20, 40)
		state.entities.PushBack(ent)
	}
	copyState()
}

// Copy all existing entities to the previous state
func copyState() {
	stateOld.players = state.players
	for e := state.entities.Front(); e != nil; e = e.Next() {
		e.Value.(*Entity).UpdatePrev()
	}
}

func (gs *GameState) Serialize(buf io.Writer, serAll bool) {
	bufTemp := &bytes.Buffer{}
	var updated uint16

	for e := state.entities.Front(); e != nil; e = e.Next() {
		if e.Value.(*Entity).Serialize(bufTemp, serAll) {
			updated++
		}
	}

	if updated > 0 {
		binary.Write(buf, binary.LittleEndian, updated)
		buf.Write(bufTemp.Bytes())
	}
}

func (gs *GameState) Tick() {
	gs.tick += 1;
}
