package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/stojg/vivere/client"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var port string

var world *World

func init() {

	rand.Seed(time.Now().UTC().UnixNano())
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	world = NewWorld(true)

	ch := client.NewClientHandler()
	world.SetNewClients(ch.NewClients())
	http.Handle("/ws/", websocket.Handler(ch.Websocket))
	http.HandleFunc("/", webserver)
}

// Main only contains the necessary wiring for bootstrapping the
// engine
func main() {

	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	for a := 0; a < 80; a++ {
		NewThingie(world)
	}

	for a := 0; a < 32; a++ {
		t := NewObstacle(world)
		t.Position[0] = float64(a*32) + 16
		t.Position[1] = -16
	}

	for a := 1; a < 32; a++ {
		t := NewObstacle(world)
		t.Position[0] = float64(a*32) + 16
		t.Position[1] = -640 + 16
	}

	for a := 1; a < 21; a++ {
		t := NewObstacle(world)
		t.Position[0] = 16
		t.Position[1] = float64(a*-32) + 16
	}

	for a := 1; a < 21; a++ {
		t := NewObstacle(world)
		t.Position[0] = 1024 - 16
		t.Position[1] = float64(a*-32) + 16
	}

	world.GameLoop()
}

func NewThingie(world *World) {
	ent := world.entities.NewEntity()
	ent.Model = 2
	physics := NewParticlePhysics()
	physics.InvMass = rand.Float64()*5 + 0.1
	ent.geometry = &Circle{Radius: 15}
	ent.physics = physics
	ent.input = NewBunnyAI(ent.physics)
	ent.graphics = NewBunnyGraphic()
	ent.Position.Set(rand.Float64()*960-32, rand.Float64()*-576-32, 0)
	ent.Orientation = (rand.Float64() * math.Pi * 2) - math.Pi
}

func NewObstacle(world *World) *Entity {
	ent := world.entities.NewEntity()
	ent.Model = 1
	physics := NewParticlePhysics()
	physics.InvMass = 0
	ent.geometry = &Rectangle{HalfSize: Vector3{16, 16, 16}}
	ent.physics = physics
	ent.Position.Set(rand.Float64()*800+64, rand.Float64()*-400+64, 0)
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
