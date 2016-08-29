package components

import (
	"testing"
)

func BenchmarkCollisionList_All(b *testing.B) {

	manager := NewEntityManager()
	list := NewCollisionList()
	list.New(manager.Create(), 0, 0, 0)
	list.New(manager.Create(), 0, 0, 0)
	list.New(manager.Create(), 0, 0, 0)
	list.New(manager.Create(), 0, 0, 0)
	list.New(manager.Create(), 0, 0, 0)

	for i := 0; i < b.N; i++ {
		list.All()
	}

}
