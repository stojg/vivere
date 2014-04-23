package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
)

var bResult io.Writer

func BenchmarkSerialization(b *testing.B) {
	e := NewEntity(5)
	buf := &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		e.Serialize(buf, true)
	}
	bResult = buf
}

func TestSerialization(t *testing.T) {
	e := NewEntity(5)
	e.model = ENTITY_BUNNY
	e.rotation = 0.4
	e.angularVel = 0.1
	e.pos = NewVec(20, 30)
	e.vel = NewVec(4, 2)
	e.size = NewVec(1, 2)
	e.action = ACTION_DIE

	buf := &bytes.Buffer{}
	e.Serialize(buf, true)

	var bitmask byte
	var eBitmask byte
	var id Id
	var model Model
	var rotation float32
	var angularVel float32
	var pos Vec
	var vel Vec
	var size Vec
	var action Action

	binary.Read(buf, binary.LittleEndian, &bitmask)

	binary.Read(buf, binary.LittleEndian, &id)

	if e.id != id {
		t.Errorf("Expected e.id %v, but got id %v", e.id, id)
	}

	eBitmask = bitmask & (1 << 0)
	// property two should have been written
	if eBitmask != 1 {
		t.Errorf("Expected bitmask %v, but got %v", 1, eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &model)
	if e.model != ENTITY_BUNNY {
		t.Errorf("Expected %v, but got %v", e.model, model)
	}

	// rotation
	eBitmask = bitmask & (1 << 1)
	if eBitmask != 2 {
		t.Errorf("Expected bitmask %v, but got %v", 2, eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &rotation)
	if e.rotation != rotation {
		t.Errorf("Expected %v, but got %v", e.rotation, rotation)
	}

	// angularVel
	eBitmask = bitmask & (1 << 2)
	if eBitmask != 4 {
		t.Errorf("Expected bitmask %v, but got %v", 4, eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &angularVel)
	if e.angularVel != angularVel {
		t.Errorf("Expected %v, but got %v", e.rotation, angularVel)
	}

	// position
	eBitmask = bitmask & (1 << 3)
	if eBitmask != 8 {
		t.Errorf("position not serialized, bitmask %v", eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &pos)
	if e.pos[0] != pos[0] {
		t.Errorf("Expected %v, but got %v", e.pos[0], pos[0])
	}

	// velocity
	eBitmask = bitmask & (1 << 4)
	if eBitmask != 16 {
		t.Errorf("Expected bitmask %v, but got %v", 16, eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &vel)
	if e.vel[0] != vel[0] {
		t.Errorf("Expected %v, but got %v", e.vel[0], vel[0])
	}

	// size
	eBitmask = bitmask & (1 << 5)
	if eBitmask != 32 {
		t.Errorf("Expected bitmask %v, but got %v", 32, eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &size)
	if e.size[0] != size[0] {
		t.Errorf("Expected %v, but got %v", e.size[0], size[0])
	}

	// action
	eBitmask = bitmask & (1 << 6)
	if eBitmask != 64 {
		t.Errorf("Expected bitmask %v, but got %v", 64, eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &action)
	if e.action != action {
		t.Errorf("Expected %v, but got %v", e.action, action)
	}
}

func TestSerializationNothingChanged(t *testing.T) {
	e := NewEntity(2)
	buf := &bytes.Buffer{}
	e.Serialize(buf, false)
	if size := buf.Len(); size != 0 {
		t.Errorf("Buffer should be empty, but got %v", size)
	}
}

func TestUpdatePrev(t *testing.T) {
	e := NewEntity(4)
	e.pos = NewVec(10, 20)
	e.UpdatePrev()

	if !e.pos.Equals(e.prev.pos) {
		t.Errorf("Current and Previous position should be the same")
	}

	e.pos.Set(5, 4)
	if e.pos.Equals(e.prev.pos) {
		t.Errorf("Current and Previous position should not be the same after resetting pos")
	}

	e.UpdatePrev()
	if !e.pos.Equals(e.prev.pos) {
		t.Errorf("Current and Previous position should be the same after UpdatePrev()")
	}
}

func TestSerializationPositionChanged(t *testing.T) {
	e := NewEntity(4)
	e.pos = NewVec(10, 20)

	buf := &bytes.Buffer{}
	e.Serialize(buf, false)

	if size := buf.Len(); size == 0 {
		t.Error("Buffer should not be empty")
	}

	var bitMask byte
	var id Id
	var pos Vec

	binary.Read(buf, binary.LittleEndian, &bitMask)
	if bitMask != 8 {
		t.Errorf("Only the position should have been serialized")
	}

	binary.Read(buf, binary.LittleEndian, &id)
	if e.id != id {
		t.Errorf("Expected %v, but got %v", e.id, id)
	}

	binary.Read(buf, binary.LittleEndian, &pos)
	if e.pos[0] != pos[0] {
		t.Errorf("Expected %v, but got %v", e.pos[0], pos[0])
	}

	next, err := buf.ReadByte()
	if err != io.EOF {
		t.Errorf("Buffer should be empty, not contain this '%v' byte", next)
	}
}
