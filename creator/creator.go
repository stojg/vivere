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

type Tile struct {
	Size  int
	x     int
	y     int
	Value float64
}

func NewTile(size, x, y int) *Tile {
	t := &Tile{}
	t.Size = size
	t.x = x
	t.y = y
	return t
}

func (tile *Tile) Position() (position [2]float64) {
	position[0] = float64(tile.x * tile.Size)
	position[1] = float64(tile.y * tile.Size)
	return position
}

type Creator struct {
	sizeX    int
	sizeZ    int
	tileSize int
	world    [][]*Tile
	seed     int64
}

func (c *Creator) Seed(seed int64 ) {
	c.seed = seed
}

func (c *Creator) Init(tileSize, sizeX, sizeZ int) {
	c.tileSize = tileSize
	c.sizeX = sizeX
	c.sizeZ = sizeZ
	c.world = make([][]*Tile, c.sizeX)

	n := NewPerlinNoise(c.seed)
	for x := range c.world {
		c.world[x] = make([]*Tile, c.sizeZ)
		for y := range c.world[x] {
			c.world[x][y] = NewTile(c.tileSize, x, y)
			v := n.At2d(float64(x)* 0.1, float64(y)* 0.1)
			c.world[x][y].Value = v * 0.5 + 0.5
		}
	}
}

func (c *Creator) GetMap() [][]*Tile {
	return c.world
}

func (c *Creator) Tile(x, y int) (tile *Tile, err error) {
	if x > c.sizeX-1 || x < 0 {
		err = errors.New("X is out of bounds, max " + string(c.sizeX-1))
	}
	if y > c.sizeZ -1 || x < 0 {
		err = errors.New("Y is out of bounds, max " + string(c.sizeZ -1))
	}
	tile = c.world[x][y]
	return
}

func (c *Creator) RandomPosition() []int {
	result := make([]int, 2)
	result[0] = rand.Intn(c.sizeX - 1)
	result[1] = rand.Intn(c.sizeZ - 1)
	return result
}
