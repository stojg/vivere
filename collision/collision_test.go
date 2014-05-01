package collision

import (
	p "github.com/stojg/vivere/physics"
	v "github.com/stojg/vivere/vec"
	"testing"
)

func TestCirclePointNo(t *testing.T) {
	circle := p.NewCircle(&v.Vec{0, 0}, 9)
	point := &v.Vec{10, 2}
	intersect := CirclePoint(circle, point)
	if intersect {
		t.Errorf("Circle shouldnt intersect point")
	}
}

func TestCirclePointYes(t *testing.T) {
	circle := p.NewCircle(&v.Vec{0, 0}, 9)
	point := &v.Vec{9, 0}
	intersect := CirclePoint(circle, point)
	if !intersect {
		t.Errorf("Circle should intersect point")
	}
}

func TestRectangleRectangleNo(t *testing.T) {
	rect1 := p.NewRectangle(&v.Vec{0, 0}, 10, 10)
	rect2 := p.NewRectangle(&v.Vec{21, 0}, 10, 10)
	intersect := RectangleRectangle(rect1, rect2)
	if intersect {
		t.Errorf("Rect1 shouldnt intersect rect2")
	}
}

func TestRectangleRectangleYes(t *testing.T) {
	rect1 := p.NewRectangle(&v.Vec{0, 0}, 10, 10)
	rect2 := p.NewRectangle(&v.Vec{10, 0}, 10, 10)
	intersect := RectangleRectangle(rect1, rect2)
	if !intersect {
		t.Errorf("Rect1 should intersect rect2")
	}
}
