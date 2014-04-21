// Package engine is the "router". It has the Update and the Render functions and
// takes care of sending all necessary information to the client about the state of
// the world
package engine

import (
	"github.com/stojg/vivere/lib/observer"
	"github.com/stojg/vivere/lib/websocket"
	"math/rand"
	"time"
)

// World is the basic struct for the game world. It coordinates and have a list of all entities
// that exists in it
type World struct {
	Type     string
	entities [10]*Entity
	Width    int
	Height   int
	events   chan interface{}
}

// NewWorld setups a new world and initializes the pub / sub on events from the client
func NewWorld(width int, height int) *World {
	w := new(World)
	w.Type = "World"
	w.Width = width
	w.Height = height
	w.events = make(chan interface{})
	rand.Seed(243)
	for i := 0; i < 10; i++ {
		x := rand.Float32() * float32(w.Width)
		y := rand.Float32() * float32(w.Height)
		rot := rand.Float32() * 360
		w.entities[i] = NewEntity(i, x, y, rot)
	}
	observer.Subscribe("World", w.events)
	go w.listen()
	return w
}

// ProcessInput recieves the commands from the client and passes that on to the
// relevant entities in the world
// @todo This might be moved to the the listen function?
func (w *World) ProcessInput() {}

// Update walks through all the entities and calls Update on them
func (w *World) Update(elapsed time.Duration) {
	for index := range w.entities {
		w.entities[index].Update(elapsed)
	}
}

// Render sends all updated entities to all clients
func (w *World) Render(now time.Time) {
	for _, element := range w.entities {
		websocket.Broadcast(element)
	}
}

// Message returns a websocket.Message that is used for sending over the websocket
func (w World) Message() *websocket.Message {
	message := new(websocket.Message)
	message.Event = "World"
	message.Message = w
	return message
}

// listen checks for new Events on the event channel
func (world *World) listen() {
	for {
		event := <-world.events
		if event.(websocket.Message).Message == "getState" {
			websocket.Broadcast(world)
		}
	}
}
