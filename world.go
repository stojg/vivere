package main

import (
	"bytes"
	"encoding/binary"
	"github.com/stojg/vivere/client"
	"github.com/stojg/vivere/creator"
	"math"
	"time"
)

type World struct {
	Frame    uint64
	debug    bool
	sizeX    float64
	sizeY    float64
	tileSize int

	entities      *EntityList
	players       []*client.Client
	collision     *CollisionDetector
	forceRegistry *ForceRegistry
	graph         *GridGraph
}

const (
	SEC_PER_UPDATE     float64 = 0.016
	SEC_PER_CLIENT_MSG float64 = 0.06
)

func NewWorld(debug bool, sizeX, sizeY float64) *World {
	return &World{
		debug:         debug,
		entities:      &EntityList{},
		collision:     &CollisionDetector{},
		forceRegistry: &ForceRegistry{},
		sizeX:         sizeX,
		sizeY:         sizeY,
	}
}

func (w *World) GameLoop() {
	previous := time.Now()
	var msgLag float64 = SEC_PER_CLIENT_MSG

	for {
		w.Frame += 1

		now := time.Now()
		elapsed := now.Sub(previous).Seconds()
		previous = now

		w.forceRegistry.UpdateForces(elapsed)

		for _, entity := range w.entities.GetAll() {
			entity.Update(elapsed)
		}

		collisions := w.collision.Collisions(w.entities.GetAll(), w.entities.QuadTree())
		for _, pair := range collisions {
			pair.Resolve(elapsed)
		}

		msgLag += elapsed
		if msgLag >= SEC_PER_CLIENT_MSG {
			state := w.Serialize(false)
			for _, player := range w.players {
				player.Update(state)
			}
			msgLag -= SEC_PER_CLIENT_MSG
		}
	}
}

func (world *World) addStaticsFromMap(heightMap [][]*creator.Tile) {
	minimalHeight := 0.2

	world.graph = NewGridGraph(100, 100)

	statics := 0
	for x := range heightMap {
		for _, tile := range heightMap[x] {
			height := tile.Value
			if height <= minimalHeight {
				world.graph.Add(tile.X, tile.Y)
			}

			if height < minimalHeight {
				continue
			}
			height = (height - (minimalHeight - 0.01)) * 10
			ent := world.entities.NewEntity()
			ent.Body.InvMass = 0
			ent.Type = ENTITY_BLOCK
			size := float64(tile.Size)
			ent.Scale.Set(size, size*height, size)
			ent.Geometry = &Rectangle{HalfSize: *ent.Scale.NewScale(0.5)}
			posX := tile.Position()[0] - float64(world.sizeX/2)
			posY := tile.Position()[1] - float64(world.sizeY/2)
			ent.Position.Set(posX, ent.Scale[1]/2, posY)
			statics++
		}
	}

	Printf("Added %d static objects", statics)
	world.graph.Init()
}

func (w *World) toTilePosition(pos *Vector3) [2]int {
	x := int(pos[0]+float64(world.sizeX/2)) / 32
	z := int(pos[2]+float64(world.sizeY/2)) / 32
	return [2]int{x, z}
}

func (w *World) toPosition(tilePos [2]int) *Vector3 {
	x := (float64(tilePos[0] * 32)) - float64(world.sizeX/2)
	z := (float64(tilePos[1] * 32)) - float64(world.sizeY/2)
	return &Vector3{x, 0, z}
}

func (w *World) isColliding(a *Entity) bool {
	checked := make(map[string]bool, 0)
	tree := world.entities.QuadTree()
	for _, b := range tree.Query(a.BoundingBox()) {
		if a == b {
			continue
		}

		hashA := string(a.ID) + ":" + string(b.(*Entity).ID)
		hashB := string(b.(*Entity).ID) + ":" + string(a.ID)
		if checked[hashA] || checked[hashB] {
			continue
		}
		checked[hashA], checked[hashB] = true, true
		_, hit := w.collision.Detect(a, b.(*Entity))
		if hit {
			return true
		}
	}
	return false
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

func (w *World) Serialize(serializeAll bool) *bytes.Buffer {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, float32(w.Frame))
	for _, entity := range w.entities.GetAll() {
		if entity.Body.isAwake || serializeAll {
			buf.Write(entity.Serialize().Bytes())
		}
	}
	return buf
}
