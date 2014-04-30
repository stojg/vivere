package vec

import "testing"

func TestNewVec(t *testing.T) {
	obj := &Vec{}

	if obj != obj {
		t.Error("Super fail, cant find struct")
	}
}

func TestEquals(t *testing.T) {
	obj := &Vec{2, 6}
	expected := &Vec{2, 6}
	notEqual := &Vec{6, 2}

	if !obj.Equals(expected) {
		t.Errorf("Vectors should be the same")
	}

	if !expected.Equals(obj) {
		t.Errorf("Vectors should be the same")
	}

	if obj.Equals(notEqual) {
		t.Errorf("Vectors should not be the same")
	}

	if notEqual.Equals(obj) {
		t.Errorf("Vectors should not be the same")
	}

	obj2 := &Vec{0.1, 0.2}
	other2 := &Vec{0.1, 0.2}
	if !obj2.Equals(other2) {
		t.Errorf("Vectors should be the same")
	}
}

func TestNormalize(t *testing.T) {
	obj := &Vec{1, 0}
	expected := &Vec{1, 0}
	actual := obj.Normalize()
	if !actual.Equals(expected) {
		t.Errorf("Test failure %v != %v", expected, actual)
	}

	obj = &Vec{1, 1}
	expected = &Vec{0.7071067811865475, 0.7071067811865475}
	actual = obj.Normalize()
	if !actual.Equals(expected) {
		t.Errorf("Test failure %v != %v", expected, actual)
	}
}

func TestLength(t *testing.T) {
	obj := &Vec{1, 1}
	var expected float64 = 1.4142135623730951
	actual := obj.Length()
	if actual != expected {
		t.Errorf("Test failure %v != %v", expected, actual)
	}

	obj = &Vec{10, 5}
	expected = 11.180339887498949
	actual = obj.Length()
	if actual != expected {
		t.Errorf("Test failure %v != %v", expected, actual)
	}
}

func TestSquareLength(t *testing.T) {
	obj := &Vec{1, 1}
	var expected float64 = 2
	actual := obj.SquareLength()
	if actual != expected {
		t.Errorf("Test failure %v != %v", expected, actual)
	}

	obj = &Vec{10, 5}
	expected = 125
	actual = obj.SquareLength()
	if actual != expected {
		t.Errorf("Test failure %v != %v", expected, actual)
	}
}

func TestScale(t *testing.T) {
	obj := &Vec{2, 3}
	expected := &Vec{4, 6}
	if !obj.Scale(2).Equals(expected) {
		t.Errorf("Test failure %v != %v", expected, obj)
	}
}

func TestComponentProduct(t *testing.T) {
	obj := &Vec{4, 3}
	expected := &Vec{16, 9}
	if !obj.ComponentProduct(obj).Equals(expected) {
		t.Errorf("Test failure %v != %v", expected, obj)
	}
}

func TestClear(t *testing.T) {
	obj := &Vec{4, 3}
	expected := &Vec{0, 0}
	if !obj.Clear().Equals(expected) {
		t.Errorf("Test failure %v != %v", expected, obj)
	}
}
