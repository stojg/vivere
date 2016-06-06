package main

import "testing"

func TestRigidBody_AddForce(t *testing.T) {

	body := NewRigidBody(0)
	force := &Vector3{0, 1, 0}
	body.AddForce(force)

	if !body.forceAccum.Equals(force) {
		t.Errorf("Expected %v, got %v", force, body.forceAccum)
	}
}

func TestRigidBody_Update(t *testing.T) {
	body := NewRigidBody(0.1)
	ent := NewEntity()
	body.AddForce(&Vector3{1, 1, 1})
	body.Update(ent, 1)
	expected := &Vector3{0.99, 0.99, 0.99}
	if ent.Position.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, ent.Position)
	}
}

func TestRigidBody_AddForceAtPoint(t *testing.T) {
	body := NewRigidBody(0.1)
	ent := NewEntity()
	body.AddForceAtPoint(ent, &Vector3{1, 0, 1}, &Vector3{1, 1, 1})
	t.Log(body.torqueAccum)

	body.Update(ent, 1)
	expected := &Vector3{0.99, 0.99, 0.99}
	expectedQuaternion := &Quaternion{1, 0, 0, 1}
	if !ent.Orientation.Equals(expectedQuaternion) {
		t.Errorf("Expected orientation %v, got %v", expectedQuaternion, ent.Orientation)
	}

	if !ent.Position.Equals(expected) {
		t.Errorf("Expected position %v, got %v", expected, ent.Position)
	}

}

//func TestRigidBody_AddForceAtBodyPoint(t *testing.T) {
//	body := NewRigidBody(0.1)
//	ent := NewEntity()
//	body.AddForceAtBodyPoint(ent, &Vector3{1, 0, 0}, &Vector3{1, 0, 0})
//	body.Update(ent, 1)
//	expectedQuaternion := &Quaternion{1,0,0,1}
//	if !ent.Orientation.Equals(expectedQuaternion) {
//		t.Errorf("Expected orientation %v, got %v", expectedQuaternion, ent.Orientation)
//	}
//
//	expected := &Vector3{0.99, 0.99, 0.99}
//	if !ent.Position.Equals(expected) {
//		t.Errorf("Expected position %v, got %v", expected, ent.Position)
//	}
//}
