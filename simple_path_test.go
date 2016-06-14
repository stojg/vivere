package main

import "testing"

func TestSimpleGraph(t *testing.T) {

	graph := NewGridGraph(4, 2)
	graph.Add(0, 0)
	graph.Add(1, 0)
	graph.Add(2, 0)
	graph.Add(3, 0)

	graph.Add(0, 1)
	graph.Add(1, 1)
	graph.Add(2, 1)
	graph.Add(3, 1)

	actual := graph.Neighbours(3, 0)

	if len(actual) != 3 {
		t.Errorf("%v expected length to be 3", actual)
	}

	actual = graph.Neighbours(3, 0)
	if len(actual) != 3 {
		t.Errorf("%v", actual)
	}

	actual = graph.Neighbours(2, 0)
	if len(actual) != 5 {
		t.Errorf("%v", actual)
	}
}
