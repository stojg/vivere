package main

import "testing"

type TestNode struct {
	id int
}

func (t *TestNode) ID() int {
	return t.id
}

func TestDijkstra(t *testing.T) {

	graph := NewGraph()

	nodes := []*TestNode{
		&TestNode{id: 0},
		&TestNode{id: 1},
		&TestNode{id: 2},
		&TestNode{id: 3},
		&TestNode{id: 4},
		&TestNode{id: 5},
		&TestNode{id: 6},
	}

	graph.Add(nodes[0], nodes[1], 1)
	graph.Add(nodes[1], nodes[2], 1)
	graph.Add(nodes[2], nodes[3], 1)
	graph.Add(nodes[3], nodes[4], 1)
	graph.Add(nodes[3], nodes[5], 1) // cul de sac
	graph.Add(nodes[5], nodes[6], 1)

	list := Dijkstra(graph, nodes[0], nodes[6])

	if list == nil {
		t.Errorf("List is empty")
		for _, record := range list {
			t.Logf("%d", record.node.ID())
		}
	}
}

func TestDijkstra_DidNotFindGoal(t *testing.T) {

	graph := NewGraph()

	nodes := []*TestNode{
		&TestNode{id: 0},
		&TestNode{id: 1},
		&TestNode{id: 2},
		&TestNode{id: 3},
		&TestNode{id: 4},
		&TestNode{id: 5},
		&TestNode{id: 6},
	}

	graph.Add(nodes[0], nodes[1], 1)
	graph.Add(nodes[1], nodes[2], 1)
	graph.Add(nodes[2], nodes[3], 1)
	graph.Add(nodes[3], nodes[4], 1)
	graph.Add(nodes[5], nodes[6], 1) // node 4 doesn't connect to 5 or 6

	list := Dijkstra(graph, nodes[0], nodes[6])

	if list != nil {
		t.Errorf("List is not empty")
		for _, record := range list {
			t.Logf("%d", record.node.ID())
		}
	}
}
