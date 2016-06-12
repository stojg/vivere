package main

type ForceRegistry struct {
	forces []*ForceRegistration
}

func (registry *ForceRegistry) Add(ent *Entity, fg ForceGenerator) {
	registry.forces = append(registry.forces, &ForceRegistration{ent, fg})
}

func (registry *ForceRegistry) UpdateForces(elapsed float64) {
	for _, f := range registry.forces {
		f.forceGenerator.UpdateForce(f.entity, elapsed)
	}
}

type ForceRegistration struct {
	entity         *Entity
	forceGenerator ForceGenerator
}

type ForceGenerator interface {
	UpdateForce(*Entity, float64)
}

type Gravity struct {
	gravity float64
}

func (gen *Gravity) UpdateForce(ent *Entity, elapsed float64) {
	if ent.Body.InvMass == 0 {
		return
	}
	linearForce := VectorY().Inverse().Scale(gen.gravity * ent.Body.Mass())
	ent.Body.AddForce(linearForce)
}

type Drag struct {
	k1 float64 // holds the velocity drag coefficient
	k2 float64 // holds the velocity squared drag coefficient
}

func (gen *Drag) UpdateForce(ent *Entity, elapsed float64) {

	linear := ent.Velocity.Clone()
	linDragCoeff := linear.Length()
	linDragCoeff = gen.k1*linDragCoeff + gen.k2*linDragCoeff
	linear.Normalize().Scale(-linDragCoeff)
	ent.Body.AddForce(linear)

	//angular := ent.Rotation.Clone()
	//angDraCoeff := angular.Length()
	//angDraCoeff = gen.k1 * angDraCoeff + gen.k2 * angDraCoeff
	//angular.Normalize().Scale(-angDraCoeff)
	//ent.Body.AddTorque(angular)
}
