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

var alignNoRotationTests = []struct {
	character *Quaternion
	target    *Quaternion
	expected  *Vector3
}{
	{QuaternionFromAngle(VectorUp(), deg2rad(0)), QuaternionFromAngle(VectorUp(), deg2rad(0)), &Vector3{0, 0, 0}},
	{QuaternionFromAngle(VectorUp(), deg2rad(34)), QuaternionFromAngle(VectorUp(), deg2rad(34)), &Vector3{0, 0, 0}},
	{QuaternionFromAngle(VectorUp(), deg2rad(90)), QuaternionFromAngle(VectorUp(), deg2rad(90)), &Vector3{0, 0, 0}},
	{QuaternionFromAngle(VectorUp(), deg2rad(180)), QuaternionFromAngle(VectorUp(), deg2rad(180)), &Vector3{0, 0, 0}},
	{QuaternionFromAngle(VectorUp(), deg2rad(234)), QuaternionFromAngle(VectorUp(), deg2rad(234)), &Vector3{0, 0, 0}},
	{QuaternionFromAngle(VectorUp(), deg2rad(270)), QuaternionFromAngle(VectorUp(), deg2rad(270)), &Vector3{0, 0, 0}},
	{QuaternionFromAngle(VectorForward(), deg2rad(270)), QuaternionFromAngle(VectorForward(), deg2rad(270)), &Vector3{0, 0, 0}},
}

var alignTests = []struct {
	character *Quaternion
	target    *Quaternion
	expected  *Vector3
}{
	{QuaternionFromAngle(VectorUp(), deg2rad(0)), QuaternionFromAngle(VectorUp(), deg2rad(45)), &Vector3{0, 0, 1.9634954084936211}},
	{QuaternionFromAngle(VectorUp(), deg2rad(0)), QuaternionFromAngle(VectorUp(), deg2rad(90)), &Vector3{0, 0, 1.9634954084936207}},
}

func TestAlignNoRotation(t *testing.T) {
	character := NewEntity()
	target := NewEntity()

	for i := range alignNoRotationTests {
		character.Orientation = alignNoRotationTests[i].character
		target.Orientation = alignNoRotationTests[i].target

		character.physics.(*RigidBody).calculateDerivedData(character)
		target.physics.(*RigidBody).calculateDerivedData(target)

		align := NewAlign(character, target, 0.001, 0.002, 0.1)
		steering := align.GetSteering()

		if !steering.angular.Equals(alignNoRotationTests[i].expected) {
			t.Errorf("Expected %v, but got %v for test %d", alignNoRotationTests[i].expected, steering.angular, i+1)
		}
	}
}

func TestAlignRotation(t *testing.T) {
	character := NewEntity()
	target := NewEntity()

	for i := range alignTests {
		character.Orientation = alignTests[i].character
		target.Orientation = alignTests[i].target

		character.physics.(*RigidBody).calculateDerivedData(character)
		target.physics.(*RigidBody).calculateDerivedData(target)

		align := NewAlign(character, target, 0.001, 0.002, 0.1)
		steering := align.GetSteering()

		if !steering.angular.Equals(alignTests[i].expected) {
			t.Errorf("Expected %v, but got %v for test %d", alignTests[i].expected, steering.angular, i+1)
		}
	}

}

