package main

type Graphics struct {
	sprite string
}

func (c *Graphics) Update(e *Entity, elapsed float64) {

}

func (c *Graphics) SetSprite(sprite string) {
	c.sprite = sprite
}

type BunnyGraphics struct {
	Graphics
	sprite string
}

func NewBunnyGraphic() *BunnyGraphics {
	b := &BunnyGraphics{}
	b.sprite = "bunny"
	return b
}
