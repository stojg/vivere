package main

import (
	"bytes"
	"encoding/binary"
	"github.com/stojg/vivere/client"
//	"github.com/volkerp/goquadtree/quadtree"
	"log"
	"math"
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

func NewWorld(debug bool) *World {
	w := &World{}
	w.entities = &EntityList{}
	w.FPS = 120
	w.debug = debug
	w.collision = &CollisionDetector{}
	return w
}

func (world *World) GameLoop() {
	ticker := time.NewTicker(time.Duration(int(1e9) / int(world.FPS)))
	previousTime := time.Now()
	for {
		select {
		case <-ticker.C:

//			qT := quadtree.NewQuadTree(quadtree.NewBoundingBox(0,1000,0,-1000))

//			for _, a := range world.entities.GetAll() {
//				qT.Add(a)
//			}

			// Get the elapsed time since the last tick
			currentTime := time.Now()
			elapsedTime := float64(currentTime.Sub(previousTime)/time.Millisecond) / 1000
			previousTime = currentTime

			world.Tick += 1

			for _, entity := range world.entities.GetAll() {
				entity.Update(elapsedTime)
			}

			world.ResolveCollisions(world.Collisions(), elapsedTime)

			// Send world state updates to the clients
			if math.Mod(float64(world.Tick), 6) == 0 {
				state := world.Serialize()
				for _, p := range world.players {
					p.Update(state)
				}
			}
			// Ping the clients every second to get the RTT
			if math.Mod(float64(world.Tick), float64(world.FPS)) == 0 {
				for _, p := range world.players {
					p.Ping()
				}
			}

			for _, entity := range world.entities.GetAll() {
				entity.physics.(*ParticlePhysics).ClearForces()
				entity.physics.(*ParticlePhysics).ClearRotations()
			}

		case newPlayer := <-world.newPlayerChan:
			world.players = append(world.players, newPlayer)
			world.Log("[+] New client connected")
		}
	}
}

func (w *World) Collisions() []*Collision {
	collisions := make([]*Collision, 0)
	for aIdx, a := range world.entities.GetAll() {
		for bIdx := aIdx + 1; bIdx <= uint16(len(world.entities.GetAll())); bIdx++ {
			collision, hit := w.collision.Detect(a, world.entities.Get(bIdx))
			if hit {
				collisions = append(collisions, collision)
			}
		}
	}
	return collisions
}

func (w *World) ResolveCollisions(collisions []*Collision, duration float64) {
	for _, pair := range collisions {
		pair.Resolve(duration)
	}
}

func (w *World) SetNewClients(e chan *client.Client) {
	w.newPlayerChan = e
}

func (w *World) Log(message string) {
	if w.debug {
		log.Println(message)
	}
}

func (w *World) Serialize() *bytes.Buffer {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, float32(w.Tick))
	for _, entity := range w.entities.GetAll() {
		if entity.Changed() {
			buf.Write(entity.Serialize().Bytes())
		}
	}
	return buf
}
