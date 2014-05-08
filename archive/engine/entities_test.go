package engine

import (
	"bytes"
	"encoding/binary"
	p "github.com/stojg/vivere/physics"
	"github.com/stojg/vivere/state"
	v "github.com/stojg/vivere/vec"
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
	e.SetRotation(0.4)
	e.Position().Set(20, 30)
	e.Velocity().Set(4, 2)

	e.SetShape(p.NewRectangle(e.Position(), 2, 1))
	e.state = state.DEAD

	buf := &bytes.Buffer{}
	e.Serialize(buf, true)

	var bitmask byte
	var eBitmask byte
	var id uint16
	var model Model
	var rotation float64
	var pos v.Vec
	var vel v.Vec
	var size v.Vec
	var state state.State

	binary.Read(buf, binary.LittleEndian, &bitmask)
	binary.Read(buf, binary.LittleEndian, &id)

	// id
	if e.id != id {
		t.Errorf("Expected e.id %v, but got id %v", e.id, id)
	}

	// model
	eBitmask = bitmask & (1 << 0)
	if eBitmask != 1 {
		t.Errorf("Expected model bitmask %v, but got %v", 1, eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &model)
	if e.model != ENTITY_BUNNY {
		t.Errorf("Expected %v, but got %v", e.model, model)
	}

	// rotation
	eBitmask = bitmask & (1 << 1)
	if eBitmask != 2 {
		t.Errorf("Expected rotation bitmask %v, but got %v", 2, eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &rotation)
	if e.Rotation() != rotation {
		t.Errorf("Expected %v, but got %v", e.Rotation(), rotation)
	}

	// position
	eBitmask = bitmask & (1 << 2)
	if eBitmask != 4 {
		t.Errorf("position not serialized, bitmask %v", eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &pos)
	if e.Position()[0] != pos[0] {
		t.Errorf("Expected %v, but got %v", e.Position()[0], pos[0])
	}

	// velocity
	eBitmask = bitmask & (1 << 3)
	if eBitmask != 8 {
		t.Errorf("Expected velocity bitmask %v, but got %v", 16, eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &vel)
	if e.Velocity()[0] != vel[0] {
		t.Errorf("Expected %v, but got %v", e.Velocity()[0], vel[0])
	}

	// size
	eBitmask = bitmask & (1 << 4)
	if eBitmask != 16 {
		t.Errorf("Expected size bitmask %v, but got %v", 32, eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &size)
	if e.Shape().Size()[0] != size[0] {
		t.Errorf("Expected %v, but got %v", e.Shape().Size()[0], size[0])
	}

	// action
	eBitmask = bitmask & (1 << 5)
	if eBitmask != 32 {
		t.Errorf("Expected action bitmask %v, but got %v", 64, eBitmask)
	}
	binary.Read(buf, binary.LittleEndian, &state)
	if e.state != state {
		t.Errorf("Expected %v, but got %v", e.state, state)
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
	e.Position().Set(10, 20)
	e.UpdatePrev()

	if !e.Position().Equals(e.prev.Position()) {
		t.Errorf("Current and Previous position should be the same")
	}

	e.Position().Set(5, 4)
	if e.Position().Equals(e.prev.Position()) {
		t.Errorf("Current and Previous position should not be the same after resetting pos")
	}

	e.UpdatePrev()
	if !e.Position().Equals(e.prev.Position()) {
		t.Errorf("Current and Previous position should be the same after UpdatePrev()")
	}
}

func TestSerializationPositionChanged(t *testing.T) {
	e := NewEntity(4)
	e.Position().Set(10, 20)

	buf := &bytes.Buffer{}
	e.Serialize(buf, false)

	if size := buf.Len(); size == 0 {
		t.Error("Buffer should not be empty")
	}

	var bitMask byte
	var id uint16
	var pos v.Vec

	binary.Read(buf, binary.LittleEndian, &bitMask)
	if bitMask != 4 {
		t.Errorf("Only the position should have been serialized")
	}

	binary.Read(buf, binary.LittleEndian, &id)
	if e.id != id {
		t.Errorf("Expected %v, but got %v", e.id, id)
	}

	binary.Read(buf, binary.LittleEndian, &pos)
	if e.Position()[0] != pos[0] {
		t.Errorf("Expected %v, but got %v", e.Position()[0], pos[0])
	}

	next, err := buf.ReadByte()
	if err != io.EOF {
		t.Errorf("Buffer should be empty, not contain this '%v' byte", next)
	}
}
