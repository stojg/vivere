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

	state = e.NewGameState()
	simulator = p.NewSimulator()

	//rand.Seed(time.Now().UTC().UnixNano())

	createWorld(state)
	state.UpdatePrev()

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

		// Every game tick
		case <-ticker.C:
			now := time.Now()
			elapsed := float64(now.Sub(tFrame)/time.Millisecond) / 1000
			tFrame = now

			state.IncTick()

			n.GetUpdates(state.Players())
			simulator.Update(state, elapsed)
			if math.Mod(float64(state.Tick()), 3) == 0 {
				SendUpdates()
			}

		// On every new connection
		case cl := <-connectionHandler.NewConn():

			player := n.NewPlayer(state.NextPlayerId(), cl)
			state.AddPlayer(player)

			ent := e.NewEntity(state.NextEntityID())
			ent.SetModel(e.ENTITY_BUNNY)
			ent.Position().Set(rand.Float64()*1000, rand.Float64()*600)

			state.AddEntity(ent)
			simulator.Forceregistry.Add(ent, player)

			buf := &bytes.Buffer{}
			state.Serialize(buf, true)

			players := make([]*n.Player, 1, 1)
			players[0] = player
			n.Send(players, buf)
		}
	}
}

// Send to clients
func SendUpdates() {
	buf := &bytes.Buffer{}
	state.Serialize(buf, false)
	n.Send(state.Players(), buf)
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

func disconnect(p *n.Player) {
	n.Disconnect(p)
	state.RemovePlayer(p)
}

func createWorld(state *e.GameState) {
	for a := 0; a < 30; a++ {
		ent := e.NewEntity(state.NextEntityID())
		ent.SetModel(e.ENTITY_BUNNY)
		ent.Position().Set(rand.Float64()*1000, rand.Float64()*600)
		ent.SetRotation(3.14)
		state.AddEntity(ent)

		simulator.Forceregistry.Add(ent, &ai.Simple{})
		//simulator.Forceregistry.Add(ent, &p.GravityGenerator{})
	}
}
