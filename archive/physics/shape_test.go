package physics

import (
	v "github.com/stojg/vivere/vec"
	"testing"
)

func TestRectangleArea(t *testing.T) {
	shape := NewRectangle(&v.Vec{0, 0}, 10, 10)
	const out = 100
	if x := shape.Area(); x != out {
		t.Errorf("shape.Area() = %v, want %v", x, out)
	}
}

func TestCircleArea(t *testing.T) {
	shape := Circle{r: 10}
	out := 314
	if x := int(shape.Area()); x != out {
		t.Errorf("shape.Area() = %v, want %v", x, out)
	}
}

func TestRectLeft(t *testing.T) {
	r1 := NewRectangle(&v.Vec{2, 2}, 4, 4)
	var out float64 = 0
	if x := r1.left(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}
func TestRectRight(t *testing.T) {
	r1 := NewRectangle(&v.Vec{2, 2}, 4, 4)
	var out float64 = 4
	if x := r1.right(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}

func TestRectBottom(t *testing.T) {
	r1 := NewRectangle(&v.Vec{2, 2}, 4, 4)
	var out float64 = 4
	if x := r1.bottom(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}

func TestRectTop(t *testing.T) {
	r1 := NewRectangle(&v.Vec{2, 2}, 4, 4)
	var out float64 = 0
	if x := r1.top(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}
