package main

import (
	"math"
)

type Graphics struct {
	sprite string
}

func (c *Graphics) Update(e *Entity, elapsed float64) {

}

func (c *Graphics) SetSprite(sprite string) {
	c.sprite = sprite
}

type BunnyGraphics struct {
	sprite string
}

func NewBunnyGraphic() *BunnyGraphics {
	b := &BunnyGraphics{}
	b.sprite = "bunny"
	return b
}

func (c *BunnyGraphics) Update(entity *Entity, elapsed float64) {

	targetDirection := math.Atan2(entity.Velocity[0], entity.Velocity[1])
	//deltaOrientation := (entity.Orientation - targetDirection)
	//log.Println(deltaOrientation)
	// Orientation is in the direction
	entity.Orientation = targetDirection
}
