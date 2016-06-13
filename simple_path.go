package main

import "container/heap"

func NewSimpleGraph(width, height int) *SimpleGraph {
	return &SimpleGraph{
		width:   width,
		height:  height,
		tiles:   make([]bool, width*height),
		weights: make([]float64, width*height),
	}
}

type GraphNode struct {
	position [2]int
	weight   float64
	index    int
	value int
}

// SimpleGraph inspired by http://www.redblobgames.com/pathfinding/a-star/implementation.html
type SimpleGraph struct {
	width   int
	height  int
	tiles   []bool
	weights []float64

}

func (graph SimpleGraph) Add(x, y int, value float64) {
	graph.tiles[graph.width*y+x] = true
	graph.weights[graph.width*y+x] = value
}

// InBounds will return true/false if an x and y is inside the graph
func (graph *SimpleGraph) InBounds(x, y int) bool {
	if 0 > x || 0 > y || x >= graph.width || y >= graph.height {
		return false
	}
	return graph.tiles[graph.width*y+x]
}

func (graph *SimpleGraph) Neighbours(x, y int) [][2]int {
	corners := [8][2]int{
		[2]int{x - 1, y - 1},
		[2]int{x - 1, y},
		[2]int{x - 1, y + 1},
		[2]int{x, y - 1},
		[2]int{x, y + 1},
		[2]int{x + 1, y - 1},
		[2]int{x + 1, y},
		[2]int{x + 1, y + 1},
	}

	var results [][2]int
	for _, corner := range corners {
		if graph.InBounds(corner[0], corner[1]) {
			results = append(results, corner)
		}
	}
	return results
}

func (graph *SimpleGraph) Cost(fromX, fromY, toX, toY int) float64 {
	return 1
}

type GraphNodeQueue []*GraphNode

func (pq GraphNodeQueue) Len() int { return len(pq) }

func (pq GraphNodeQueue) Less(i, j int) bool {
	return pq[i].weight < pq[j].weight
}

func (pq GraphNodeQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
// http://dave.cheney.net/2014/06/07/five-things-that-make-go-fast
func (pq *GraphNodeQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*GraphNode)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *GraphNodeQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *GraphNodeQueue) update(item *GraphNode, value int, priority float64) {
	item.value = value
	item.weight = priority
	heap.Fix(pq, item.index)
}

