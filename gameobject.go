package main

import ()

type Component interface{}

type GameObjectList struct {
	set    map[uint16]*GameObject
	nextID uint16
}

func NewGameObjectList() *GameObjectList {
	gol := &GameObjectList{}
	gol.set = make(map[uint16]*GameObject)
	gol.nextID = 0
	return gol
}

func (gol *GameObjectList) NewGameObject() *GameObject {
	gol.nextID++
	g := &GameObject{}
	g.id = (gol.nextID)
	gol.Add(g)
	return g
}

func (gol *GameObjectList) Add(i *GameObject) bool {
	_, found := gol.set[i.ID()]
	gol.set[i.ID()] = i
	return !found
}

func (gol *GameObjectList) Get(i uint16) *GameObject {
	return gol.set[i]
}

func (gol *GameObjectList) Remove(i uint16) {
	delete(gol.set, i)
}

func (gol *GameObjectList) Length() int {
	return len(gol.set)
}

type GameObject struct {
	id        uint16
	physics   Component
	graphics  Component
	input     Component
	collision Component
}

func (g *GameObject) ID() uint16 {
	return g.id
}
