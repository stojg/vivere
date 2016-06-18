package main

import (
	"encoding/json"
	"fmt"
	"github.com/stojg/vivere/creator"
	"io/ioutil"
	"math"
	"os"
	"testing"
)

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
	graph.Init()
	actual := graph.Neighbours([2]int{0, 0})
	if len(actual) != 3 {
		t.Errorf("expected 3 neighbours for 0,0, got %d: %v", len(actual), actual)
		return
	}

	if actual[0][0] != 1 || actual[0][1] != 0 {
		t.Errorf("expected first neighbour for 0,0 to be 0,1, got %d,%d", actual[0][0], actual[0][1])
	}

	actual = graph.Neighbours([2]int{3, 0})
	if len(actual) != 3 {
		t.Errorf("expected 3 neighbours for 3,0, got %d: %v", len(actual), actual)
	}

	actual = graph.Neighbours([2]int{2, 0})
	if len(actual) != 5 {
		t.Errorf("expected 5 neighbours for 2,0, got %d: %v", len(actual), actual)
	}

	cost := graph.Cost([2]int{2, 0}, 0)
	if cost != 1 {
		t.Errorf("expected cost between 2,0 and it's first neighbours to be 1, got %f", cost)
	}

	cost = graph.Cost([2]int{2, 0}, 1)
	if cost != 1 {
		t.Errorf("expected cost between 2,0 and it's third neighbours to be %f, got %f", 1, cost)
	}

	cost = graph.Cost([2]int{2, 0}, 2)
	if cost != math.Sqrt(2) {
		t.Errorf("expected cost between 2,0 and it's third neighbours to be %f, got %f", math.Sqrt(2), cost)
	}
}

func BenchmarkGraph_Init(b *testing.B) {
	var tiles [][]*creator.Tile
	file, e := ioutil.ReadFile("./testdata/map.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal(file, &tiles)
	graph := NewGridGraph(100, 100)
	for x := range tiles {
		for _, tile := range tiles[x] {
			if tile.Value <= 0 {
				graph.Add(tile.X, tile.Y)
			}
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.Init()
	}
}

func TestPathFinder(t *testing.T) {

	graph := NewGridGraph(4, 2)
	graph.Add(0, 0)
	graph.Add(1, 0)
	graph.Add(2, 0)
	graph.Add(3, 0)
	graph.Add(0, 1)
	graph.Add(1, 1)
	graph.Add(2, 1)
	graph.Add(3, 1)
	graph.Init()
	actual := graph.Neighbours([2]int{0, 0})
	if len(actual) != 3 {
		t.Errorf("expected 3 neighbours for 0,0, got %d: %v", len(actual), actual)
		return
	}
	list, cost := PathFinder(graph, [2]int{0, 0}, [2]int{3, 1})
	if len(list) != 4 {
		t.Errorf("wrong path: %v\n", list)
		for i := range list {
			t.Errorf("wrong path: %v %f\n", list[i], cost[i])
		}
	}
}

func TestPathFinder_Bigger(t *testing.T) {
	var tiles [][]*creator.Tile
	file, e := ioutil.ReadFile("./testdata/map.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	json.Unmarshal(file, &tiles)

	graph := NewGridGraph(100, 100)
	for x := range tiles {
		for _, tile := range tiles[x] {
			if tile.Value <= 0 {
				graph.Add(tile.X, tile.Y)
			}
		}
	}
	graph.Init()

	cameFrom, _ := PathFinder(graph, [2]int{30, 20}, [2]int{50, 50})

	pathLength := len(cameFrom)
	expected := 57
	if pathLength != expected {
		t.Errorf("PathFinder should have found a path with %d steps, got %d", expected, pathLength)
	}
}

var benchList [][2]int

// BenchmarkPath_Bigger-4                	   10000	    122626 ns/op
// BenchmarkPath_Bigger-4                	   10000	    175879 ns/op
func BenchmarkPath_Bigger(b *testing.B) {
	var tiles [][]*creator.Tile
	file, e := ioutil.ReadFile("./testdata/map.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal(file, &tiles)
	graph := NewGridGraph(100, 100)
	for x := range tiles {
		for _, tile := range tiles[x] {
			if tile.Value <= 0 {
				graph.Add(tile.X, tile.Y)
			}
		}
	}
	graph.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchList, _ = PathFinder(graph, [2]int{0, 0}, [2]int{50, 50})
	}
}
