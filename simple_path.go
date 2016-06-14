package main

import (
	"container/heap"
)

func NewGridGraph(width, height int) *GridGraph {
	return &GridGraph{
		width:   width,
		height:  height,
		tiles:   make([]bool, width*height),
		weights: make([]float64, width*height),
	}
}

type GridNode struct {
	x int
	y int
}

// GridGraph contains all the nodes that should be searched over and can find
// neighbours of a node and the cost between node.
//
// inspired by http://www.redblobgames.com/pathfinding/a-star/implementation.html
type GridGraph struct {
	width   int
	height  int
	tiles   []bool
	weights []float64
}

func (graph GridGraph) Add(x, y int) {
	graph.tiles[graph.width*y+x] = true
}

// Neighbours
func (graph *GridGraph) Neighbours(x, y int) [][2]int {
	edges := [8][2]int{
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
	for _, corner := range edges {
		if graph.inGrid(corner[0], corner[1]) {
			results = append(results, corner)
		}
	}
	return results
}

func (graph *GridGraph) Cost(fromX, fromY, toX, toY int) float64 {
	return 1
}

// InBounds will return true/false if an x and y is inside the graph
// and is 'walkable'
func (graph *GridGraph) inGrid(x, y int) bool {
	if 0 > x || 0 > y || x >= graph.width || y >= graph.height {
		return false
	}
	return graph.tiles[graph.width*y+x]
}

// PathFindingNode gets inserted into a priority queue.
type PathFindingNode struct {
	// Cost contains the total cost from the start node to this node
	Cost float64
	// This contains the estimated cost to the target
	EstimatedCost float64
	// Unique ID that references the item in a navigational graph
	ID int
	// Closed is set to true if this is node is in the closed list, maybe not needed?
	Closed bool
	// index is an internal reference to which index in the PathFindingQueue
	// this items sits in
	index int
}

// PathFindingQueue is priority queue that keeps PathFindingNodes sorted
// by their smallest cost
type PathFindingQueue []*PathFindingNode

func (pq PathFindingQueue) Len() int {
	return len(pq)
}

func (pq PathFindingQueue) Less(i, j int) bool {
	return pq[i].Cost < pq[j].Cost
}

func (pq PathFindingQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// http://dave.cheney.net/2014/06/07/five-things-that-make-go-fast
func (pq *PathFindingQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PathFindingNode)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PathFindingQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PathFindingQueue) Update(item *PathFindingNode, priority float64) {
	item.Cost = priority
	heap.Fix(pq, item.index)
}
