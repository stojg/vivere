package main

import (
	"bytes"
	"encoding/binary"
	"math/rand"
)

func NewLevel() *Level {

	x := 3200.0
	y := 3200.0

	ground := entityManager.CreateEntity()
	groundBody := NewBodyComponent(0,0,0, x, 0.1, y)
	groundBody.Model = ENTITY_BLOCK
	entityManager.AddComponent(ground, groundBody)

	for i := 0; i < 100; i++ {
		e := entityManager.CreateEntity()
		entityManager.AddComponent(e, NewBodyComponent(x*rand.Float64()-x/2, 8, rand.Float64()*y-y/2, 8, 24, 8))
		entityManager.AddComponent(e, NewMoveComponent(1))
	}

	lvl := &Level{}
	lvl.systems = append(lvl.systems, &PhysicSystem{})
	lvl.systems = append(lvl.systems, &AI{})
	return lvl
}

type Level struct {
	systems []System
}

func (l *Level) Update(elapsed float64) {
	for i := range l.systems {
		l.systems[i].Update(elapsed)
	}
}

func (l *Level) Draw() *bytes.Buffer {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, float32(Frame))

	entities := entityManager.EntitiesWith("*main.BodyComponent")
	for i, id := range entities {
		component := entityManager.EntityComponent(entities[i], "*main.BodyComponent")
		binaryStream(buf, INST_ENTITY_ID, *id)
		binaryStream(buf, INST_SET_POSITION, component.(*BodyComponent).Position)
		binaryStream(buf, INST_SET_ORIENTATION, component.(*BodyComponent).Orientation)
		binaryStream(buf, INST_SET_TYPE, component.(*BodyComponent).Model)
		binaryStream(buf, INST_SET_SCALE, component.(*BodyComponent).Scale)
	}

	return buf
}
