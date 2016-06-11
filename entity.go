package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	. "github.com/volkerp/goquadtree/quadtree"
	"math"
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

func NewEntity() *Entity {
	return &Entity{
		Position:    &Vector3{},
		Orientation: &Quaternion{1, 0, 0, 0},
		Velocity:    &Vector3{},
		Rotation:    &Vector3{},

		MaxAcceleration:        &Vector3{1, 1, 1},
		MaxSpeed:               10,
		MaxAngularAcceleration: &Vector3{1, 1, 1},
		MaxRotation:            math.Pi / 2,

		Scale: &Vector3{15, 15, 15},

		Input: &NullComponent{},
		Body:  NewRigidBody(1),

		prevPosition:    &Vector3{0, 0, 0},
		prevOrientation: &Quaternion{1, 0, 0, 0},
	}
}

type Entity struct {
	ID          uint16
	Position    *Vector3    // Holds the linear position of the rigid body in world space.
	Orientation *Quaternion // Holds the angular orientation of the rigid body in world space.
	Velocity    *Vector3    // Holds the linear velocity of the rigid body in world space.
	Rotation    *Vector3    // Holds the angular velocity, or rotation for the rigid body in world space.

	MaxAcceleration        *Vector3 // limits the linear acceleration
	MaxSpeed               float64  // limits the linear velocity
	MaxAngularAcceleration *Vector3 // limits the linear acceleration
	MaxRotation            float64  // limits the angular velocity

	Scale    *Vector3    // the size of this entity
	bBox     BoundingBox // for broad phase collision detection
	geometry interface{} // for narrow phase collision detection

	Type EntityType
	Dead bool

	Changed         bool // did this entity change position or rotate since last frame
	prevPosition    *Vector3
	prevOrientation *Quaternion

	// Components
	Input Component
	Body  *RigidBody
}

func (g *Entity) BoundingBox() BoundingBox {
	g.bBox.MinX = g.Position[0] - g.Scale[0]
	g.bBox.MaxX = g.Position[0] + g.Scale[0]
	g.bBox.MinY = g.Position[1] - g.Scale[1]
	g.bBox.MaxY = g.Position[1] + g.Scale[1]
	g.bBox.MinZ = g.Position[2] - g.Scale[2]
	g.bBox.MaxZ = g.Position[2] + g.Scale[2]
	return g.bBox
}

func (ent *Entity) Update(elapsed float64) {
	ent.prevPosition.Set(ent.Position[0], ent.Position[1], ent.Position[2])
	ent.prevOrientation = ent.Orientation
	ent.Changed = false

	ent.Input.Update(ent, elapsed)
	ent.Body.Update(ent, elapsed)

	if !ent.prevPosition.Equals(ent.Position) || !ent.prevOrientation.Equals(ent.Orientation) {
		ent.Changed = true
	}
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
