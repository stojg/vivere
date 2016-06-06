package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	. "github.com/volkerp/goquadtree/quadtree"
)

type EntityType uint16

const (
	ENTITY_NONE EntityType = iota
	ENTITY_BLOCK
	ENTITY_PRAY
	ENTITY_HUNTER
	ENTITY_SCARED
	ENTITY_CAMO
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
	g.ID = (gol.nextID)
	gol.Add(g)
	return g
}

func NewEntity() *Entity {
	ent := &Entity{}
	ent.Position = &Vector3{0, 0, 0}
	ent.Orientation = &Quaternion{}
	ent.Velocity = &Vector3{}
	ent.Rotation = 0
	ent.MaxAcceleration = 10
	ent.MaxSpeed = 40
	ent.MaxRotation = &Vector3{0.1, 0.1, 0.1}
	ent.Scale = &Vector3{15, 15, 15}
	ent.physics = &NullComponent{}
	ent.graphics = &NullComponent{}
	ent.input = &NullComponent{}
	ent.physics = NewRigidBody(5)
	ent.prevPosition = &Vector3{0, 0, 0}
	return ent
}

func (gol *EntityList) Add(i *Entity) bool {
	_, found := gol.set[i.ID]
	if gol.set == nil {
		gol.set = make(map[uint16]*Entity)
	}
	gol.set[i.ID] = i
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
	ID uint16
	// Holds the linear position of the rigid body in world space.
	Position *Vector3
	// Holds the angular orientation of the rigid body in world space.
	Orientation *Quaternion
	// Holds the linear velocity of the rigid body in world space.
	Velocity        *Vector3
	Rotation        float64
	MaxAcceleration float64
	MaxSpeed        float64
	MaxRotation     *Vector3
	Type            EntityType
	Scale           *Vector3
	Dead            bool
	geometry        interface{}
	input           Component
	physics         Component
	graphics        Component
	changed         bool
	prevPosition    *Vector3
	prevOrientation *Quaternion
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

func (ent *Entity) Update(elapsed float64) {
	//ent.Position[1] = ent.Scale[1]/2 - 1
	ent.prevPosition.Set(ent.Position[0], ent.Position[1], ent.Position[2])
	ent.prevOrientation = ent.Orientation
	ent.changed = false

	ent.input.Update(ent, elapsed)
	ent.physics.Update(ent, elapsed)
	ent.graphics.Update(ent, elapsed)

	//if ent.prevPosition.Equals(ent.Position) == false || !ent.prevOrientation.Equals(ent.Orientation) {
	ent.changed = true
	//}
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
	ent.binaryStream(buf, INST_ENTITY_ID, ent.ID)
	ent.binaryStream(buf, INST_SET_POSITION, ent.Position)
	ent.binaryStream(buf, INST_SET_ORIENTATION, ent.Orientation)
	ent.binaryStream(buf, INST_SET_TYPE, ent.Type)
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
	case EntityType:
		binary.Write(buf, binary.LittleEndian, float32(val.(EntityType)))
	case float32:
		binary.Write(buf, binary.LittleEndian, float32(val.(float32)))
	case float64:
		binary.Write(buf, binary.LittleEndian, float32(val.(float64)))
	case *Vector3:
		binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[0]))
		binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[1]))
		binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[2]))
	case *Quaternion:
		binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).r))
		binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).i))
		binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).j))
		binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).k))
	default:
		panic(fmt.Errorf("%c", val))
	}

}
