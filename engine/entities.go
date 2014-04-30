package engine

import (
	"bytes"
	"encoding/binary"
	"io"
	p "github.com/stojg/vivere/physics"
)

type Model uint16

const (
	ENTITY_DELETE Model = 0
	ENTITY_WORLD  Model = 1
	ENTITY_BUNNY  Model = 2
)

type Unique interface {
	Id() uint16
}

type Entity struct {
	p.Body
	id         uint16
	model      Model
	prev       *Entity
	action     Action
	left       bool
}

func (e *Entity) SetModel(m Model) {
	e.model = m
}

func (e *Entity) Action() Action {
	return e.action
}

func NewEntity(id uint16) *Entity {
	e := &Entity{}
	e.id = id

	e.action = ACTION_NONE
	e.Position().Set(0, 0)
	e.SetRotation(0.0)
	e.SetMass(1)
	e.SetInertia(4)
	e.SetDamping(0.999)
	e.SetShape(&p.Rectangle{H: 10, W: 20})

	e.prev = &Entity{}
	e.UpdatePrev()

	return e
}

func (e *Entity) UpdatePrev() {
	e.prev.model = e.model
	e.prev.action = e.action

	e.prev.SetShape(e.Shape())
	e.prev.Position().Copy(e.Position())
	e.prev.SetRotation(e.Rotation())
	e.prev.Velocity().Copy(e.Velocity())
}

// Serialize writes a binary representation of this object into a writer
func (e *Entity) Serialize(buf io.Writer, serAll bool) bool {

	bufTemp := &bytes.Buffer{}
	bitMask := make([]byte, 1)

	bitMask[0] |= 0 << uint(0)
	if serAll || e.model != e.prev.model {
		bitMask[0] |= 1 << uint(0)
		binary.Write(bufTemp, binary.LittleEndian, e.model)
	}

	bitMask[0] |= 0 << uint(1)
	if serAll || e.Rotation() != e.prev.Rotation() {
		bitMask[0] |= 1 << uint(1)
		binary.Write(bufTemp, binary.LittleEndian, e.Rotation())
	}

	bitMask[0] |= 0 << uint(2)
	if serAll || !e.Position().Equals(e.prev.Position()) {
		bitMask[0] |= 1 << uint(2)
		for i := range e.Position() {
			binary.Write(bufTemp, binary.LittleEndian, &e.Position()[i])
		}
	}

	bitMask[0] |= 0 << uint(3)
	if serAll || !e.Velocity().Equals(e.prev.Velocity()) {
		bitMask[0] |= 1 << uint(3)
		for i := range e.Velocity() {
			binary.Write(bufTemp, binary.LittleEndian, &e.Velocity()[i])
		}
	}

	bitMask[0] |= 0 << uint(4)
	if serAll || !e.Shape().Size().Equals(e.prev.Shape().Size()) {
		bitMask[0] |= 1 << uint(4)
		size := e.Shape().Size()
		binary.Write(bufTemp, binary.LittleEndian, size[0])
		binary.Write(bufTemp, binary.LittleEndian, size[1])
	}

	bitMask[0] |= 0 << uint(5)
	if serAll || e.action != e.prev.action {
		bitMask[0] |= 1 << uint(5)
		binary.Write(bufTemp, binary.LittleEndian, e.action)
	}

	// Only write if something changed
	if bitMask[0] != 0 {
		// Add the bitmask
		buf.Write(bitMask)
		// Add the id
		binary.Write(buf, binary.LittleEndian, e.id)
		// Write the rest of the values
		buf.Write(bufTemp.Bytes())
		return true
	}
	return false
}
