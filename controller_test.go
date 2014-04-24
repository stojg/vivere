package main

import "testing"

func TestNewNPController(t *testing.T) {
	obj := NewNPController(&Perception{})
	if obj != obj {
		t.Error("Super fail, cant find struct")
	}
}
