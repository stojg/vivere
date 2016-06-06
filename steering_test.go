package main

import (
	"math"
	"testing"
)

var rotationTest = []struct {
	in  float64
	out float64
}{
	{math.Pi, math.Pi},
	{-math.Pi, -math.Pi},
	{math.Pi / 2, math.Pi / 2},
	{-math.Pi / 2, -math.Pi / 2},
	{math.Pi * 2, 0},
	{-math.Pi * 2, 0},
	{-math.Pi*2 + math.Pi/2, math.Pi / 2},
}

//func TestMapToRange(t *testing.T) {
//	a := &Align{}
//	for i, tt := range rotationTest {
//		s := a.MapToRange(tt.in)
//		if s != tt.out {
//			t.Errorf("%d MapToRange(%f) => %f, want %f", i, tt.in, s, tt.out)
//		}
//	}
//}

func TestSomething(t *testing.T) {
	origin := NewEntity()
	origin.Position = &Vector3{0, 0, 0}
	origin.physics.(*RigidBody).calculateDerivedData(origin)
	target := NewEntity()
	target.Position = &Vector3{-10, 0, 0}
	target.physics.(*RigidBody).calculateDerivedData(target)
	f := &Face{
		character: origin,
		target: target,
		baseOrientation: VectorUp(),
	}
	f.Align.timeToTarget = 1

	st := f.GetSteering()

	expected := &Vector3{1,2,3}
	if !st.linear.Equals(expected) {
		t.Errorf("%v != %v", expected, st.linear)
	}


}
