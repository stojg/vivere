package engine

import (
	"github.com/stojg/vivere/lib/websocket"
	"github.com/stojg/vivere/lib/entity"
	"encoding/json"
	"time"
)

type World struct {
	entities [1024]entity.Entity
}

func (w *World) Init() {
	w.entities[0] = entity.Entity{"bunny", 0, time.Now()}
}

func (w *World) ProcessInput() {

}

func (w *World) Update(elapsed time.Duration) {
	w.entities[0].Rotation += 0.1
}

func (w *World) Render(now time.Time) {
	w.entities[0].Timestamp = now
	b, _ := json.Marshal(w.entities[0])
	websocket.H.Broadcast <- b
}
