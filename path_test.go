package main

import (
	"testing"
	"fmt"
	"encoding/json"
	"os"
	"github.com/stojg/vivere/creator"
	"io/ioutil"
	"math"
)

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

func TestDijkstra_map(t *testing.T) {
	var tiles [][]*creator.Tile
	file, e := ioutil.ReadFile("./testdata/map.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal(file, &tiles)
	graph := NewGraph()
	for x := range tiles {
		for _, tile := range tiles[x] {
			if tile.Value <= 0 {
				addConnsToGraph(graph, tiles, tile)
			}
		}
	}

	list := Dijkstra(graph, tiles[99][99], tiles[50][50])

	pathLength := len(list)
	expected := 164
	if pathLength != expected {
		t.Errorf("Dijkstra should have found a path with %d steps, got %d", expected, pathLength)
	}

}

func addConnsToGraph(graph *Graph, tiles [][]*creator.Tile, tile *creator.Tile) {

	maxX := len(tiles) - 1
	maxY := len(tiles[0]) - 1

	axes := []int{-1, 0, 1,}

	tilePos := tile.Position()

	for _, x := range axes {
		if tile.X + x < 0 || tile.X + x > maxX {
			continue
		}
		for _, y := range axes {
			if tile.Y + y < 0 || tile.Y + y > maxY {
				continue
			}
			connTile := tiles[tile.X + x][tile.Y + y]
			diffX := tilePos[0] - connTile.Position()[0]
			diffY := tilePos[1] - connTile.Position()[1]

			cost := math.Sqrt(diffX * diffX + diffY * diffY)

			graph.Add(tile, connTile, cost)
		}
	}

}

func savemap() {
	mapFile, err := os.Create("./map.json")
	if err != nil {
		fmt.Errorf("opening map file %s\n", err.Error())
	}

	j, jerr := json.MarshalIndent(world.heightMap, "", "  ")
	if jerr != nil {
		fmt.Println("jerr:", jerr.Error())
	}

	_, werr := mapFile.Write(j)
	if werr != nil {
		fmt.Println("werr:", werr.Error())
	}
}
