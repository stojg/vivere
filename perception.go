package main

type Perception struct{}

func (p *Perception) WorldDimension() *Vec {
	return &Vec{1000, 600}
}
