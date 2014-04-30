package physics

import (
	v "github.com/stojg/vivere/vec"
)

type ForceGenerator interface {
	UpdateForce(e Kinematic, duration float64)
}

type GravityGenerator struct{}

func (gg *GravityGenerator) UpdateForce(e Kinematic, duration float64) {

	center := &v.Vec{500,300}
	center.Sub(e.Position())
	center.Normalize().Scale(4)
	e.AddForce(center)
}

//type SpringGenerator struct{
//	other Forceable
//	restLength float64
//	springConstant float64
//}
//
//func (gg *SpringGenerator) UpdateForce(e Forceable, duration float64) {
//	force := &Vec{}
//	force.Sub(&gg.other.position)
//
//	length := force.Length()
//	length = math.Abs(length - gg.restLength)
//	length *= gg.springConstant
//	force.Normalize().Scale(-length)
//	e.AddForce(force)
//}
