package main

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"io"
	//	"log"
)

type Model uint16
type Id uint16

const (
	ENTITY_DELETE Model = 0
	ENTITY_WORLD  Model = 1
	ENTITY_BUNNY  Model = 2
)

type MassData struct {
	mass            float64
	invMass         float64
	inertia         float64
	inverse_inertia float64
}

func (md *MassData) InvMass() float64 {
	if md.invMass == 0 {
		md.invMass = 1 / md.mass
	}
	return md.invMass
}

type Material struct {
	density     float64
	restitution float64
}

type Transform struct {
	position Vec
	rotation float64
}

type Body struct {
	shape        Shape
	tx           Transform
	material     Material
	massData     MassData
	velocity     Vec
	maxVelocity  float64
	force        Vec
	gravityScale float64
}

type Entity struct {
	Body
	id         Id
	model      Model
	prev       *Entity
	controller Controller
	action     Action
	element    *list.Element
}

func NewEntity(id Id) *Entity {
	e := &Entity{}
	e.id = id
	e.controller = &PController{}
	e.action = ACTION_NONE

	e.shape = &Rectangle{h: 10, w: 20}
	e.tx = Transform{position: Vec{0, 0}, rotation: 0.0}
	e.material = Material{density: 0.3, restitution: 0.3}
	e.massData = MassData{mass: 20, inertia: 4}

	e.prev = &Entity{}
	e.prev.id = id
	e.prev.action = ACTION_NONE
	e.prev.shape = e.shape
	e.prev.tx = e.tx
	e.prev.material = e.material
	e.prev.massData = e.massData

	return e
}

// Update will call this entity's controller to find an action and then
// update the internal state
func (e *Entity) Update(elapsed int64) {
	input := e.controller.GetAction(e)

	// Symplectic Euler
	velocity := input.force.Scale(e.massData.InvMass() * float64(elapsed))

	e.velocity.Add(velocity)

	e.tx.position.Add(e.velocity.Scale(float64(elapsed)))
	e.action = input.action
	//	e.rotation = e.rotation + (e.angularVel * elapsedSecond)
}

func (e *Entity) UpdatePrev() {
	e.prev.model = e.model
	e.prev.action = e.action
	e.prev.tx.position = e.tx.position
	e.prev.tx.rotation = e.tx.rotation
	e.prev.tx.rotation = e.prev.tx.rotation
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
	if serAll || e.tx.rotation != e.prev.tx.rotation {
		bitMask[0] |= 1 << uint(1)
		binary.Write(bufTemp, binary.LittleEndian, e.tx.rotation)
	}

	bitMask[0] |= 0 << uint(2)
	if serAll || !e.tx.position.Equals(&e.prev.tx.position) {
		bitMask[0] |= 1 << uint(2)
		for i := range e.tx.position {
			binary.Write(bufTemp, binary.LittleEndian, &e.tx.position[i])
		}
	}

	bitMask[0] |= 0 << uint(3)
	if serAll || !e.velocity.Equals(&e.prev.velocity) {
		bitMask[0] |= 1 << uint(3)
		for i := range e.velocity {
			binary.Write(bufTemp, binary.LittleEndian, &e.velocity[i])
		}
	}

	bitMask[0] |= 0 << uint(4)
	if serAll || !e.shape.Size().Equals(e.prev.shape.Size()) {
		bitMask[0] |= 1 << uint(4)
		size := e.shape.Size()
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
