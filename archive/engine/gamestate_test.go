package engine

import (
	"bytes"
	"github.com/stojg/vivere/net"
	"testing"
)

func TestNewGameState(t *testing.T) {
	obj := &GameState{}

	if obj != obj {
		t.Error("Super fail, cant find struct")
	}
}

func TestNextIdentityID(t *testing.T) {
	gamestate := &GameState{}
	id := gamestate.NextEntityID()
	if id != 1 {
		t.Errorf("Wrong next Entity returned: %d", id)
	}
	id = gamestate.NextEntityID()
	if id != 2 {
		t.Errorf("Wrong next Entity returned: %d", id)
	}
	id = gamestate.NextEntityID()
	if id != 3 {
		t.Errorf("Wrong next Entity returned: %d", id)
	}
}

func TestAddPlayer(t *testing.T) {
	gamestate := NewGameState()
	gamestate.AddPlayer(&net.Player{})
}

func TestSerialize(t *testing.T) {
	gamestate := NewGameState()

	gamestate.AddEntity(NewEntity(1))
	buf := &bytes.Buffer{}
	gamestate.Serialize(buf, true)
}
