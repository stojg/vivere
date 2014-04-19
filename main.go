package main

import (
	"os"
	"log"
	"net/http"
	"time"
	"encoding/json"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// start the hub
	go h.run()

	// Create a bunny
	m := Entity{"bunny", 0, time.Now().UnixNano()}
	go func() {
		// 60frames a second
		timer := time.Tick(16 * time.Millisecond)
		//for now := range timer {
		for now := range timer {
			m.Rotation = m.Rotation + 0.01
			m.Timestamp = now.UnixNano()
			b, _ := json.Marshal(m)
			h.broadcast <- b
		}
	}()

	// Serve the static site content
	http.HandleFunc("/", serveHome)
	// Serve the websocket service
	http.HandleFunc("/ws", serveWebsocket)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
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
