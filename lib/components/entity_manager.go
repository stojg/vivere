package components

import (
	"math"
)

type Entity uint32

func NewEntityManager() *EntityManager {
	return &EntityManager{}
}

type EntityManager struct {
	nextID   Entity
	entities []*Entity
}

func (e *EntityManager) Create() *Entity {
	id := e.generateNewID()
	e.entities = append(e.entities, &id)
	return &id
}

func (e *EntityManager) generateNewID() Entity {
	if e.nextID < math.MaxUint32 {
		e.nextID++
		return e.nextID
	}
	panic("Out of entity ids, implement GC")
}
