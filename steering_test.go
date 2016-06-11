package main

import (
	"math"
	"testing"
)

func deg2rad(degree float64) float64 {
	return degree * (math.Pi / 180)
}

func rad2deg(radians float64) float64 {
	return radians * (180 / math.Pi)
}

func TestAlignNoRotation(t *testing.T) {
	character := NewEntity()
	target := NewEntity()

	var alignNoRotationTests = []struct {
		character *Quaternion
		target    *Quaternion
		expected  *Vector3
	}{
		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(0)), &Vector3{0, 0, 0}},
		{QuaternionFromAxisAngle(VectorY(), deg2rad(34)), QuaternionFromAxisAngle(VectorY(), deg2rad(34)), &Vector3{0, 0, 0}},
		{QuaternionFromAxisAngle(VectorY(), deg2rad(90)), QuaternionFromAxisAngle(VectorY(), deg2rad(90)), &Vector3{0, 0, 0}},
		{QuaternionFromAxisAngle(VectorY(), deg2rad(180)), QuaternionFromAxisAngle(VectorY(), deg2rad(180)), &Vector3{0, 0, 0}},
		{QuaternionFromAxisAngle(VectorY(), deg2rad(234)), QuaternionFromAxisAngle(VectorY(), deg2rad(234)), &Vector3{0, 0, 0}},
		{QuaternionFromAxisAngle(VectorY(), deg2rad(270)), QuaternionFromAxisAngle(VectorY(), deg2rad(270)), &Vector3{0, 0, 0}},
		{QuaternionFromAxisAngle(VectorY(), deg2rad(270)), QuaternionFromAxisAngle(VectorY(), deg2rad(270)), &Vector3{0, 0, 0}},
	}

	for i := range alignNoRotationTests {
		character.Orientation = alignNoRotationTests[i].character
		target.Orientation = alignNoRotationTests[i].target

		character.physics.calculateDerivedData(character)
		target.physics.calculateDerivedData(target)

		align := NewAlign(character, target, 0.5, 0.01, 0.1)
		steering := align.GetSteering()

		if !steering.angular.Equals(alignNoRotationTests[i].expected) {
			t.Errorf("Expected %v, but got %v for test %d", alignNoRotationTests[i].expected, steering.angular, i+1)
		}
	}
}

func TestAlignRotation(t *testing.T) {
	character := NewEntity()
	target := NewEntity()

	var alignTests = []struct {
		character *Quaternion
		target    *Quaternion
		expected  *Vector3
	}{
		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(45)), &Vector3{0, 15.707963267948966, 0}},
		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(90)), &Vector3{0, 15.707963267948966, 0}},
		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(179)), &Vector3{0, 15.707963267948966, 0}},
		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(180)), &Vector3{0, 15.707963267948966, 0}},
		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(181)), &Vector3{0, -15.707963267948966, 0}},
		{QuaternionFromAxisAngle(VectorY(), deg2rad(0)), QuaternionFromAxisAngle(VectorY(), deg2rad(270)), &Vector3{0, -15.707963267948966, 0}},
	}

	for i := range alignTests {
		character.Orientation = alignTests[i].character
		target.Orientation = alignTests[i].target

		character.physics.calculateDerivedData(character)
		target.physics.calculateDerivedData(target)

		align := NewAlign(character, target, 0.001, 0.002, 0.1)
		steering := align.GetSteering()

		if !steering.angular.Equals(alignTests[i].expected) {
			t.Errorf("Expected %v, but got %v for test %d", alignTests[i].expected, steering.angular, i+1)
		}
	}

}

func TestFace_calculateOrientation(t *testing.T) {
	character := NewEntity()
	target := NewEntity()

	face := NewFace(character, target)

	var alignTests = []struct {
		target   *Vector3
		expected *Quaternion
	}{
		{&Vector3{1, 0, 0}, QuaternionFromAxisAngle(VectorY(), deg2rad(0))},
		{&Vector3{-1, 0, 0}, QuaternionFromAxisAngle(VectorY(), deg2rad(180))},
		{&Vector3{0, 0, 1}, QuaternionFromAxisAngle(VectorY(), deg2rad(-90))},
		{&Vector3{0, 0, -1}, QuaternionFromAxisAngle(VectorY(), deg2rad(90))},
		{&Vector3{1, 0, 1}, QuaternionFromAxisAngle(VectorY(), deg2rad(-45))},
		{&Vector3{-1, 0, -1}, QuaternionFromAxisAngle(VectorY(), deg2rad(135))},
		{&Vector3{1, 0, -1}, QuaternionFromAxisAngle(VectorY(), deg2rad(45))},
		{&Vector3{-1, 0, 1}, QuaternionFromAxisAngle(VectorY(), deg2rad(-135))},
		{&Vector3{-0.5, 0, 0}, QuaternionFromAxisAngle(VectorY(), deg2rad(180))},
		{&Vector3{0.5, 0, 0}, QuaternionFromAxisAngle(VectorY(), deg2rad(0))},
		{&Vector3{0.5, 0, 0.5}, QuaternionFromAxisAngle(VectorY(), deg2rad(-45))},
		{&Vector3{-0.5, 0, -0.5}, QuaternionFromAxisAngle(VectorY(), deg2rad(135))},
		{&Vector3{1, 0, 0.99}, QuaternionFromAxisAngle(VectorY(), deg2rad(-44.712083933442905))},
	}

	for i, test := range alignTests {

		actual := face.calculateOrientation(test.target)
		if !actual.Equals(test.expected) {
			t.Errorf("Expected %v, but got %v for test %d", alignTests[i].expected, actual, i+1)
		}
	}
}
