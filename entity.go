package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Component interface {
	Update(ent *Entity, elapsed float64)
}

type NullComponent struct{}

func (nc *NullComponent) Update(ent *Entity, elapsed float64) {}

type EntityList struct {
	set    map[uint16]*Entity
	nextID uint16
}

func NewEntityList() *EntityList {
	gol := &EntityList{}
	gol.set = make(map[uint16]*Entity)
	gol.nextID = 0
	return gol
}

func (gol *EntityList) NewEntity() *Entity {
	gol.nextID++
	g := &Entity{}
	g.id = (gol.nextID)
	g.Position = &Vector3{0, 0, 0}
	g.Orientation = 0
	g.scale = &Vector3{1, 1, 1}
	g.physics = &NullComponent{}
	g.graphics = &NullComponent{}
	g.geometry = &Circle{Position: g.Position, Radius: 15}
	g.input = &NullComponent{}
	gol.Add(g)
	return g
}

func (gol *EntityList) Add(i *Entity) bool {
	_, found := gol.set[i.ID()]
	if gol.set == nil {
		gol.set = make(map[uint16]*Entity)
	}
	gol.set[i.ID()] = i
	return !found
}

func (gol *EntityList) GetAll() map[uint16]*Entity {
	return gol.set
}

func (gol *EntityList) Get(i uint16) *Entity {
	return gol.set[i]
}

func (gol *EntityList) Remove(i uint16) {
	delete(gol.set, i)
}

func (gol *EntityList) Length() int {
	return len(gol.set)
}

type Entity struct {
	id          uint16
	Position    *Vector3
	Orientation float64
	scale       *Vector3
	geometry	Geometry
	input       Component
	physics     Component
	graphics    Component
}

func (g *Entity) ID() uint16 {
	return g.id
}

func (ent *Entity) Update(elapsed float64) {
	ent.input.Update(ent, elapsed)
	ent.physics.Update(ent, elapsed)
	ent.graphics.Update(ent, elapsed)
}

func (ent *Entity) Changed() bool {
	return true
}

type Literal byte

const (
	INST_ENTITY_ID    Literal = 1
	INST_SET_POSITION Literal = 2
	INST_SET_ORIENTATION Literal = 3
)

func (ent *Entity) Serialize() *bytes.Buffer {
	buf := &bytes.Buffer{}
	ent.binaryStream(buf, INST_ENTITY_ID, ent.id)
	ent.binaryStream(buf, INST_SET_POSITION, ent.Position)
	ent.binaryStream(buf, INST_SET_ORIENTATION, ent.Orientation)
	return buf
}

func (ent *Entity) binaryStream(buf *bytes.Buffer, lit Literal, val interface{}) {
	binary.Write(buf, binary.LittleEndian, lit)
	switch val.(type) {
	case uint8:
		binary.Write(buf, binary.LittleEndian, byte(val.(uint8)))
	case uint16:
		binary.Write(buf, binary.LittleEndian, float32(val.(uint16)))
	case float32:
		binary.Write(buf, binary.LittleEndian, float32(val.(float32)))
	case float64:
		binary.Write(buf, binary.LittleEndian, float32(val.(float64)))
	case *Vector3:
		binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[0]))
		binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[1]))
	default:
		panic(fmt.Errorf("%c", val))
	}

}
