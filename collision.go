package main

type Collision struct{}

func (c *Collision) Update(e *Entity, elapsed float64) {
	// do nothing
}

type ArcadeCollision struct {
	Collision
}

func NewArcadeCollision() *ArcadeCollision {
	ac := &ArcadeCollision{}
	return ac
}
