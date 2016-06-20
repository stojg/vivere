package components

import (
	. "github.com/stojg/vivere/lib/vector"
)

type EntityType uint16

const (
	ENTITY_NONE EntityType = iota
	ENTITY_GROUND
	ENTITY_BLOCK
	ENTITY_PRAY
	ENTITY_HUNTER
	ENTITY_SCARED
	ENTITY_CAMO
)

func NewModelList() *ModelList {
	return &ModelList{
		entity: make(map[*Entity]*Model),
	}
}

type ModelList struct {
	entity map[*Entity]*Model
}

func (b *ModelList) All() map[*Entity]*Model {
	return b.entity
}

func (b *ModelList) New(toEntity *Entity, w, h, d float64, model EntityType) *Model {
	b.entity[toEntity] = NewModel(w, h, d, model)
	return b.entity[toEntity]
}

func (b *ModelList) Get(fromEntity *Entity) *Model {
	return b.entity[fromEntity]
}

func NewModel(w, h, d float64, model EntityType) *Model {
	return &Model{
		Position:    &Vector3{0, 0, 0},
		Orientation: NewQuaternion(1, 0, 0, 0),
		Model:       model,
		Scale:       &Vector3{w, h, d},
	}
}

type Model struct {
	Position    *Vector3    // Holds the linear position of the rigid body in world space.
	Orientation *Quaternion // Holds the angular orientation of the rigid body in world space.
	Scale       *Vector3    // the size of this entity
	Model       EntityType
}
