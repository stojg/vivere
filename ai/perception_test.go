package ai

import (
	v "github.com/stojg/vivere/vec"
	"testing"
)

func TestWorldDimension(t *testing.T) {
	obj := &Perception{}
	dims := obj.WorldDimension()
	expected := &v.Vec{1000, 600}
	if !expected.Equals(dims) {
		t.Error("Super fail, cant find struct")
	}
}
