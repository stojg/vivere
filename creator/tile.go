package creator

type Tile struct {
	Size  int
	x     int
	y     int
	Value float64
	maxX  int
}

func NewTile(size, x, y int, value float64, maxX int) *Tile {
	return &Tile{
		Size: size,
		x: x,
		y: y,
		Value: value,
		maxX: maxX,
	}
}

func (tile *Tile) X() int {
	return tile.x
}

func (tile *Tile) Y() int {
	return tile.y
}

func (tile *Tile) Position() (position [2]float64) {
	position[0] = float64(tile.x * tile.Size)
	position[1] = float64(tile.y * tile.Size)
	return position
}

func (tile *Tile) ID() int {
	return tile.maxX*tile.y + tile.x
}
