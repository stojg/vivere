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

func TestMapToRange(t *testing.T) {
	a := &Align{}
	for i, tt := range rotationTest {
		s := a.MapToRange(tt.in)
		if s != tt.out {
			t.Errorf("%d MapToRange(%f) => %f, want %f", i, tt.in, s, tt.out)
		}
	}

}
