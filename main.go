package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"github.com/stojg/vivere/client"
	"github.com/stojg/vivere/creator"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"runtime/pprof"
	"time"
)

var port string

var world *World

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	world = NewWorld(true, 3200, 3200)
	ch := client.NewClientHandler()
	world.SetNewClients(ch.NewClients())

	c := &creator.Creator{}
	c.Seed(time.Now().UnixNano())
	c.Init(32, int(world.sizeX/32), int(world.sizeY/32))
	world.SetMap(c.GetMap())

	for a := 0; a < 20; a++ {
		ent := NewThingie(world)
		for world.Collision(ent) {
			ent.Position.Set(rand.Float64()*960, rand.Float64()*-576-32, 0)
		}
	}
	http.Handle("/ws/", websocket.Handler(ch.Websocket))
	http.HandleFunc("/", webserver)
}

// Main only contains the necessary wiring for bootstrapping the
// engine
func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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

func NewThingie(world *World) *Entity {
	ent := world.entities.NewEntity()
	ent.Model = 2
	ent.geometry = &Circle{Radius: 15}
	ent.physics = NewParticlePhysics(rand.Float64()*5 + 0.1)
	ent.input = NewSimpleAI(ent.physics)
	ent.graphics = NewBunnyGraphic()
	ent.Position.Set(rand.Float64()*960, rand.Float64()*-576-32, 0)
	ent.Orientation = (rand.Float64() * math.Pi * 2) - math.Pi
	return ent
}

func NewObstacle(world *World) *Entity {
	ent := world.entities.NewEntity()
	ent.Model = 1
	ent.geometry = &Rectangle{HalfSize: Vector3{16, 16, 16}}
	ent.physics = NewParticlePhysics(0)
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
