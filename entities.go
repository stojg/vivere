package main

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"io"
)

type Model uint16
type Id uint16

const (
	ENTITY_DELETE Model = 0
	ENTITY_WORLD Model = 1
	ENTITY_BUNNY Model = 2
)

type Entity struct {
	id         Id
	model      Model
	rotation   float32
	angularVel float32
	pos        *Vec
	vel        *Vec
	size       *Vec
	prev       *Entity
	controller Controller
	element    *list.Element
}

func NewEntity(id Id) *Entity {
	e := &Entity{}
	e.id = id
	e.pos = NewVec(0, 0)
	e.vel = NewVec(0, 0)
	e.size = NewVec(0, 0)
	e.prev = &Entity{}
	e.prev.id = id
	e.prev.pos = NewVec(0, 0)
	e.prev.vel = NewVec(0, 0)
	e.prev.size = NewVec(0, 0)
	e.controller = &PlayerController{}
	return e
}

// Update will call this entity's controller to find an action and then
// update the internal state
func (e *Entity) Update(elapsed int64) {
	elapsedSecond := float32(elapsed) / 1000

	e.controller.GetAction(e)

	e.rotation = e.rotation + (e.angularVel * elapsedSecond)
	e.pos.Add(e.pos, e.vel.Scale(float64(elapsedSecond), e.vel))
}

func (e *Entity) UpdatePrev() {
	e.prev.model = e.model
	e.prev.rotation = e.rotation
	e.prev.angularVel = e.angularVel
	e.prev.pos.Copy(e.pos)
	e.prev.vel.Copy(e.vel)
	e.prev.size.Copy(e.size)
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
	if serAll || e.rotation != e.prev.rotation {
		bitMask[0] |= 1 << uint(1)
		binary.Write(bufTemp, binary.LittleEndian, e.rotation)
	}

	bitMask[0] |= 0 << uint(2)
	if serAll || e.angularVel != e.prev.angularVel {
		bitMask[0] |= 1 << uint(2)
		binary.Write(bufTemp, binary.LittleEndian, e.angularVel)
	}

	bitMask[0] |= 0 << uint(3)
	if serAll || !e.pos.Equals(e.prev.pos) {
		bitMask[0] |= 1 << uint(3)
		for i := range e.pos {
			binary.Write(bufTemp, binary.LittleEndian, &e.pos[i])
		}
	}

	bitMask[0] |= 0 << uint(4)
	if serAll || !e.vel.Equals(e.prev.vel) {

		bitMask[0] |= 1 << uint(4)
		for i := range e.vel {
			binary.Write(bufTemp, binary.LittleEndian, &e.vel[i])
		}
	}

	bitMask[0] |= 0 << uint(5)
	if serAll || !e.size.Equals(e.prev.size) {
		bitMask[0] |= 1 << uint(5)
		for i := range e.size {
			binary.Write(bufTemp, binary.LittleEndian, &e.size[i])
		}
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
