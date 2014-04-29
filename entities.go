package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Model uint16
type Id uint16

const (
	ENTITY_DELETE Model = 0
	ENTITY_WORLD  Model = 1
	ENTITY_BUNNY  Model = 2
)

type Material struct {
	density     float64
	restitution float64
}

type Body struct {
	shape           Shape
	position        Vec
	rotation        float64
	mass            float64
	invMass         float64
	inertia         float64
	invInertia      float64
	velocity        Vec
	maxVelocity     float64
	acceleration    Vec
	macAcceleration float64
	gravityScale    float64
	damping         float64
}

func (body *Body) InvMass() float64 {
	if body.invMass == 0 {
		body.invMass = 1 / body.mass
	}
	return body.invMass
}

func (body *Body) SetMass(m float64) *Body {
	body.mass = m
	body.invMass = 1 / body.mass
	return body
}

func (body *Body) SetInertia(i float64) *Body {
	body.inertia = i
	body.invInertia = 1 / body.inertia
	return body
}

type Entity struct {
	Body
	id         Id
	model      Model
	prev       *Entity
	controller Controller
	action     Action
	left       bool
}

func NewEntity(id Id) *Entity {
	e := &Entity{}
	e.id = id
	e.controller = &PController{}
	e.action = ACTION_NONE

	e.shape = &Rectangle{h: 10, w: 20}
	e.position = Vec{0, 0}
	e.rotation = 0.0
	e.SetMass(1)
	e.SetInertia(4)
	e.damping = 0.999
	e.gravityScale = 10

	e.prev = &Entity{}
	e.prev.model = e.model
	e.prev.id = id
	e.prev.position = e.position
	e.prev.rotation = e.rotation
	e.prev.action = ACTION_NONE
	e.prev.shape = e.shape
	e.prev.mass = e.mass
	e.prev.invMass = e.invMass
	e.prev.inertia = e.inertia
	e.prev.invInertia = e.invInertia

	return e
}

func (e *Entity) UpdatePrev() {
	e.prev.model = e.model
	e.prev.action = e.action
	e.prev.position = e.position
	e.prev.rotation = e.rotation
	e.prev.velocity = e.prev.velocity
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
	if serAll || !e.position.Equals(&e.prev.position) {
		bitMask[0] |= 1 << uint(2)
		for i := range e.position {
			binary.Write(bufTemp, binary.LittleEndian, &e.position[i])
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
