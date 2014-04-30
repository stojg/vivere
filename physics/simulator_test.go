package physics

import "testing"

func TestNewSimulator(t *testing.T) {
	obj := &Simulator{}

	if obj != obj {
		t.Error("Super fail, cant find struct")
	}
}
