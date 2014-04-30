// http://www.gamedev.net/page/resources/_/technical/game-programming/multiplayer-pong-with-go-websockets-and-webgl-r3112
package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	ai "github.com/stojg/vivere/ai"
	e "github.com/stojg/vivere/engine"
	n "github.com/stojg/vivere/net"
	p "github.com/stojg/vivere/physics"
	"log"
	"math"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

const (
	FRAMES_PER_SECOND = 60
)

var state *e.GameState
var simulator *p.Simulator



func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rand.Seed(time.Now().UTC().UnixNano())

	state = e.NewGameState()
	simulator = p.NewSimulator()

	createWorld(state)

	connectionHandler := n.NewConnectionHandler()
	http.Handle("/ws/", websocket.Handler(connectionHandler.WsHandler))
	http.HandleFunc("/", serveStatic)

	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	ticker := time.NewTicker(time.Duration(int(1e9) / FRAMES_PER_SECOND))
	tFrame := time.Now()

	for {
		select {

		// Main game loop
		case <-ticker.C:
			now := time.Now()
			elapsed := float64(now.Sub(tFrame)/time.Millisecond) / 1000
			tFrame = now

			state.IncTick()

			n.GetUpdates(state.Players())
			simulator.Update(state, elapsed)

			// only send updates to the clients every third tick (20Hz)
			if math.Mod(float64(state.Tick()), 3) == 0 {
				SendUpdates()
			}

		// New connection
		case cl := <-connectionHandler.NewConn():
			player := connect(cl)
			ent := e.NewEntity(state.NextEntityID())
			ent.SetModel(e.ENTITY_BUNNY)
			ent.Position().Set(rand.Float64()*1000, rand.Float64()*600)

			state.AddEntity(ent)
			simulator.Forceregistry.Add(ent, player)
		}
	}
}

// Send to clients
func SendUpdates() {
	buf := &bytes.Buffer{}
	state.Serialize(buf, false)
	for _, player := range state.Players() {
		err := n.Send(player, buf)
		if err != nil {
			log.Printf("[!] ws.Send() for Player %d - '%s'\n", player.Id(), err)
			disconnect(player)
		}
	}
	state.UpdatePrev()
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
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

func connect(connection *n.ClientConn) *n.Player {
	player := n.NewPlayer(state.NextPlayerId(), connection)
	state.AddPlayer(player)
	return player
}

func disconnect(p *n.Player) {
	n.Disconnect(p)
	state.RemovePlayer(p)
}

func createWorld(state *e.GameState) {
	for a := 0; a < 1; a++ {
		ent := e.NewEntity(state.NextEntityID())
		ent.SetModel(e.ENTITY_BUNNY)
		ent.Position().Set(rand.Float64()*1000, rand.Float64()*600)
		ent.SetRotation(3.14)
		state.AddEntity(ent)
		simulator.Forceregistry.Add(ent, &ai.Simple{})
	}
	state.UpdatePrev()
}
