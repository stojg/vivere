package main

import (
	"container/heap"
	"math"
)

func NewGridGraph(width, height int) *GridGraph {
	return &GridGraph{
		width:  width,
		height: height,
		nodes:  make([]bool, width*height),
		edges:  make([][][2]int, width*height),
		costs:  make([][]float64, width*height),
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
	width  int
	height int
	nodes  []bool
	edges  [][][2]int
	costs  [][]float64
}

// Add adds a node to the graph
func (graph *GridGraph) Add(x, y int) {
	graph.nodes[x+y*graph.width] = true
}

// Init will pre calculate costs and neighbours, this will be
// need to be called after all nodes has been added to the grid
func (graph *GridGraph) Init() {
	// pre-calculate neighbours
	for y := 0; y < graph.height; y++ {
		for x := 0; x < graph.width; x++ {
			// populate neighbours in row order
			edges := [8][2]int{
				[2]int{x - 1, y - 1},
				[2]int{x, y - 1},
				[2]int{x + 1, y - 1},
				[2]int{x - 1, y},
				[2]int{x + 1, y},
				[2]int{x - 1, y + 1},
				[2]int{x, y + 1},
				[2]int{x + 1, y + 1},
			}

			var corners [][2]int
			for _, corner := range edges {
				if graph.inGrid(corner[0], corner[1]) {
					corners = append(corners, corner)
				}
			}
			graph.edges[x+y*graph.width] = corners
			graph.costs[x+y*graph.width] = make([]float64, len(corners))
			for i := range corners {
				xDiff := math.Abs(float64(x - corners[i][0]))
				yDiff := math.Abs(float64(y - corners[i][1]))

				graph.costs[x+y*graph.width][i] = math.Sqrt(xDiff*xDiff + yDiff*yDiff)
			}
		}
	}
}

// Neighbours
func (graph *GridGraph) Neighbours(id [2]int) [][2]int {
	return graph.edges[id[0]+ id[1]*graph.width]
}

// Cost
func (graph *GridGraph) Cost(id [2]int, neighbour int) float64 {
	return graph.costs[id[0]+ id[1]*graph.width][neighbour]
}

// inGrid will return true/false if an x and y is inside the graph
// and is 'walkable' (added to the graph)
func (graph *GridGraph) inGrid(x, y int) bool {
	if 0 > x || 0 > y || x >= graph.width || y >= graph.height {
		return false
	}
	return graph.nodes[x+y*graph.width]
}

// PathFindingNode gets inserted into a priority queue.
type PathFindingNode struct {
	// From is the node where we came from, for the start node, this will be nil
	From *PathFindingNode
	// Cost contains the total cost from the start node to this node
	Cost float64
	// This contains the estimated cost to the target
	EstimatedCost float64
	// Unique ID that references the item in a navigational graph
	ID [2]int
	// Closed is set to true if this is node is in the closed list, maybe not needed?
	Closed bool
	// index is an internal reference to which index in the PathFindingQueue
	// this items sits in
	index int
}

func PathFinder(graph *GridGraph, start, goal [2]int) []*PathFindingNode {

	closed := make(map[[2]int]bool)

	open := make(PathFindingQueue, 0)
	startNode := &PathFindingNode {
		ID: start,
		From: nil,
	}
	heap.Push(&open, startNode)
	heap.Init(&open)

	var current *PathFindingNode

	for open.Len() > 0 {
		current = heap.Pop(&open).(*PathFindingNode)

		if current.ID == goal {
			break
		}
		neighbours := graph.Neighbours(current.ID)
		for i := range neighbours {
			toNode := neighbours[i]
			// skip if the node is closed
			if _, ok := closed[toNode]; ok {

				continue
			}
			// get the cost estimate for the end node
			costToNeighbour := current.Cost + graph.Cost(current.ID, i)

			// we are here if we need to update the record, update the cost and connection
			// we don't care about duplicate entries
			record := &PathFindingNode {
				ID: toNode,
				Cost: costToNeighbour,
				From: current,
			}
			heap.Push(&open, record)
		}
		closed[current.ID] = true
	}

	// we're here either found the goal or we have no more nodes to search,
	// find out which
	if current == nil || current.ID != goal {
		return nil
	}

	var path []*PathFindingNode

	for current != startNode {
		path = append(path, current)
		current = current.From
	}


	// reverse the list
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

// PathFindingQueue is priority queue that keeps PathFindingNodes sorted
// by their smallest cost
type PathFindingQueue []*PathFindingNode

func (pq PathFindingQueue) Len() int { return len(pq) }

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
