package engine

import (
	"github.com/stojg/vivere/lib/websocket"
	"math/rand"
	"time"
)

type World struct {
	entities [10]Entity
}

func (w *World) Init() {
	rand.Seed(42)
	for i := 0; i < 10; i++ {
		pos := Position{rand.Float32() * 1000, rand.Float32() * 600}
		w.entities[i] = Entity{Id: i, Name: "bunny", Rotation: 0, Position: pos}
		websocket.Send(w.entities[i].ToMessage())
	}
}

func (w *World) ProcessInput() {

}

func (w *World) Update(elapsed time.Duration) {
	for index := range w.entities {
		w.entities[index].Rotation += 0.1
		w.entities[index].Position.X += 0.1
	}
}

func (w *World) Render(now time.Time) {
	for _, element := range w.entities {
		websocket.Send(element.ToMessage())
	}
}
