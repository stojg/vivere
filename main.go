package main

import (
	"github.com/stojg/vivere/client"
	"github.com/stojg/vivere/creator"
	"golang.org/x/net/websocket"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var world *World

func main() {
	rand.Seed(time.Now().UnixNano())

	ch := client.NewClientHandler()
	world = NewWorld(true, 3200, 3200)

	world.newPlayerChan = ch.NewClients()

	c := creator.NewCreator(time.Now().UnixNano(), 32, int(world.sizeX/32), int(world.sizeY/32))
	c.Create()
	world.SetMap(c.GetMap())

	spawnZone := []float64{
		0.8 * world.sizeX,
		0.8 * world.sizeY,
	}
	spawnZone[1] -= spawnZone[1] / 2

	world.forceRegistry = &ForceRegistry{}

	drag := &Drag{
		k1: 0.01,
		k2: 0.01 * 0.01,
	}

	for a := 0; a < 0; a++ {
		ent := NewPray(world, rand.Float64()*spawnZone[0]-spawnZone[0]/2, 15/2-1, rand.Float64()*spawnZone[1]-spawnZone[1]/2)
		//ent := NewPray(world, 0, 15/2-1, 0)

		for world.Collision(ent) {
			ent.Position.Set(rand.Float64()*spawnZone[0]-spawnZone[0]/2, ent.Scale[1]/2, rand.Float64()*spawnZone[1]-spawnZone[1]/2)
		}

		ent.Orientation = QuaternionFromAxisAngle(VectorY(), rand.Float64()*(2*math.Pi)-math.Pi)
		world.forceRegistry.Add(ent, drag)
	}

	log.Println("world has been generated")

	http.Handle("/ws/", websocket.Handler(ch.Websocket))
	http.HandleFunc("/", staticAssets)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	go func() {
		for {
			select {
			case newPlayer := <-world.newPlayerChan:
				newPlayer.Update(world.Serialize(true))
				world.players = append(world.players, newPlayer)
				world.Log("[+] New client connected")
			}
		}
	}()

	for _, ent := range world.entities.GetAll() {
		ent.Body.ClearAccumulators()
		ent.Body.calculateDerivedData(ent)
	}

	world.GameLoop()
}

func NewPray(world *World, x, y, z float64) *Entity {

	ent := world.entities.NewEntity()
	ent.Position.Set(x, y, z)
	ent.Type = 2
	ent.Scale.Set(15, 15, 15)
	ent.geometry = &Rectangle{
		HalfSize: *ent.Scale.NewScale(0.5),
	}
	mass := 10.0
	ent.Body = NewRigidBody(mass)

	it := &Matrix3{}
	it.SetBlockInertiaTensor(&Vector3{1, 1, 1}, mass)
	ent.Body.SetInertiaTensor(it)
	ent.Input = NewSimpleAI(world)

	return ent
}

func staticAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	if r.URL.Path[1:] == "" {
		http.ServeFile(w, r, "static/index.html")
		return
	}
	http.ServeFile(w, r, "static/"+r.URL.Path[1:])
}
