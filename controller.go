package main

type ControllerSystem struct{}

func (s *ControllerSystem) Update(elapsed float64) {

	for _, move := range controllerList.All() {
		move.Update(elapsed)
	}

}
