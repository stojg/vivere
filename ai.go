package main

type AI struct {
}

func (s *AI) Update(elapsed float64) {

	entities := entityManager.EntitiesWith("*main.MoveComponent")
	for i := range entities {
		move := entityManager.EntityComponent(entities[i], "*main.MoveComponent").(*MoveComponent)
		move.AddForce(&Vector3{20, 0, 0})
	}
}
