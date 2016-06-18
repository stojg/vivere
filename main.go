package main

import (
	"github.com/stojg/vivere/client"
	"github.com/stojg/vivere/creator"
	"golang.org/x/net/websocket"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	//"encoding/json"
	//"fmt"
)

const (
	VERB_NORM  int = 0
	VERB_INFO  int = 1
	VERB_DEBUG int = 2
)

var (
	world     *World
	verbosity int = VERB_NORM
)

func main() {

	envVerbosity := os.Getenv("VERBOSITY")
	if envVerbosity == "" {
		verbosity = 0
	} else {
		verbosity, _ = strconv.Atoi(envVerbosity)
	}

	rand.Seed(time.Now().UnixNano())

	world = NewWorld(true, 3200, 3200)

	var seed int64 = 1465762025024741914
	seed = rand.Int63()
	Printf("Creating world with seed %d\n", seed)
	c := creator.NewCreator(seed, 32, int(world.sizeX/32), int(world.sizeY/32))
	world.addStaticsFromMap(c.Create())

	Println("Creating creatures")

	dragForce := &Drag{k1: 0.05, k2: 0.05 * 0.05}
	for a := 0; a < 100; a++ {
		ent := NewAnt(world, rand.Float64()*world.sizeX-world.sizeX/2, 15/2-1, rand.Float64()*world.sizeY-world.sizeY/2)
		ent.Orientation = QuaternionFromAxisAngle(VectorY(), rand.Float64()*(2*math.Pi)-math.Pi)
		for world.isColliding(ent) {
			dPrintln("Rerolling initial position")
			ent.Position.Set(rand.Float64()*world.sizeX-world.sizeX/2, 15/2-1, rand.Float64()*world.sizeY-world.sizeY/2)
		}
		world.forceRegistry.Add(ent, dragForce)
	}

	Println("Setting up networking")
	ch := client.NewClientHandler()
	newClientChan := ch.NewClients()
	http.Handle("/ws/", websocket.Handler(ch.Websocket))
	http.HandleFunc("/", staticAssets)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	go func() {
		for {
			select {
			case client := <-newClientChan:
				client.Update(world.Serialize(true))
				world.players = append(world.players, client)
				Println("New client connected")
			}
		}
	}()

	Printf("First frame for %d entities", len(world.entities.GetAll()))
	for _, ent := range world.entities.GetAll() {
		ent.Body.ClearAccumulators()
		ent.Body.calculateDerivedData(ent)
	}

	Println("Running gameloop")
	PrintFPS(world)
	world.GameLoop()
}

func NewAnt(world *World, x, y, z float64) *Entity {
	ent := world.entities.NewEntity()
	ent.Position.Set(x, y, z)
	ent.MaxAcceleration = &Vector3{10, 1, 10}
	ent.Type = 2
	ent.Scale.Set(15, 15, 15)
	ent.Geometry = &Rectangle{
		HalfSize: *ent.Scale.NewScale(0.5),
	}
	mass := 10.0
	ent.Body = NewRigidBody(mass)

	it := &Matrix3{}
	it.SetBlockInertiaTensor(&Vector3{1, 1, 1}, mass)
	ent.Body.SetInertiaTensor(it)
	ent.Input = NewSimpleAI()

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
