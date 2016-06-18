package creator

import (
	"errors"
	"math/rand"
)

type Literal byte

const (
	INST_TILE_ID Literal = 1
	INST_TILE_POSITION    Literal = 2
	INST_TILE_TYPE        Literal = 3
)

func NewCreator(seed int64, tileSize, sizeX, sizeY int) *Creator {
	return &Creator{
		seed: seed,
		tileSize: tileSize,
		sizeX: sizeX,
		sizeY: sizeY,
		tiles: make([][]*Tile, sizeX),
	}
}

type Creator struct {
	sizeX    int
	sizeY    int
	tileSize int
	tiles    [][]*Tile
	seed     int64
}

func (c *Creator) Create() [][]*Tile{
	perlin := NewPerlinNoise(c.seed)

	for x := range c.tiles {
		c.tiles[x] = make([]*Tile, c.sizeY)
		for y := range c.tiles[x] {
			v := perlin.At2d(float64(x)* 0.1, float64(y)* 0.1)
			c.tiles[x][y] = NewTile(c.tileSize, x, y, v, c.sizeY)
		}
	}

	return c.tiles
}

func (c *Creator) Tile(x, y int) (tile *Tile, err error) {
	if x > c.sizeX-1 || x < 0 {
		err = errors.New("X is out of bounds, max " + string(c.sizeX-1))
	}
	if y > c.sizeY -1 || x < 0 {
		err = errors.New("Y is out of bounds, max " + string(c.sizeY -1))
	}
	tile = c.tiles[x][y]
	return
}

func (c *Creator) RandomPosition() []int {
	result := make([]int, 2)
	result[0] = rand.Intn(c.sizeX - 1)
	result[1] = rand.Intn(c.sizeY - 1)
	return result
}
