package main

import (
	"fmt"
	"math"
)

type Node interface {
	ID() int
	Position() [2]float64
}

func Dijkstra(graph *Graph, start, goal Node) []*NodeRecord {

	open := &DijkstraPathList{
		list: make(map[int]*NodeRecord, 0),
	}
	closed := &DijkstraPathList{
		list: make(map[int]*NodeRecord, 0),
	}

	open.add(&NodeRecord{
		node: start,
	})

	var current *NodeRecord

	for open.len() > 0 {

		current = open.smallest()
		if current == nil {
			fmt.Printf("current empty?\n")
			break
		}
		if current.node == goal {
			break
		}

		connections := graph.getConnections(current.node)

		for _, connection := range connections {
			toNode := connection.To
			// skip if the node is closed
			if closed.contains(toNode) {
				continue
			}

			// get the cost estimate for the end node
			endNodeCost := current.costSoFar + connection.Cost

			var record *NodeRecord

			if record = open.find(toNode); record != nil {
				// here we find the record in the open list corresponding to the endNode

				// but the cost calculated for it is lower than the existing one, so we continue
				if record.costSoFar <= endNodeCost {
					continue
				}
			} else {
				// otherwise we know we've got an unvisited node now, so make a record for it
				record = &NodeRecord{
					node: toNode,
				}
			}

			// we are here if we need to update the record, update the cost and connection
			record.costSoFar = endNodeCost
			record.connection = current

			// and add it to the open list if it's a new record
			if !open.contains(toNode) {
				open.add(record)
			}
		}

		if current != nil {
			open.remove(current)
			closed.add(current)
		}
	}

	// we're here either found the goal or we have no more nodes to search,
	// find out which
	if current == nil || current.node != goal {
		return nil
	}

	path := make([]*NodeRecord, 0)

	for current.node != start {
		path = append(path, current)
		current = current.connection
	}

	// reverse the list
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

type NodeRecord struct {
	node               Node
	connection         *NodeRecord
	costSoFar          float64
	estimatedTotalCost float64
	closed             bool
}

// http://theory.stanford.edu/~amitp/GameProgramming/Heuristics.html
func heuristic(from, to Node) float64 {
	diffX := math.Abs(from.Position()[0] - to.Position()[0])
	diffY := math.Abs(from.Position()[1] - to.Position()[1])
	return math.Sqrt(diffX*diffX + diffY*diffY)
}

func AStar(graph *Graph, start, goal Node) []*NodeRecord {

	open := &AstarPathList{
		list: make(map[int]*NodeRecord, 0),
	}
	closed := &AstarPathList{
		list: make(map[int]*NodeRecord, 0),
	}

	open.add(&NodeRecord{
		node:               start,
		estimatedTotalCost: heuristic(start, goal),
	})

	var current *NodeRecord

	opened := 0

	for open.len() > 0 {

		current = open.smallest()

		if current == nil {
			fmt.Printf("current empty?\n")
			break
		}

		if current.node == goal {
			break
		}

		connections := graph.getConnections(current.node)

		for _, connection := range connections {
			toNode := connection.To

			// get the cost estimate for the end node
			endNodeCost := current.costSoFar + connection.Cost
			var endNodeHeuristicCost float64
			var record *NodeRecord
			// skip if the node is closed
			if record = closed.find(toNode); record != nil {
				if record.costSoFar <= endNodeCost {
					continue
				}
				closed.remove(record)
				endNodeHeuristicCost = current.costSoFar - connection.Cost
			} else if record = open.find(toNode); record != nil {
				// here we find the record in the open list corresponding to the endNode

				// but the cost calculated for it is lower than the existing one, so we continue
				if record.costSoFar <= endNodeCost {
					continue
				}
				endNodeHeuristicCost = current.costSoFar - connection.Cost
			} else {
				// otherwise we know we've got an unvisited node now, so make a record for it
				record = &NodeRecord{
					node: toNode,
				}
				endNodeHeuristicCost = heuristic(toNode, goal)
			}

			// we are here if we need to update the record, update the cost and connection
			record.costSoFar = endNodeCost
			record.connection = current
			record.estimatedTotalCost = endNodeCost + endNodeHeuristicCost

			// and add it to the open list if it's a new record
			if !open.contains(toNode) {
				open.add(record)
			}
		}

		if current != nil {
			opened++
			open.remove(current)
			closed.add(current)
		}
	}

	// we're here either found the goal or we have no more nodes to search,
	// find out which
	if current == nil || current.node != goal {
		return nil
	}

	path := make([]*NodeRecord, 0)

	for current.node != start {
		path = append(path, current)
		current = current.connection
	}

	// reverse the list
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func NewGraph() *Graph {
	return &Graph{
		Connections: make(map[int][]*Connection),
	}
}

type Graph struct {
	Connections map[int][]*Connection
}

func (g *Graph) Add(from, to Node, cost float64) {
	g.Connections[from.ID()] = append(g.Connections[from.ID()], &Connection{
		From: from,
		To:   to,
		Cost: cost,
	})
}

func (g *Graph) getConnections(node Node) []*Connection {
	return g.Connections[node.ID()]
}

type Connection struct {
	From Node
	To   Node
	Cost float64
}

// todo, make it more performant
type DijkstraPathList struct {
	list map[int]*NodeRecord
}

// returns the NodeRecord with the smallest cost so far
func (p *DijkstraPathList) smallest() *NodeRecord {
	var smallestRecord *NodeRecord
	var minCost = math.MaxFloat64
	for i := range p.list {
		if p.list[i].costSoFar < minCost {
			smallestRecord = p.list[i]
			minCost = p.list[i].costSoFar
		}
	}
	return smallestRecord
}

func (p *DijkstraPathList) add(n *NodeRecord) {
	if n == nil {
		panic(fmt.Sprintf("list.add() trying to add nil NodeRecord"))
	}
	p.list[n.node.ID()] = n
}

func (p *DijkstraPathList) remove(n *NodeRecord) {
	delete(p.list, n.node.ID())
}

func (p *DijkstraPathList) len() int {
	return len(p.list)
}

func (p DijkstraPathList) contains(n Node) bool {
	if _, ok := p.list[n.ID()]; ok {
		return true
	}
	return false
}

func (p DijkstraPathList) find(n Node) *NodeRecord {
	return p.list[n.ID()]
}

type AstarPathList struct {
	list map[int]*NodeRecord
}

// returns the NodeRecord with the smallest cost so far
func (p *AstarPathList) smallest() *NodeRecord {
	var smallestRecord *NodeRecord
	var minCost = math.MaxFloat64
	for i := range p.list {
		if p.list[i].estimatedTotalCost < minCost {
			smallestRecord = p.list[i]
			minCost = p.list[i].estimatedTotalCost
		}
	}
	return smallestRecord
}

func (p *AstarPathList) add(n *NodeRecord) {
	if n == nil {
		panic(fmt.Sprintf("list.add() trying to add nil NodeRecord"))
	}
	p.list[n.node.ID()] = n
}

func (p *AstarPathList) remove(n *NodeRecord) {
	delete(p.list, n.node.ID())
}

func (p *AstarPathList) len() int {
	return len(p.list)
}

func (p AstarPathList) contains(n Node) bool {
	if _, ok := p.list[n.ID()]; ok {
		return true
	}
	return false
}

func (p AstarPathList) find(n Node) *NodeRecord {
	return p.list[n.ID()]
}
