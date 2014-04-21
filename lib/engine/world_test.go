package engine

import "testing"

func TestNewWorldWidth(t *testing.T) {
	world := NewWorld(600, 500)
	const expected = 600
	if actual := world.Width; actual != expected {
		t.Errorf("Expected %v, but go %v", expected, actual)
	}
}

func TestNewWorldHeight(t *testing.T) {
	world := NewWorld(600, 500)
	const expected = 500
	if actual := world.Height; actual != expected {
		t.Errorf("Expected %v, but go %v", expected, actual)
	}
}
