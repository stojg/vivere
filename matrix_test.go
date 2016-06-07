package main

import "testing"

func TestTransformVector3(t *testing.T) {

	m := &Matrix4{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	v := &Vector3{1, 2, 3}
	actual := m.TransformVector3(v)

	expected := &Vector3{18, 46, 74}

	if !actual.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

}
