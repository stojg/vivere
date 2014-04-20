package engine

import (
	"encoding/json"
	"log"
	"time"
)

type Position struct {
	X float32
	Y float32
}

type Entity struct {
	Id        int
	Name      string
	Rotation  float32
	Position  Position
	Timestamp time.Time
}

func (e *Entity) ToMessage() []byte {
	json, err := json.Marshal(e)
	if err != nil {
		log.Println("error:", err)
	}
	return json
}

// Example commands:
//
// - CreateEntity
// - UpdateEntity
// - DestroyEntity
// - WorldInit
// - Ping
//
// Example Entities
// Creature
// Obstacle
