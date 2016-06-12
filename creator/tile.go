package creator

type Tile struct {
	Size  int
	X     int
	Y     int
	Value float64
	MaxX  int
}

func NewTile(size, x, y int, value float64, maxX int) *Tile {
	return &Tile{
		Size: size,
		X: x,
		Y: y,
		Value: value,
		MaxX: maxX,
	}
}

func (tile *Tile) Position() (position [2]float64) {
	position[0] = float64(tile.X * tile.Size)
	position[1] = float64(tile.Y * tile.Size)
	return position
}

func (tile *Tile) ID() int {
	return tile.MaxX *tile.Y + tile.X
}
