package main

import "testing"

func TestTransformVector3(t *testing.T) {

	m := &Matrix4{0, 1, 2, 3}
	v := &Vector3{1, 2, 3}
	actual := m.TransformVector3(v)

	expected := &Vector3{10, 3, 3}

	if !actual.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

}
