package main

import (
	"github.com/stojg/vivere/client"
	"github.com/stojg/vivere/creator"
	"golang.org/x/net/websocket"
	"log"
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
	//c.Init(32, int(world.sizeX/32), int(world.sizeY/32))
	//world.SetMap(c.GetMap())

	//for a := 0; a < 2; a++ {

	//for world.Collision(ent) {
	//	ent.Position.Set(rand.Float64()*1000-500, ent.Scale[1]/2, rand.Float64()*-1000-500)
	//}
	ent := NewPray(world, 0, 0, 0)
	ent.physics.(*RigidBody).ClearAccumulators()
	ent.physics.(*RigidBody).calculateDerivedData(ent)

	//ent = NewPray(world, 0, 0, 0)
	//ent.physics.(*RigidBody).ClearAccumulators()
	//ent.physics.(*RigidBody).calculateDerivedData(ent)
	//}

	//hunter := NewHunter(world)
	//for world.Collision(hunter) {
	//	hunter.Position.Set(rand.Float64()*1000-500, hunter.Scale[1]/2, rand.Float64()*-1000-500)
	//}

	log.Println("World generated!")

	http.Handle("/ws/", websocket.Handler(ch.Websocket))
	http.HandleFunc("/", webserver)

}

func NewPray(world *World, x, y, z float64) *Entity {

	//spawnSizeX := float64(world.sizeX) * 0.8
	//spawnSizeY := float64(world.sizeY) * 0.8
	//halfX := spawnSizeX / 2
	//halfY := spawnSizeY / 2

	//rotationAngle := 0.0
	//RotationAxis := VectorUp()

	ent := world.entities.NewEntity()
	ent.Position.Set(x, y, z)
	ent.Type = 2
	ent.MaxAcceleration = 100
	ent.MaxSpeed = 50
	ent.Scale.Set(15, 15, 15)
	ent.geometry = &Rectangle{HalfSize: *ent.Scale.Clone().Scale(0.5)}
	mass := 10.0
	ent.physics = NewRigidBody(mass)

	it := &Matrix3{}
	it.SetBlockInertiaTensor(&Vector3{1, 1, 1}, mass)
	ent.physics.(*RigidBody).SetInertiaTensor(it)
	ent.input = NewSimpleAI(world)
	ent.graphics = NewBunnyGraphic()

	//ent.Position.Set(rand.Float64()*spawnSizeX-halfX, ent.Scale[1]/2-1, rand.Float64()*spawnSizeY-halfY)
	// @todo: fix for rigidbody
	return ent
}

func NewHunter(world *World) *Entity {

	spawnSizeX := float64(world.sizeX) * 0.8
	spawnSizeY := float64(world.sizeY) * 0.8

	halfX := spawnSizeX / 2
	halfY := spawnSizeY / 2

	ent := world.entities.NewEntity()
	ent.Type = 3
	ent.Scale.Set(30, 30, 30)
	//ent.geometry = &Circle{Radius: 15}
	ent.geometry = &Rectangle{HalfSize: *ent.Scale.Clone().Scale(0.5)}
	ent.physics = NewParticlePhysics(0.1)
	ent.input = NewHunterAI(world)
	ent.MaxAcceleration = 100
	ent.MaxSpeed = 100
	ent.graphics = NewBunnyGraphic()
	ent.Position.Set(rand.Float64()*spawnSizeX-halfX, ent.Scale[1]/2-1, rand.Float64()*spawnSizeY-halfY)
	// @todo: fix for rigidbody
	// ent.Orientation = (rand.Float64() * math.Pi * 2) - math.Pi
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
