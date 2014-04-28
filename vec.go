package main

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
	return NewVec(a[0]/length, a[1]/length)
}

func (a *Vec) Length() float64 {
	return math.Sqrt(a[0]*a[0] + a[1]*a[1])
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
	return &Vec{-(*vec)[0], -(*vec)[1]}
}

// Inverted returns an inverted copy of the vector.
func (vec *Vec) Inverted() *Vec {
	vec[0] = -vec[0]
	vec[1] = -vec[1]
	return vec
}

func (res *Vec) Sub(b *Vec) *Vec {
	(*res)[0] = (*res)[0] - (*b)[0]
	(*res)[1] = (*res)[1] - (*b)[1]
	return res
}

func (a *Vec) Equals(b *Vec) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func (v *Vec) ClampLength(length float64) *Vec {
	return v.Normalize().Scale(length)
}

func Dot(a, b *Vec) float64 {
	return (*a)[0]*(*b)[0] + (*a)[1]*(*b)[1]
}

func (v *Vec) Nrm2Sq() float64 {
	return Dot(v, v)
}

func (v *Vec) Scaled(alpha float64) *Vec {
	(*v)[0] = alpha * (*v)[0]
	(*v)[1] = alpha * (*v)[1]
	return v
}

func (v *Vec) Scale(alpha float64) *Vec {
	return &Vec{alpha * (*v)[0], alpha * (*v)[1]}
}
