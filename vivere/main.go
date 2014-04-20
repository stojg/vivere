package main

import (
	"github.com/stojg/vivere/lib/webserver"
	"github.com/stojg/vivere/lib/websocket"
	"github.com/stojg/vivere/lib/engine"
	"log"
	"net/http"
	"os"
	"time"
	//"fmt"
)

const FRAMES_PER_SECOND = 30

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// start the hub
	go websocket.H.Run()

	go func() {
		// Serve the static site content
		http.HandleFunc("/", webserver.ServeStatic)
		// Serve the websocket service
		http.HandleFunc("/ws", websocket.Serve)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	w := new(engine.World)
	w.Init()


	previous := time.Now();
	c := time.Tick(time.Second / FRAMES_PER_SECOND)
	for now := range c {
		elapsed := now.Sub(previous)
		previous = now
		w.ProcessInput()
		w.Update(elapsed)
		w.Render(now)
	}
}


