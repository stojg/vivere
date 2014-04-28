package main

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"io"
	//"log"
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
	inv_mass        float64
	inertia         float64
	inverse_inertia float64
}

//Rock       Density : 0.6  Restitution : 0.1
//Wood       Density : 0.3  Restitution : 0.2
//Metal      Density : 1.2  Restitution : 0.05
//BouncyBall Density : 0.3  Restitution : 0.8
//SuperBall  Density : 0.3  Restitution : 0.95
//Pillow     Density : 0.1  Restitution : 0.2
//Static     Density : 0.0  Restitution : 0.4

type Material struct {
	density     float64
	restitution float64
}

type Transform struct {
	position Vec
	rotation float64
}

type Body struct {
	shape        *Shape
	tx           Transform
	material     Material
	mass_data    MassData
	velocity     Vec
	force        Vec
	gravityScale float64
}

type Entity struct {
	id         Id
	model      Model
	rotation   float32
	angularVel float32
	pos        *Vec
	vel        *Vec
	size       *Vec

	mass    float64
	invMass float64

	maxVel    float64
	maxVelAcc float64

	prev       *Entity
	controller Controller
	action     Action
	element    *list.Element
}

func NewEntity(id Id) *Entity {
	e := &Entity{}
	e.id = id
	e.pos = &Vec{0, 0}
	e.vel = &Vec{0, 0}
	e.size = &Vec{0, 0}
	e.mass = 100
	e.invMass = 1 / e.mass
	e.maxVel = 10
	e.maxVelAcc = 1
	e.controller = &PController{}
	e.action = ACTION_NONE

	e.prev = &Entity{}
	e.prev.id = id
	e.prev.pos = &Vec{0, 0}
	e.prev.vel = &Vec{0, 0}
	e.prev.size = &Vec{0, 0}
	e.prev.action = ACTION_NONE
	return e
}

// Update will call this entity's controller to find an action and then
// update the internal state
func (e *Entity) Update(elapsed int64) {
	elapsedSecond := float32(elapsed) / 1000
	input := e.controller.GetAction(e)

	//log.Println(input.force)

	// Symplectic Euler
	//	velocity := input.force.Scale(e.invMass).Scale(float64(elapsed))
	velocity := input.force.Scale(e.invMass).Scale(float64(elapsed))
	e.pos.Add(velocity.Scale(float64(elapsed)))

	e.rotation = e.rotation + (e.angularVel * elapsedSecond)
	//	e.pos = e.pos.Add(e.vel)
}

func (e *Entity) UpdatePrev() {
	e.prev.model = e.model
	e.prev.rotation = e.rotation
	e.prev.angularVel = e.angularVel
	e.prev.pos.Copy(e.pos)
	e.prev.vel.Copy(e.vel)
	e.prev.size.Copy(e.size)
	e.prev.action = e.action
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

	bitMask[0] |= 0 << uint(6)
	if serAll || e.action != e.prev.action {
		bitMask[0] |= 1 << uint(6)
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
