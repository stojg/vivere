package vec

import (
	"errors"
	"math"
)

var (
	// Zero holds a zero vector.
	VecZero = Vec{}

	// UnitX holds a vector with X set to one.
	VecUnitX = Vec{1, 0}
	// UnitY holds a vector with Y set to one.
	VecUnitY = Vec{0, 1}
	// UnitXY holds a vector with X and Y set to one.
	VecUnitXY = Vec{1, 1}

	// MinVal holds a vector with the smallest possible component values.
	VecMinVal = Vec{-math.MaxFloat64, -math.MaxFloat64}
	// MaxVal holds a vector with the highest possible component values.
	VecMaxVal = Vec{+math.MaxFloat64, +math.MaxFloat64}
)

type Vec [2]float64

func NewVec(x, y float64) *Vec {
	e := &Vec{}
	e[0] = x
	e[1] = y
	return e
}

func (a *Vec) Normalize() *Vec {
	length := a.Length()
	if length > 0 {
		a.Scale(1 / length)
	}
	return a
}

func (a *Vec) Length() float64 {
	return math.Sqrt(a[0]*a[0] + a[1]*a[1])
}

func (a *Vec) SquareLength() float64 {
	return (a[0]*a[0] + a[1]*a[1])
}

func (a *Vec) Set(x, y float64) *Vec {
	(*a)[0] = x
	(*a)[1] = y
	return a
}

func (a *Vec) Add(b *Vec) *Vec {
	(*a)[0] += (*b)[0]
	(*a)[1] += (*b)[1]
	return a
}

func (v *Vec) Sub(b *Vec) *Vec {
	(*v)[0] -= (*b)[0]
	(*v)[1] -= (*b)[1]
	return v
}

func (a *Vec) Copy(b *Vec) (*Vec, error) {
	if len(a) != len(b) {
		return a, errors.New("Vec: Can't copy values between two Vec with different size")
	}
	(*a)[0] = (*b)[0]
	(*a)[1] = (*b)[1]
	return a, nil
}

// Invert inverts the vector.
func (vec *Vec) Invert() *Vec {
	vec[0] = -vec[0]
	vec[1] = -vec[1]
	return vec
}

func (a *Vec) Equals(b *Vec) bool {
	if len(a) != len(b) {
		return false
	}
	if (*a)[0] != (*b)[0] {
		return false
	}
	if (*a)[0] != (*b)[0] {
		return false
	}
	return true
}

func (v *Vec) ClampLength(length float64) *Vec {
	return v.Normalize().Scale(length)
}

func Dot(a, b *Vec) float64 {
	return (*a)[0]*(*b)[0] + (*a)[1]*(*b)[1]
}

func (v *Vec) SquaredNormal() float64 {
	return Dot(v, v)
}

func (v *Vec) ComponentProduct(b *Vec) *Vec {
	return &Vec{(*v)[0] * (*b)[0], (*v)[1] * (*b)[1]}
}

func (v *Vec) ComponentProductUpdate(b *Vec) *Vec {
	(*v)[0] *= (*b)[0]
	(*v)[1] *= (*b)[1]
	return v
}

func (v *Vec) Scale(alpha float64) *Vec {
	(*v)[0] *= alpha
	(*v)[1] *= alpha
	return v
}

func (v *Vec) AddScaledVector(b *Vec, t float64) *Vec {
	(*v)[0] += (*b)[0] * t
	(*v)[1] += (*b)[1] * t
	return v
}

func (v *Vec) Clear() *Vec {
	return v.Set(0, 0)
}
