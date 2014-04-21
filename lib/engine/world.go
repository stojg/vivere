package engine

import (
	"github.com/stojg/vivere/lib/observer"
	"github.com/stojg/vivere/lib/websocket"
	"math/rand"
	"time"
)

type World struct {
	Type     string
	entities [10]*Entity
	Width    int
	Height   int
	events   chan interface{}
}

func NewWorld() *World {
	w := new(World)
	w.Type = "World"
	w.Width = 1000
	w.Height = 600
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

func (world *World) listen() {
	for {
		event := <-world.events
		if event.(websocket.Message).Message == "getState" {
			websocket.Broadcast(world)
		}
	}
}

func (w *World) ProcessInput() {}

func (w *World) Update(elapsed time.Duration) {
	for index := range w.entities {
		w.entities[index].Update(elapsed)
	}
}

func (w *World) Render(now time.Time) {
	for _, element := range w.entities {
		websocket.Broadcast(element)
	}
}

func (w World) Message() *websocket.Message {
	message := new(websocket.Message)
	message.Event = "World"
	message.Message = w
	return message
}
