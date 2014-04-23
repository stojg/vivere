// http://www.gamedev.net/page/resources/_/technical/game-programming/multiplayer-pong-with-go-websockets-and-webgl-r3112
package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	FRAMES_PER_SECOND = 30
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rand.Seed(time.Now().UTC().UnixNano())
	http.Handle("/ws/", websocket.Handler(wsHandler))
	http.HandleFunc("/", serveStatic)
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	ticker := time.NewTicker(time.Duration(int(1e9) / FRAMES_PER_SECOND))
	current := time.Now()
	for {
		select {
		// Every game tick
		case <-ticker.C:
			now := time.Now()
			elapsed := int64(now.Sub(current) / time.Millisecond)
			current = now
			state.Tick()
			getClientInputs()
			//processInput()
			update(elapsed)
			render()
		// On every new connection
		case cl := <-newConn:
			login(cl)
			buf := &bytes.Buffer{}
			state.Serialize(buf, true)
			if buf.Len() > 0 {
				websocket.Message.Send(cl.ws, buf.Bytes())
			}

		}
	}
}

var removeList = make([]Id, 0)

// Send to clients
func render() {
	buf := &bytes.Buffer{}
	state.Serialize(buf, false)
	if buf.Len() == 0 {
		return
	}
	// trunc the removeList
	removeList = removeList[0:0]
	for _, player := range state.players {
		err := websocket.Message.Send(player.conn.ws, buf.Bytes())
		if err != nil {
			removeList = append(removeList, player.id)
			log.Printf("[!] ws.Send() for Player %d - '%s'\n", player.id, err)
		}
	}
	for _, id := range removeList {
		disconnect(id)
	}
	copyState()
}

// Update the state of all entities
func update(elapsed int64) {
	for e := state.entities.Front(); e != nil; e = e.Next() {
		e.Value.(*Entity).Update(elapsed)
	}
}

type Controller interface {
	GetAction(e *Entity) Action
}

type PlayerController struct {
	player *Player
}

// GetAction
func (p *PlayerController) GetAction(e *Entity) Action {
	if ActiveCommand(p.player, ACTION_UP) {
		e.vel[1] = -100
	}
	if ActiveCommand(p.player, ACTION_DOWN) {
		e.vel[1] = 100
	}
	if ActiveCommand(p.player, ACTION_LEFT) {
		e.vel[0] = -100
	}
	if ActiveCommand(p.player, ACTION_RIGHT) {
		e.vel[0] = 100
	}
	ClearCommand(p.player)
	return ACTION_NONE
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
