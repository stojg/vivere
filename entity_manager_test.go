package main

import (
	"testing"
)

func BenchmarkEntityManager_AddComponent(b *testing.B) {

	e := entityManager.CreateEntity()
	health := &BodyComponent{
		Model: ENTITY_PRAY,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entityManager.AddComponent(e, health)
	}
}

func BenchmarkEntityManager_EntityComponent(b *testing.B) {

	e := entityManager.CreateEntity()
	health := &BodyComponent{
		Model: ENTITY_PRAY,
	}
	entityManager.AddComponent(e, health)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entityManager.EntityComponent(e, "*main.BodyComponent")
	}
}

func TestEntityManager_EntitiesWith(t *testing.T) {

	for i := 0; i < 500; i++ {
		e := entityManager.CreateEntity()
		health := &BodyComponent{
			Model: ENTITY_PRAY,
		}
		entityManager.AddComponent(e, health)
	}
	for i := 0; i < 500; i++ {
		e := entityManager.CreateEntity()
		render := &MoveComponent{}
		entityManager.AddComponent(e, render)
	}

	expected := 500
	actual := len(entityManager.EntitiesWith("*main.BodyComponent"))

	if expected != actual {
		t.Errorf("expected list to be %d long, got %d entries", expected, actual)
	}
}

func BenchmarkEntityManager_EntitiesWith(b *testing.B) {

	for i := 0; i < 1; i++ {
		e := entityManager.CreateEntity()
		health := &BodyComponent{
			Model: ENTITY_PRAY,
		}
		entityManager.AddComponent(e, health)
	}
	for i := 0; i < 1; i++ {
		entityManager.CreateEntity()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entityManager.EntitiesWith("*main.BodyComponent")
	}
}
