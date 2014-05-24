package main

type Circle struct {
	Radius float64
}

type Rectangle struct {
	SizeX    float64
	SizeY    float64
	SizeZ    float64
	MinPoint struct{ X, Y, Z float64 }
	MaxPoint struct{ X, Y, Z float64 }
}

func (r *Rectangle) ToWorld(position *Vector3) {
	halfX := r.SizeX / 2
	halfY := r.SizeY / 2
	halfZ := r.SizeZ / 2
	r.MinPoint.X = position[0] - halfX
	r.MaxPoint.X = position[0] + halfX
	r.MinPoint.Y = position[1] - halfY
	r.MaxPoint.Y = position[1] + halfY
	r.MinPoint.Z = position[2] - halfZ
	r.MaxPoint.Z = position[2] + halfZ
}
