package main

type Command struct {

}

func (t *Command) toJSON() {

}

type Entity struct {
	Name string
	Rotation float32
	Timestamp int64
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

