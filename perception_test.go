package main

import "testing"

func TestNewPerception(t *testing.T) {
	obj := &Perception{}

	if obj != obj {
		t.Error("Super fail, cant find struct")
	}
}
