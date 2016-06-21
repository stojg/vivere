package main

import (
	"bytes"
	"encoding/binary"
	. "github.com/stojg/vivere/lib/components"
	"github.com/stojg/vivere/lib/vector"
	"math"
	"math/rand"
)

var (
	entities      *EntityManager
	modelList     *ModelList
	collisionList *CollisionList
	rigidList     *RigidBodyList
)

func NewLevel() *Level {

	x := 3200.0
	y := 3200.0

	entities = NewEntityManager()
	modelList = NewModelList()
	rigidList = NewRigidBodyManager()
	collisionList = NewCollisionList()

	ground := entities.Create()
	modelList.New(ground, x, 0.1, y, ENTITY_GROUND)

	var dudeList []*Entity
	for i := 0; i < 100; i++ {
		e := entities.Create()
		dudeList = append(dudeList, e)

		body := modelList.New(e, 8, 24, 8, ENTITY_PRAY)
		body.Position.Set(x*rand.Float64()-x/2, 8, rand.Float64()*y-y/2)
		phi := rand.Float64() * math.Pi * 2
		body.Orientation.RotateByVector(&vector.Vector3{math.Cos(phi), 0, math.Sin(phi)})

		rig := rigidList.New(e, 1)
		rig.MaxAcceleration = &vector.Vector3{10, 0, 10}

		collisionList.New(e, 8, 24, 8)
	}

	e := entities.Create()
	m := modelList.New(e, 12, 32, 12, ENTITY_BLOCK)
	m.Position.Set(0, 16, 0)
	rigidList.New(e, 0.0)
	collisionList.New(e, 12, 32, 12)

	lvl := &Level{}
	lvl.systems = append(lvl.systems, &PhysicSystem{})
	//@todo make an AI component?
	lvl.systems = append(lvl.systems, NewAI(dudeList))
	lvl.systems = append(lvl.systems, &CollisionSystem{})
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

	for id, component := range modelList.All() {
		binaryStream(buf, INST_ENTITY_ID, *id)
		binaryStream(buf, INST_SET_POSITION, component.Position)
		binaryStream(buf, INST_SET_ORIENTATION, component.Orientation)
		binaryStream(buf, INST_SET_TYPE, component.Model)
		binaryStream(buf, INST_SET_SCALE, component.Scale)
	}

	return buf
}
