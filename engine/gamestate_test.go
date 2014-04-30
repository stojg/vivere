package engine

import (
	"testing"
	"github.com/stojg/vivere/net"
)

func TestNewGameState(t *testing.T) {
	obj := &GameState{}

	if obj != obj {
		t.Error("Super fail, cant find struct")
	}
}

func TestNextIdentityID(t *testing.T) {
	state = &GameState{}
	id := state.NextEntityID()
	if id != 1 {
		t.Errorf("Wrong next Entity returned: %d", id)
	}
	id = state.NextEntityID()
	if id != 2 {
		t.Errorf("Wrong next Entity returned: %d", id)
	}
	id = state.NextEntityID()
	if id != 3 {
		t.Errorf("Wrong next Entity returned: %d", id)
	}
}

func TestAddPlayer(t *testing.T) {
	state := NewGameState()
	state.AddPlayer(&net.Player{})
}
