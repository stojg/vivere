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
	//main loop
	for {
		select {
		// Every game tick
		case <-ticker.C:
			getClientInputs()
			processInput()
			update()
			render()
		// On every new connection
		case cl := <-newConn:
			id := newId()
			clients[id] = cl
			login(id)
			buf := &bytes.Buffer{}
			state.Serialize(buf, true)
			if buf.Len() > 0 {
				websocket.Message.Send(cl.ws, buf.Bytes())
			}

		}
	}
}

var removeList = make([]PlayerId, 0)

// Send to clients
func render() {
	buf := &bytes.Buffer{}
	state.Serialize(buf, false)
	if buf.Len() == 0 {
		return
	}
	// trunc the removeList
	removeList = removeList[0:0]
	for id, cl := range clients {
		err := websocket.Message.Send(cl.ws, buf.Bytes())
		if err != nil {
			removeList = append(removeList, id)
			log.Printf("render: Error - '%s'\n", err)
		}
	}
	for _, id := range removeList {
		delete(clients, id)
		disconnect(id)
	}
	copyState()
}

// Update the state of all entities
func update() {
	for e := state.entities.Front(); e != nil; e = e.Next() {
		e.Value.(*Entity).Update()
	}
}

func processInput() {
	for index, playerId := range state.players {
		if active(state.players[index], ACTION_UP) {
			log.Println(playerId, " pressed up")
		}
		if active(state.players[index], ACTION_DOWN) {
			log.Println(playerId, " pressed down")
		}
	}
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
