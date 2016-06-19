package main

import (
	"math"
)

type PhysicSystem struct{}

func (s *PhysicSystem) Update(elapsed float64) {

	entities := entityManager.EntitiesWith("*main.MoveComponent")
	for i := range entities {
		move := entityManager.EntityComponent(entities[i], "*main.MoveComponent").(*MoveComponent)
		if !move.isAwake {
			return
		}

		en := entityManager.EntityComponent(entities[i], "*main.BodyComponent")
		if en == nil {
			panic("Physic system requires that a *main.BodyComponent have been set")
		}

		body := en.(*BodyComponent)

		// Calculate linear acceleration from force inputs.
		move.lastFrameAcceleration = move.acceleration.Clone()
		move.lastFrameAcceleration.AddScaledVector(move.forceAccum, move.InvMass)

		// Calculate angular acceleration from torque inputs.
		angularAcceleration := move.inverseInertiaTensorWorld.TransformVector3(move.torqueAccum)

		// Adjust velocities
		// Update linear velocity from both acceleration and impulse.
		move.Velocity.AddScaledVector(move.lastFrameAcceleration, elapsed)

		// Update angular velocity from both acceleration and impulse.
		move.Rotation.AddScaledVector(angularAcceleration, elapsed)

		// Impose drag
		move.Velocity.Scale(math.Pow(move.LinearDamping, elapsed))
		move.Rotation.Scale(math.Pow(move.AngularDamping, elapsed))

		// Adjust positions
		// Update linear position
		body.Position.AddScaledVector(move.Velocity, elapsed)
		// Update angular position
		body.Orientation.AddScaledVector(move.Rotation, elapsed)

		// Normalise the orientation, and update the matrices with the new position and orientation
		move.calculateDerivedData(body)

		// Clear accumulators.
		move.ClearAccumulators()

		// Update the kinetic energy store, and possibly put the body to sleep.
		if move.canSleep {
			currentMotion := move.Velocity.ScalarProduct(move.Velocity) + move.Rotation.ScalarProduct(move.Rotation)
			bias := math.Pow(0.5, elapsed)
			motion := bias*move.motion + (1-bias)*currentMotion
			if motion < move.sleepEpsilon {
				move.isAwake = false
			}
		} else if move.motion > 10*move.sleepEpsilon {
			move.motion = 10 * move.sleepEpsilon
		}
	}
}
