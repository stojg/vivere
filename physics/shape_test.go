package physics

import "testing"

func TestRectangleArea(t *testing.T) {
	shape := Rectangle{H: 10, W: 10}
	const out = 100
	if x := shape.Area(); x != out {
		t.Errorf("shape.Area() = %v, want %v", x, out)
	}
}

func TestCircleArea(t *testing.T) {
	shape := Circle{R: 10}
	out := 314
	if x := int(shape.Area()); x != out {
		t.Errorf("shape.Area() = %v, want %v", x, out)
	}
}

func TestRectLeft(t *testing.T) {
	r1 := Rectangle{x: 2, y: 2, H: 4, W: 4}
	var out float64 = 0
	if x := r1.left(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}
func TestRectRight(t *testing.T) {
	r1 := Rectangle{x: 2, y: 2, H: 4, W: 4}
	var out float64 = 4
	if x := r1.right(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}

func TestRectBottom(t *testing.T) {
	r1 := Rectangle{x: 2, y: 2, H: 4, W: 4}
	var out float64 = 4
	if x := r1.bottom(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}

func TestRectTop(t *testing.T) {
	r1 := Rectangle{x: 2, y: 2, H: 4, W: 4}
	var out float64 = 0
	if x := r1.top(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}
