package main

import (
	"bytes"
	"encoding/binary"
	"github.com/stojg/vivere/client"
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
	collision     *Collision
}

func NewWorld(debug bool) *World {
	w := &World{}
	w.entities = &EntityList{}
	w.FPS = 120
	w.debug = debug
	w.collision = &Collision{}
	return w
}

func (w *World) GameLoop() {
	ticker := time.NewTicker(time.Duration(int(1e9) / int(w.FPS)))
	previousTime := time.Now()
	for {
		select {
		case <-ticker.C:
			// Get the elapsed time since the last tick
			currentTime := time.Now()
			elapsedTime := float64(currentTime.Sub(previousTime)/time.Millisecond) / 1000
			previousTime = currentTime

			w.Tick += 1

			for _, entity := range w.entities.GetAll() {
				entity.Update(elapsedTime)
			}

			w.ResolveCollisions(w.Collisions(), elapsedTime)

			// Send world state updates to the clients
			if math.Mod(float64(w.Tick), 6) == 0 {
				state := w.Serialize()
				for _, p := range w.players {
					p.Update(state)
				}
			}
			// Ping the clients every second to get the RTT
			if math.Mod(float64(w.Tick), float64(w.FPS)) == 0 {
				for _, p := range w.players {
					p.Ping()
				}
			}

			for _, entity := range w.entities.GetAll() {
				entity.physics.(*ParticlePhysics).ClearForces()
			}

		case newPlayer := <-w.newPlayerChan:
			w.players = append(w.players, newPlayer)
			w.Log("[+] New client connected")
		}
	}
}

func (w *World) Collisions() []*CollisionPair {
	collisions := make([]*CollisionPair, 0)
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

func (w *World) ResolveCollisions(collisions []*CollisionPair, duration float64) {
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
