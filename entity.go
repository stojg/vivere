package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	v "github.com/stojg/vivere/vec"
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
	g.Position = &v.Vec{0, 0}
	g.scale = &v.Vec{1, 1}
	g.physics = &NullComponent{}
	g.graphics = &NullComponent{}
	g.input = &NullComponent{}
	g.collision = &NullComponent{}
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
	Position    *v.Vec
	orientation float64
	scale       *v.Vec
	input       Component
	physics     Component
	collision   Component
	graphics    Component
}

func (g *Entity) ID() uint16 {
	return g.id
}

func (ent *Entity) Update(elapsed float64) {
	ent.input.Update(ent, elapsed)
	ent.physics.Update(ent, elapsed)
	ent.collision.Update(ent, elapsed)
	ent.graphics.Update(ent, elapsed)
}

func (ent *Entity) Changed() bool {
	return true
}

type Literal byte

const (
	INST_ENTITY_ID    Literal = 1
	INST_SET_POSITION Literal = 2
	INST_SET_ROTATION Literal = 3
)

func (ent *Entity) Serialize() *bytes.Buffer {
	buf := &bytes.Buffer{}
	ent.binaryStream(buf, INST_ENTITY_ID, ent.id)
	ent.binaryStream(buf, INST_SET_POSITION, ent.Position)
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
		binary.Write(buf, binary.LittleEndian, float32(val.(float64)))
	case float64:
		binary.Write(buf, binary.LittleEndian, float32(val.(float32)))
	case *v.Vec:
		binary.Write(buf, binary.LittleEndian, float32(val.(*v.Vec)[0]))
		binary.Write(buf, binary.LittleEndian, float32(val.(*v.Vec)[1]))
	default:
		panic(fmt.Errorf("%c", val))
	}

}
