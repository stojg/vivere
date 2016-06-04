package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	. "github.com/volkerp/goquadtree/quadtree"
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
	g := NewEntity()
	g.id = (gol.nextID)
	gol.Add(g)
	return g
}

func NewEntity() *Entity {
	ent := &Entity{}
	ent.Position = &Vector3{0, 0, 0}
	ent.Orientation = 0
	ent.Velocity = &Vector3{}
	ent.Rotation = 0
	ent.MaxAcceleration = 10
	ent.MaxSpeed = 40
	ent.MaxRotation = 10
	ent.Scale = &Vector3{15, 15, 15}
	ent.physics = &NullComponent{}
	ent.graphics = &NullComponent{}
	ent.input = &NullComponent{}
	ent.prevPosition = &Vector3{0, 0, 0}
	return ent
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
	id              uint16
	Position        *Vector3
	Orientation     float64
	Velocity        *Vector3
	Rotation        float64
	MaxAcceleration float64
	MaxSpeed        float64
	MaxRotation     float64
	Scale           *Vector3
	geometry        interface{}
	input           Component
	physics         Component
	graphics        Component
	Model           uint16
	changed         bool
	prevPosition    *Vector3
	prevOrientation float64
	bBox            BoundingBox
}

func (g *Entity) BoundingBox() BoundingBox {
	g.bBox.MinX = g.Position[0] - g.Scale[0]/2
	g.bBox.MaxX = g.Position[0] + g.Scale[0]/2
	g.bBox.MinY = g.Position[1] - g.Scale[1]/2
	g.bBox.MaxY = g.Position[1] + g.Scale[1]/2
	g.bBox.MinZ = g.Position[2] - g.Scale[2]/2
	g.bBox.MaxZ = g.Position[2] + g.Scale[2]/2
	return g.bBox
}

func (g *Entity) ID() uint16 {
	return g.id
}

func (ent *Entity) Update(elapsed float64) {
	ent.Position[1] = ent.Scale[1]/2 - 1
	ent.prevPosition.Set(ent.Position[0], ent.Position[1], ent.Position[2])
	ent.prevOrientation = ent.Orientation
	ent.changed = false

	ent.input.Update(ent, elapsed)
	ent.physics.Update(ent, elapsed)
	ent.graphics.Update(ent, elapsed)

	if ent.prevPosition.Equals(ent.Position) == false || ent.prevOrientation != ent.Orientation {
		ent.changed = true
	}
}

func (ent *Entity) Changed() bool {
	return ent.changed
}

type Literal byte

const (
	INST_ENTITY_ID       Literal = 1
	INST_SET_POSITION    Literal = 2
	INST_SET_ORIENTATION Literal = 3
	INST_SET_TYPE        Literal = 4
	INST_SET_SCALE       Literal = 5
)

func (ent *Entity) Serialize() *bytes.Buffer {
	buf := &bytes.Buffer{}
	ent.binaryStream(buf, INST_ENTITY_ID, ent.id)
	ent.binaryStream(buf, INST_SET_POSITION, ent.Position)
	ent.binaryStream(buf, INST_SET_ORIENTATION, ent.Orientation)
	ent.binaryStream(buf, INST_SET_TYPE, ent.Model)
	ent.binaryStream(buf, INST_SET_SCALE, ent.Scale)
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
		binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[2]))
	default:
		panic(fmt.Errorf("%c", val))
	}

}
