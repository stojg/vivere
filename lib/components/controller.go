package components

import "sync"

func NewControllerList() *ControllerList {
	return &ControllerList{
		entity: make(map[*Entity]Controller),
	}
}

type ControllerList struct {
	sync.Mutex
	entity map[*Entity]Controller
}

func (b *ControllerList) All() map[*Entity]Controller {
	b.Lock()
	defer b.Unlock()
	return b.entity
}

func (b *ControllerList) New(toEntity *Entity, cont Controller) Controller {
	b.Lock()
	defer b.Unlock()
	b.entity[toEntity] = cont
	return cont
}

func (b *ControllerList) Get(fromEntity *Entity) Controller {
	b.Lock()
	defer b.Unlock()
	return b.entity[fromEntity]
}

type Controller interface {
	Update(float64)
}
