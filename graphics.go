package main

import (
//	"math"
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

}
