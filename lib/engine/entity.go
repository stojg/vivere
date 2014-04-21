package engine

import (
	"github.com/stojg/vivere/lib/websocket"
	"time"
)

type Position struct {
	X float32
	Y float32
}

type Entity struct {
	Id       int
	Name     string
	Rotation float32
	Position Position
	Created  time.Time
}

func NewEntity(id int, posX, posY, rotation float32) *Entity {
	e := new(Entity)
	pos := Position{posX, posY}
	e.Id = id
	e.Position = pos
	e.Rotation = rotation
	e.Created = time.Now()
	return e
}

func (e Entity) Message() *websocket.Message {
	message := new(websocket.Message)
	message.Event = "Entity"
	message.Message = e
	return message
}

func (e *Entity) Update(elapsed time.Duration) {
	e.Rotation += 0.01
}
