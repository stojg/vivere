package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/stojg/vivere/client"
	"log"
	"math/rand"
	"net/http"
	"os"
)

var port string

var world *World

func init() {
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

	for a := 0; a < 250; a++ {
		NewRabbit(world)
	}
	world.GameLoop()
}

func NewRabbit(world *World) {
	ent := world.entities.NewEntity()
	ent.physics = NewParticlePhysics()
	ent.input = NewBunnyAI(ent.physics)
	ent.collision = NewArcadeCollision()
	ent.graphics = NewBunnyGraphic()
	ent.Position.Set(rand.Float64()*1000, rand.Float64()*600)
	ent.orientation = 0
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
