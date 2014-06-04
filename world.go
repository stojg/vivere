package main

import (
	"bytes"
	"encoding/binary"
	"github.com/stojg/vivere/client"
	"github.com/volkerp/goquadtree/quadtree"
	"log"
	"time"
)

type World struct {
	entities      *EntityList
	players       []*client.Client
	FPS           uint8
	Tick          uint64
	newPlayerChan chan *client.Client
	debug         bool
	collision     *CollisionDetector
}

const (
	SEC_PER_UPDATE  float64 = 0.016
	SEC_PER_MESSAGE float64 = 0.05
)

func NewWorld(debug bool) *World {
	w := &World{}
	w.entities = &EntityList{}
	w.FPS = 120
	w.debug = debug
	w.collision = &CollisionDetector{}
	return w
}

func (world *World) GameLoop() {
	previousTime := time.Now()
	var updateLag float64 = 0
	var msgLag float64 = SEC_PER_UPDATE
	for {
		// Get the elapsed time since the last tick
		currentTime := time.Now()
		elapsedTime := currentTime.Sub(previousTime).Seconds()
		previousTime = currentTime

		updateLag -= elapsedTime
		msgLag += elapsedTime

		qT := quadtree.NewQuadTree(quadtree.NewBoundingBox(0, 1000, 0, -1000))
		for _, entity := range world.entities.GetAll() {
			qT.Add(entity)
		}

		for _, entity := range world.entities.GetAll() {
			entity.Update(elapsedTime)
		}

		// Collisions
		collisions := world.Collisions(&qT)
		for _, pair := range collisions {
			pair.Resolve(elapsedTime)
		}

		for _, entity := range world.entities.GetAll() {
			entity.physics.(*ParticlePhysics).ClearForces()
			entity.physics.(*ParticlePhysics).ClearRotations()
		}
		updateLag -= SEC_PER_UPDATE

		// //Ping the clients every second to get the RTT
		// if math.Mod(float64(world.Tick), float64(world.FPS)) == 0 {
		// for _, p := range world.players {
		// p.Ping()
		// }
		// }

		if msgLag >= SEC_PER_MESSAGE {
			state := world.Serialize(false)
			for _, player := range world.players {
				player.Update(state)
			}
			msgLag -= SEC_PER_MESSAGE
		}

		// Check if the game loop took longer than 16ms
		cycleTime := time.Now().Sub(previousTime).Seconds()
		reminder := SEC_PER_UPDATE - cycleTime
		if reminder > 0 {
			time.Sleep(time.Duration(reminder*1000) * time.Millisecond)
		} else {

		}
	}
}

func (w *World) Collisions(tree *quadtree.QuadTree) []*Collision {
	collisions := make([]*Collision, 0)
	checked := make(map[string]bool, 0)

	for _, a := range world.entities.GetAll() {
		if a.Changed() == false {
			continue
		}
		t := tree.Query(a.BoundingBox())
		for _, b := range t {
			hashA := string(a.id) + ":" + string(b.(*Entity).id)
			hashB := string(b.(*Entity).id) + ":" + string(a.id)
			if checked[hashA] || checked[hashB] {
				continue
			}
			checked[hashA], checked[hashB] = true, true
			collision, hit := w.collision.Detect(a, b.(*Entity))
			if hit {
				collisions = append(collisions, collision)
			}
		}
	}

	return collisions
}

func (w *World) SetNewClients(e chan *client.Client) {
	w.newPlayerChan = e
}

func (w *World) Log(message string) {
	if w.debug {
		log.Println(message)
	}
}

func (w *World) Serialize(serializeAll bool) *bytes.Buffer {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, float32(w.Tick))
	for _, entity := range w.entities.GetAll() {
		if entity.Changed() || serializeAll {
			buf.Write(entity.Serialize().Bytes())
		}
	}
	return buf
}
