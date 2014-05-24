package main

type Circle struct {
	Radius float64
}

type Rectangle struct {
	HalfSize Vector3
	MinPoint Vector3
	MaxPoint Vector3
}

func (r *Rectangle) ToWorld(position *Vector3) {
	r.MinPoint[0] = position[0] - r.HalfSize[0]
	r.MaxPoint[0] = position[0] + r.HalfSize[0]
	r.MinPoint[1] = position[1] - r.HalfSize[1]
	r.MaxPoint[1] = position[1] + r.HalfSize[1]
	r.MinPoint[2] = position[2] - r.HalfSize[2]
	r.MaxPoint[2] = position[2] + r.HalfSize[2]
}
