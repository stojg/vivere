package engine

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	n "github.com/stojg/vivere/net"
	p "github.com/stojg/vivere/physics"
)

var state *GameState

type GameState struct {
	entities     []p.Kinematic
	players      []*n.Player
	tick         uint32
	nextPlayerId uint16
	nextEntityId uint16
	prevState    *GameState
	simulator    *p.Simulator
}

type Stater interface {
	Action() Action
}

type Snapshotable interface {
	UpdatePrev()
}

type Serializer interface {
	Serialize(buf io.Writer, serAll bool) bool
}

func NewGameState() *GameState {
	st := &GameState{}
	st.entities = make([]p.Kinematic, 0)
	st.players = make([]*n.Player, 0)
	st.tick = 0
	st.nextPlayerId = 0
	st.nextEntityId = 0

	st.prevState = &GameState{}
	st.prevState.entities = make([]p.Kinematic, 0)
	st.prevState.players = make([]*n.Player, 0)
	return st
}

// NextPlayerId returns the next id
func (gs *GameState) NextPlayerId() (nextPlayerId uint16) {
	gs.nextPlayerId += 1
	return gs.nextPlayerId
}

func (state *GameState) NextEntityID() (nextEntityId uint16) {
	state.nextEntityId += 1
	return state.nextEntityId
}

func (gs *GameState) AddPlayer(p *n.Player) {
	log.Printf("[+] Player %d logged in\n", p.Id())
	gs.players = append(gs.players, p)
}

func (gs *GameState) RemovePlayer(p *n.Player) {
	for index, pInList := range gs.players {
		if p.Id() != pInList.Id() {
			continue
		}
		p.Disconnect()
		// Copy the last entry to the PlayerID position
		gs.players[index] = gs.players[len(gs.players)-1]
		// Shrink the list
		gs.players = gs.players[:len(gs.players)-1]
		log.Printf("[-] Player %d was disconnected \n", p.Id())
		return
	}
}

func (gs *GameState) SetSimulator(s *p.Simulator) {
	gs.simulator = s
}

func (gs *GameState) Entities() []p.Kinematic {
	return gs.entities
}

func (gs *GameState) Players() []*n.Player {
	return gs.players
}

// Copy all existing entities to the previous state
func (gs *GameState) UpdatePrev() {
	gs.prevState.entities = gs.entities
	gs.prevState.players = gs.players
	gs.prevState.tick = gs.tick
	gs.nextPlayerId = gs.nextPlayerId

	for _, entity := range gs.entities {
		entity.(Snapshotable).UpdatePrev()
	}
}

func (gs *GameState) AddEntity (e p.Kinematic) int {
	gs.entities = append(gs.entities, e)
	log.Printf("[+] Entity added")
	return len(gs.entities)-1
}

func (gs *GameState) RemoveEntity(newEntity Unique) {
	index := -1
	for i, entity := range gs.entities {
		if entity.(Unique).Id() == newEntity.Id() {
			index = i
			break
		}
	}
	if index < 0 {
		log.Printf("[!] Couldnt remove entity with id ", newEntity.Id())
		return
	}
	// Copy the last entry to the index position
	gs.entities[index] = gs.entities[len(gs.entities)-1]
	// Shrink the list
	gs.entities = gs.entities[:len(gs.entities)-1]
	log.Printf("[-] Entity #%v removed", newEntity.Id())
}

// Serialize the game state
func (gs *GameState) Serialize(buf io.Writer, serAll bool) {
	bufTemp := &bytes.Buffer{}
	var updated uint16

	for _, entity := range gs.Entities() {
		if entity.(Serializer).Serialize(bufTemp, true) {
			updated++
			if entity.(Stater).Action() == ACTION_DIE {
				gs.RemoveEntity(entity.(Unique))
			}
		}
	}

	if updated > 0 {
		binary.Write(buf, binary.LittleEndian, gs.tick)
		binary.Write(buf, binary.LittleEndian, updated)
		buf.Write(bufTemp.Bytes())
	}
}

// Tick increases the internal game state tick counter
func (gs *GameState) Tick() uint32 {
	return gs.tick
}

func (gs *GameState) IncTick() {
	gs.tick += 1
}
