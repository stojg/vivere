package main

import (
	"math"
)

type PhysicSystem struct{}

func (s *PhysicSystem) Update(elapsed float64) {

	entities := rigidList.All()
	for i, move := range entities {
		if !move.IsAwake {
			return
		}

		body := modelList.Get(i)
		if body == nil {
			panic("Physic system requires that a *main.BodyComponent have been set")
		}

		// Calculate linear acceleration from force inputs.
		move.LastFrameAcceleration = move.Acceleration.Clone()
		move.LastFrameAcceleration.AddScaledVector(move.ForceAccum, move.InvMass)

		// Calculate angular acceleration from torque inputs.
		angularAcceleration := move.InverseInertiaTensorWorld.TransformVector3(move.TorqueAccum)

		// Adjust velocities
		// Update linear velocity from both acceleration and impulse.
		move.Velocity.AddScaledVector(move.LastFrameAcceleration, elapsed)

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
		move.CalculateDerivedData(body)

		// Clear accumulators.
		move.ClearAccumulators()

		// Update the kinetic energy store, and possibly put the body to sleep.
		if move.CanSleep {
			currentMotion := move.Velocity.ScalarProduct(move.Velocity) + move.Rotation.ScalarProduct(move.Rotation)
			bias := math.Pow(0.5, elapsed)
			motion := bias*move.Motion + (1-bias)*currentMotion
			if motion < move.SleepEpsilon {
				move.IsAwake = false
			}
		} else if move.Motion > 10*move.SleepEpsilon {
			move.Motion = 10 * move.SleepEpsilon
		}
	}
}
