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
	"time"
)

var port string

var world *World

// Main only contains the necessary wiring for bootstrapping the
// engine
func main() {
	initWorld()

	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
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
	world.GameLoop()
}

func initWorld() {
	rand.Seed(time.Now().UTC().UnixNano())
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	world = NewWorld(true, 3200, 3200)
	ch := client.NewClientHandler()
	world.newPlayerChan = ch.NewClients()

	c := &creator.Creator{}
	c.Seed(time.Now().UnixNano())
	c.Init(32, int(world.sizeX/32), int(world.sizeY/32))
	world.SetMap(c.GetMap())

	for a := 0; a < 100; a++ {
		ent := NewPray(world)
		for world.Collision(ent) {
			ent.Position.Set(rand.Float64()*1000-500, ent.Scale[1]/2, rand.Float64()*-1000-500)
		}
	}

	hunter := NewHunter(world)
	for world.Collision(hunter) {
		hunter.Position.Set(rand.Float64()*1000-500, hunter.Scale[1]/2, rand.Float64()*-1000-500)
	}

	log.Println("World generated!")

	http.Handle("/ws/", websocket.Handler(ch.Websocket))
	http.HandleFunc("/", webserver)

}

func NewPray(world *World) *Entity {

	spawnSizeX := float64(world.sizeX) * 0.8
	spawnSizeY := float64(world.sizeY) * 0.8

	halfX := spawnSizeX / 2
	halfY := spawnSizeY / 2

	ent := world.entities.NewEntity()
	ent.Model = 2
	ent.MaxAcceleration = 100
	ent.MaxSpeed = 50
	ent.Scale.Set(15, 15, 15)
	ent.geometry = &Rectangle{HalfSize: *ent.Scale.Clone().Scale(0.5)}
	ent.physics = NewParticlePhysics(rand.Float64()*5 + 0.1)
	ent.input = NewSimpleAI(world)
	ent.graphics = NewBunnyGraphic()
	ent.Position.Set(rand.Float64()*spawnSizeX-halfX, ent.Scale[1]/2-1, rand.Float64()*spawnSizeY-halfY)
	ent.Orientation = (rand.Float64() * math.Pi * 2) - math.Pi
	return ent
}

func NewHunter(world *World) *Entity {

	spawnSizeX := float64(world.sizeX) * 0.8
	spawnSizeY := float64(world.sizeY) * 0.8

	halfX := spawnSizeX / 2
	halfY := spawnSizeY / 2

	ent := world.entities.NewEntity()
	ent.Model = 3
	ent.Scale.Set(30, 30, 30)
	//ent.geometry = &Circle{Radius: 15}
	ent.geometry = &Rectangle{HalfSize: *ent.Scale.Clone().Scale(0.5)}
	ent.physics = NewParticlePhysics(0.1)
	ent.input = NewHunterAI(world)
	ent.MaxAcceleration = 100
	ent.MaxSpeed = 100
	ent.graphics = NewBunnyGraphic()
	ent.Position.Set(rand.Float64()*spawnSizeX-halfX, ent.Scale[1]/2-1, rand.Float64()*spawnSizeY-halfY)
	ent.Orientation = (rand.Float64() * math.Pi * 2) - math.Pi
	return ent
}

// webserver is a http.HandleFunc for serving static files over http
func webserver(w http.ResponseWriter, r *http.Request) {
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
