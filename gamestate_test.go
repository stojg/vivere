package main

import "testing"

func TestNewGameState(t *testing.T) {
	obj := &GameState{}

	if obj != obj {
		t.Error("Super fail, cant find struct")
	}
}

func TestNextID(t *testing.T) {
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

func TestAddController(t *testing.T) {
	//state = &GameState{}
	//state.AddController(&NPController{})
}
