package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/stojg/vivere/client"
	"log"
	"net/http"
	"os"
)

// Main only contains the necessary wiring for bootstrapping the
// engine
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	world := NewWorld(true)
	ch := client.NewClientHandler()
	world.SetNewClients(ch.NewClients())
	http.Handle("/ws/", websocket.Handler(ch.Websocket))
	http.HandleFunc("/", webserver)
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()
	world.GameLoop()
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
