package components

import (
	. "github.com/stojg/vivere/lib/vector"
	"sync"
)

func NewCollisionList() *CollisionList {
	return &CollisionList{
		entity: make(map[*Entity]*Collision),
	}
}

type CollisionList struct {
	sync.Mutex
	entity map[*Entity]*Collision
}

func (b *CollisionList) All() map[*Entity]*Collision {
	result := make(map[*Entity]*Collision)
	b.Lock()
	for k,v := range b.entity {
		result[k] = v
	}
	b.Unlock()
	return result
}

func (b *CollisionList) New(toEntity *Entity, x, y, z float64) *Collision {
	b.Lock()
	b.entity[toEntity] = &Collision{
		Geometry: &Rectangle{
			HalfSize: Vector3{x / 2, y / 2, z / 2},
		},
		//Geometry: &Circle{
		//	Radius: x,
		//},
	}
	b.Unlock()
	return b.entity[toEntity]
}

func (b *CollisionList) Get(fromEntity *Entity) *Collision {
	b.Lock()
	t := b.entity[fromEntity]
	b.Unlock()
	return t
}

type Collision struct {
	Geometry interface{}
}

type Circle struct {
	Radius float64
}

type Rectangle struct {
	HalfSize Vector3
	MinPoint Vector3
	MaxPoint Vector3
}

// ToWorld sets the min and max points of this rectangle in world coordinates
func (r *Rectangle) ToWorld(position *Vector3) {
	r.MinPoint[0] = position[0] - r.HalfSize[0]
	r.MaxPoint[0] = position[0] + r.HalfSize[0]
	r.MinPoint[1] = position[1] - r.HalfSize[1]
	r.MaxPoint[1] = position[1] + r.HalfSize[1]
	r.MinPoint[2] = position[2] - r.HalfSize[2]
	r.MaxPoint[2] = position[2] + r.HalfSize[2]
}
