package main

import (
	"fmt"
	"math"
)

var (
	entityManager *EntityManager
)

func init() {
	Println("Inititalising EntityManager")
	entityManager = &EntityManager{
		componentsByClass: make(map[*Entity]map[string]Component),
	}
}

type Entity uint32

type EntityManager struct {
	nextID            Entity
	entities          []*Entity
	componentsByClass map[*Entity]map[string]Component
}

func (e *EntityManager) generateNewID() Entity {
	if e.nextID < math.MaxUint32 {
		e.nextID++
		return e.nextID
	}
	panic("Out of entity ids, implement GC")
}

func (e *EntityManager) CreateEntity() *Entity {
	id := e.generateNewID()
	e.entities = append(e.entities, &id)
	return &id
}

func (e *EntityManager) AddComponent(toEntity *Entity, component Component) {
	typeName := fmt.Sprintf("%T", component)
	if _, ok := e.componentsByClass[toEntity]; !ok {
		e.componentsByClass[toEntity] = make(map[string]Component)
	}
	//if _, ok := e.componentsByClass[toEntity][typeName]; !ok {
	//	e.componentsByClass[toEntity][typeName] = make(map[string][typeName])
	//}
	e.componentsByClass[toEntity][typeName] = component
}

func (e *EntityManager) EntityComponent(forEntity *Entity, t string) Component {
	if _, ok := e.componentsByClass[forEntity][t]; ok {
		return e.componentsByClass[forEntity][t]
	}
	return nil
}

func (e *EntityManager) EntitiesWith(t string) []*Entity {
	var list []*Entity
	for ent, val := range e.componentsByClass {
		if _, ok := val[t]; ok {
			list = append(list, ent)
		}
	}
	return list
}
