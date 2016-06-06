package main

import (
	"bytes"
	"encoding/binary"
	"github.com/stojg/vivere/client"
	"github.com/stojg/vivere/creator"
	"github.com/volkerp/goquadtree/quadtree"
	"log"
	"math"
	"time"
)

type World struct {
	entities      *EntityList
	players       []*client.Client
	Tick          uint64
	newPlayerChan chan *client.Client
	debug         bool
	collision     *CollisionDetector
	heightMap     [][]*creator.Tile
	sizeX         float64
	sizeY         float64
}

const (
	SEC_PER_UPDATE  float64 = 0.016
	SEC_PER_MESSAGE float64 = 0.05
)

func NewWorld(debug bool, sizeX, sizeY float64) *World {
	w := &World{}
	w.entities = &EntityList{}
	w.debug = debug
	w.collision = &CollisionDetector{}
	w.sizeX = sizeX
	w.sizeY = sizeY
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

		world.Tick += 1

		qT := quadtree.NewQuadTree(quadtree.NewBoundingBox(-6400.0, 6400.0, -6400.0, 6400.0))
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
			//} else if world.debug {
			//log.Printf("lag %f", reminder*1000)
		}
	}
}

func (world *World) SetMap(heightMap [][]*creator.Tile) {
	world.heightMap = heightMap
	for x := range world.heightMap {
		world.heightMap[x][0].Value = 0.71
		world.heightMap[x][len(world.heightMap[x])-1].Value = 0.71
	}
	for y := range world.heightMap[0] {
		world.heightMap[0][y].Value = 0.71
		world.heightMap[len(world.heightMap[0])-1][y].Value = 0.71
	}

	for x := range world.heightMap {
		for y := range world.heightMap[x] {
			height := world.heightMap[x][y].Value
			if height < 0.7 {
				continue
			}
			height = (height - 0.70) * 20
			ent := world.entities.NewEntity()
			ent.Type = ENTITY_BLOCK
			size := float64(world.heightMap[x][y].Size)
			ent.Scale.Set(size, size*height, size)
			ent.geometry = &Rectangle{HalfSize: *ent.Scale.Clone().Scale(0.5)}
			ent.physics = NewParticlePhysics(0)
			posX := world.heightMap[x][y].Position()[0] - float64(world.sizeX/2)
			posY := world.heightMap[x][y].Position()[1] - float64(world.sizeY/2)
			ent.Position.Set(posX, ent.Scale[1]/2, posY)
		}
	}
}

func (w *World) Collision(a *Entity) bool {
	qT := quadtree.NewQuadTree(quadtree.NewBoundingBox(-w.sizeX/2, w.sizeX/2, -w.sizeY/2, w.sizeY/2))
	for _, b := range world.entities.GetAll() {
		qT.Add(b)
	}
	for _, b := range world.entities.GetAll() {
		if a == b {
			continue
		}
		_, hit := w.collision.Detect(a, b)
		if hit {
			return true
		}
	}
	return false
}

func (w *World) Collisions(tree *quadtree.QuadTree) []*Collision {
	collisions := make([]*Collision, 0)
	checked := make(map[string]bool, 0)

	for _, a := range world.entities.GetAll() {
		if !a.Changed() {
			continue
		}

		t := tree.Query(a.BoundingBox())
		for _, b := range t {
			if a == b {
				continue
			}

			hashA := string(a.ID) + ":" + string(b.(*Entity).ID)
			hashB := string(b.(*Entity).ID) + ":" + string(a.ID)
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

func (w *World) findClosest(me *Entity, t EntityType) (*Entity, float64) {
	set := w.entities.GetAll()
	var closest *Entity
	closestDist := math.Inf(+1)
	for _, ent := range set {
		if ent.Type != t {
			continue
		}
		distance := ent.Position.NewSub(me.Position).Length()
		if distance < closestDist {
			closest = ent
			closestDist = distance
		}
	}
	if closest != nil {
		return closest, closestDist
	}
	return nil, 0
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
