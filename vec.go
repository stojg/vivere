package main

import (
	"errors"
)

type Vec [2]float64

func NewVec(x, y float64) *Vec {
	e := &Vec{}
	e[0] = x
	e[1] = y
	return e
}

func (a *Vec) Set(x, y float64) *Vec {
	(*a)[0] = x
	(*a)[1] = y
	return a
}

func (res *Vec) Add(a, b *Vec) *Vec {
	(*res)[0] = (*a)[0] + (*b)[0]
	(*res)[1] = (*a)[1] + (*b)[1]
	return res
}

func (a *Vec) Copy(b *Vec) (*Vec, error) {

	if len(a) != len(b) {
		return a, errors.New("Vec: Can't copy values between two Vec with different size")
	}

	(*a)[0] = (*b)[0]
	(*a)[1] = (*b)[1]

	return a, nil
}

func (res *Vec) Sub(a, b *Vec) *Vec {
	(*res)[0] = (*a)[0] - (*b)[0]
	(*res)[1] = (*a)[1] - (*b)[1]
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

func (res *Vec) Clamp(s *Vec) {
	for i := range *res {
		if (*res)[i] > (*s)[i]/2 {
			(*res)[i] = (*s)[i] / 2
		}
		if (*res)[i] < -(*s)[i]/2 {
			(*res)[i] = -(*s)[i] / 2
		}
	}
}

func Dot(a, b *Vec) float64 {
	return (*a)[0]*(*b)[0] + (*a)[1]*(*b)[1]
}

func (v *Vec) Nrm2Sq() float64 {
	return Dot(v, v)
}

func (res *Vec) Scale(alpha float64, v *Vec) *Vec {
	(*res)[0] = alpha * (*v)[0]
	(*res)[1] = alpha * (*v)[1]
	return res
}
