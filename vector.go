package main

import (
	"math"
)

type Vector3 [3]float64

var (
	UnitX = Vector3{1, 0, 0}
	UnitY = Vector3{0, 1, 0}
	UnitZ = Vector3{0, 0, 1}
)

func NewVector3(x, y, z float64) *Vector3 {
	e := &Vector3{}
	e[0] = x
	e[1] = y
	e[2] = z
	return e
}

func VectorZ() *Vector3 {
	return &Vector3{0, 0, 1}
}

func VectorY() *Vector3 {
	return &Vector3{0, 1, 0}
}

func VectorX() *Vector3 {
	return &Vector3{1, 0, 0}
}

func (v *Vector3) Clone() *Vector3 {
	result := &Vector3{}
	result[0] = v[0]
	result[1] = v[1]
	result[2] = v[2]
	return result
}

func (a *Vector3) Set(x, y, z float64) {
	(*a)[0] = x
	(*a)[1] = y
	(*a)[2] = z
}

func (a *Vector3) Copy(b *Vector3) {
	if len(a) != len(b) {
		panic("Can't copy values between two Vec2 with different size")
	}
	(*a)[0] = (*b)[0]
	(*a)[1] = (*b)[1]
	(*a)[2] = (*b)[2]
}

func (a *Vector3) Add(b *Vector3) *Vector3 {
	(*a)[0] += (*b)[0]
	(*a)[1] += (*b)[1]
	(*a)[2] += (*b)[2]
	return a
}

func (a *Vector3) NewAdd(b *Vector3) *Vector3 {
	vec := &Vector3{}
	vec[0] = (*a)[0] + (*b)[0]
	vec[1] = (*a)[1] + (*b)[1]
	vec[2] = (*a)[2] + (*b)[2]
	return vec
}

func (v *Vector3) AddScaledVector(b *Vector3, t float64) *Vector3 {
	(*v)[0] += (*b)[0] * t
	(*v)[1] += (*b)[1] * t
	(*v)[2] += (*b)[2] * t
	return v
}

func (a *Vector3) Sub(b *Vector3) *Vector3 {
	(*a)[0] -= (*b)[0]
	(*a)[1] -= (*b)[1]
	(*a)[2] -= (*b)[2]
	return a
}

func (a *Vector3) NewSub(b *Vector3) *Vector3 {
	vec := &Vector3{}
	vec[0] = (*a)[0] - (*b)[0]
	vec[1] = (*a)[1] - (*b)[1]
	vec[2] = (*a)[2] - (*b)[2]
	return vec
}

func (a *Vector3) Length() float64 {
	return math.Sqrt(a[0]*a[0] + a[1]*a[1] + a[2]*a[2])
}

func (a *Vector3) SquareLength() float64 {
	return a[0]*a[0] + a[1]*a[1] + a[2]*a[2]
}

func (a *Vector3) Normalize() *Vector3 {
	length := a.Length()
	if length > 0 {
		a.Scale(1 / length)
	}
	return a
}

func (v *Vector3) Scale(alpha float64) *Vector3 {
	(*v)[0] *= alpha
	(*v)[1] *= alpha
	(*v)[2] *= alpha
	return v
}

func (v *Vector3) Dot(b *Vector3) float64 {
	return (*v)[0]*(*b)[0] + (*v)[1]*(*b)[1] + (*v)[2]*(*b)[2]
}

func (v *Vector3) Clear() *Vector3 {
	v[0] = 0
	v[1] = 0
	v[2] = 0
	return v
}

// VectorProduct aka cross product
func (v *Vector3) VectorProduct(vector *Vector3) *Vector3 {
	return &Vector3{
		v[1]*vector[2] - v[2]*vector[1],
		v[2]*vector[0] - v[0]*vector[2],
		v[0]*vector[1] - v[1]*vector[0],
	}
}

func (v *Vector3) Cross(vector *Vector3) *Vector3 {
	return v.VectorProduct(vector)
}

// ScalarProduct calculates and returns the scalar product of this vector
// with the given vector.
func (v *Vector3) ScalarProduct(vector *Vector3) float64 {
	return v[0]*vector[0] + v[1]*vector[1] + v[2]*vector[2]
}

func (v *Vector3) ComponentProduct(vector *Vector3) *Vector3 {
	result := &Vector3{}
	result[0] = v[0] * vector[0]
	result[1] = v[1] * vector[1]
	result[2] = v[2] * vector[2]
	return result
}

func (v *Vector3) Equals(z *Vector3) bool {
	if v[0] != z[0] {
		return false
	}
	if v[1] != z[1] {
		return false
	}
	if v[2] != z[2] {
		return false
	}
	return true
}

func (v *Vector3) AsOrientation() float64 {
	return math.Atan2(v[0], v[1])
}

func OrientationAsVector(orientation float64) *Vector3 {

	vec := &Vector3{}
	vec[0] = math.Sin(orientation)
	vec[1] = 0
	vec[2] = math.Cos(orientation)
	return vec
}
