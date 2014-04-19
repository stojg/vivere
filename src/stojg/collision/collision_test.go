package collision

import "testing"

func TestRectangleArea(t *testing.T) {
	shape := Rectangle{ h: 10, w: 10 }
	const out = 100
	if x := shape.Area(); x != out {
		t.Errorf("shape.Area() = %v, want %v", x, out)
	}
}

func TestCircleArea(t *testing.T) {
	shape := Circle{ r: 10}
	out := 314
	if x := int(shape.Area()); x != out {
		t.Errorf("shape.Area() = %v, want %v", x, out)
	}
}

func TestPointInCircle(t *testing.T) {
	shape := Circle{0,0,10}
	point := Point{0,0}
	expected := true
	if result := shape.Intersect(&point); result != expected {
		t.Errorf("%v, expected %v", result, expected)
	}
}

func TestPointNotInCircle(t *testing.T) {
	shape := Circle{0,0,10}
	point := Point{11,0}
	expected := false
	if result := shape.Intersect(&point); result != expected {
		t.Errorf("%v, expected %v", result, expected)
	}
}

func TestRectLeft(t *testing.T) {
	r1 := Rectangle{ x:2,y:2, h:4, w:4 }
	var out float32  = 0
	if x := r1.left(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}
func TestRectRight(t *testing.T) {
	r1 := Rectangle{ x:2, y:2, h:4, w:4 }
	var out float32  = 4
	if x := r1.right(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}

func TestRectBottom(t *testing.T) {
	r1 := Rectangle{ x:2, y:2, h:4, w:4 }
	var out float32  = 4
	if x := r1.bottom(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}

func TestRectTop(t *testing.T) {
	r1 := Rectangle{ x:2, y:2, h:4, w:4 }
	var out float32  = 0
	if x := r1.top(); x != out {
		t.Errorf("%v, want %v", x, out)
	}
}

func TestRectAgainstRectDontIntersect(t *testing.T) {
	r1 := Rectangle{0,0,0,4}
	r2 := Rectangle{100,100,2,2}
	out := false
	if x := r1.Intersects(&r2); x != out {
		t.Errorf("r1.Intersects(r2) = %v, want %v", x, out)
	}
}

func TestRectAgainstRectDoIntersect(t *testing.T) {
	r1 := Rectangle{0,0,4,4}
	r2 := Rectangle{2,2,4,4}
	out := true
	if x := r1.Intersects(&r2); x != out {
		t.Errorf("r1.Intersects(r2) = %v, want %v", x, out)
	}
}
