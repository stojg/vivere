package physics

import (
	v "github.com/stojg/vivere/vec"
	. "gopkg.in/check.v1"
	"testing"
)

type TestKinematic struct {
	position *v.Vec
	velocity *v.Vec
	forces   *v.Vec
}

func NewTestKinematic(posX, posY, velX, velY float64) *TestKinematic {
	tk := &TestKinematic{}
	tk.position = &v.Vec{posX, posY}
	tk.velocity = &v.Vec{velX, velY}
	tk.forces = &v.Vec{0, 0}
	return tk
}

func (tk *TestKinematic) InvMass() float64 {
	return 1
}
func (tk *TestKinematic) Position() *v.Vec {
	return tk.position
}
func (tk *TestKinematic) Velocity() *v.Vec {
	return tk.velocity
}
func (tk *TestKinematic) AddForce(v *v.Vec) {
	tk.forces.Add(v)
}

func (tk *TestKinematic) Forces() *v.Vec {
	return tk.forces
}

func (tk *TestKinematic) ClearForces() {
	tk.forces = &v.Vec{0, 0}
}

func (tk *TestKinematic) Damping() float64 {
	return 1
}

type EntityList struct {
	entitities []Kinematic
}

func (tl *EntityList) Entities() []Kinematic {
	return tl.entitities
}

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})
