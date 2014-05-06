package main

import (
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
	g.physics = &NullComponent{}
	g.graphics = &NullComponent{}
	g.input = &NullComponent{}
	g.collision = &NullComponent{}
	gol.Add(g)
	return g
}

func (gol *EntityList) Add(i *Entity) bool {
	_, found := gol.set[i.ID()]
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
	position    *v.Vec
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
