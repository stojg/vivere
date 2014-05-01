package engine

import (
	"bytes"
	"encoding/binary"
	n "github.com/stojg/vivere/net"
	p "github.com/stojg/vivere/physics"
	"github.com/stojg/vivere/state"
	"io"
	"log"
)

type Snapshotable interface {
	UpdatePrev()
}

type Serializer interface {
	Serialize(buf io.Writer, serAll bool) bool
}

type GameState struct {
	entities     []p.Kinematic
	players      []*n.Player
	tick         uint32
	nextPlayerId uint16
	nextEntityId uint16
	prevState    *GameState
	simulator    *p.Simulator
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

func (gamestate *GameState) NextEntityID() (nextEntityId uint16) {
	gamestate.nextEntityId += 1
	return gamestate.nextEntityId
}

func (gs *GameState) Entities() []p.Kinematic {
	return gs.entities
}

func (gs *GameState) AddEntity(e p.Kinematic) int {
	gs.entities = append(gs.entities, e)
	log.Printf("[+] Entity #%v added ", len(gs.entities)-1)
	return len(gs.entities) - 1
}

func (gamestate *GameState) RemoveDeadEntities() {
	for i, entity := range gamestate.entities {
		if entity.(state.Stater).State() != state.DEAD {
			continue
		}
		gamestate.simulator.Remove(entity)
		gamestate.entities[i] = gamestate.entities[len(gamestate.entities)-1]
		gamestate.entities = gamestate.entities[:len(gamestate.entities)-1]
		log.Println("[-] Entity removed from gamestate")
	}
}

func (gs *GameState) NextPlayerId() (nextPlayerId uint16) {
	gs.nextPlayerId += 1
	return gs.nextPlayerId
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
		gs.players[index] = gs.players[len(gs.players)-1]
		gs.players = gs.players[:len(gs.players)-1]
		log.Printf("[-] Player %d was disconnected \n", p.Id())
		return
	}
}

func (gs *GameState) Players() []*n.Player {
	return gs.players
}

func (gamestate *GameState) Simulator() *p.Simulator {
	return gamestate.simulator
}

func (gs *GameState) SetSimulator(s *p.Simulator) {
	gs.simulator = s
}

func (gs *GameState) Tick() uint32 {
	return gs.tick
}

func (gs *GameState) IncTick() {
	gs.tick += 1
}

func (gs *GameState) UpdatePrev() {
	gs.prevState.entities = gs.entities
	gs.prevState.players = gs.players
	gs.prevState.tick = gs.tick
	gs.nextPlayerId = gs.nextPlayerId

	for _, entity := range gs.entities {
		entity.(Snapshotable).UpdatePrev()
	}
}

func (gamestate *GameState) Serialize(buf io.Writer, serAll bool) {
	bufTemp := &bytes.Buffer{}
	var updated uint16

	for _, entity := range gamestate.Entities() {
		if entity.(Serializer).Serialize(bufTemp, serAll) {
			updated++
		}
	}

	if updated > 0 {
		binary.Write(buf, binary.LittleEndian, gamestate.tick)
		binary.Write(buf, binary.LittleEndian, updated)
		buf.Write(bufTemp.Bytes())
	}
}
