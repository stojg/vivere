package main

import (
	"fmt"
	"math"
)

type NodeRecord struct {
	node       Node
	connection *NodeRecord
	costSoFar  float64
}

func Dijkstra(graph *Graph, start, goal Node) []*NodeRecord {

	open := &PathFindingList{}
	closed := &PathFindingList{}

	open.add(&NodeRecord{
		node:       start,
		connection: nil,
		costSoFar:  0,
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
			// .. or if it is open and we've found a worse route

			if open.contains(toNode) {
				// here we find the record in the open list corresponding to the endNode
				record = open.find(toNode)
				// but the cost calculated for it is lower than the current one
				if record.costSoFar <= endNodeCost {
					continue
				}
			} else {
				// otherwise we know we've got an unvisited node node, so make a record for it
				record = &NodeRecord{
					node: toNode,
				}
			}

			// we are here if we need to update the node.  update the cost and connection
			record.costSoFar = endNodeCost
			record.connection = current

			// and add it to the open list
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
	// found which
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

type Node interface {
	ID() int
}

func NewGraph() *Graph {
	return &Graph{
		connections: make(map[int][]*Connection),
	}
}

type Graph struct {
	connections map[int][]*Connection
}

func (g *Graph) Add(from, to Node, cost float64) {
	g.connections[from.ID()] = append(g.connections[from.ID()], &Connection{
		From: from,
		To:   to,
		Cost: cost,
	})
}

func (g *Graph) getConnections(node Node) []*Connection {
	return g.connections[node.ID()]
}

type Connection struct {
	From Node
	To   Node
	Cost float64
}

// todo, make it more performant
type PathFindingList struct {
	list []*NodeRecord
}

// returns the NodeRecord with the smallest cost so far
func (p *PathFindingList) smallest() *NodeRecord {
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

func (p *PathFindingList) add(n *NodeRecord) {
	if n == nil {
		panic(fmt.Sprintf("list.add() trying to add nil NodeRecord"))
	}
	p.list = append(p.list, n)
}

func (p *PathFindingList) remove(n *NodeRecord) {
	if len(p.list) == 0 {
		return
	}
	if n.node == nil {
		panic("list.remove() trying to remove NodeRecord without a Node")
	}
	for i, record := range p.list {
		if record == nil {
			panic(fmt.Sprintf("list.remove() - found nil NodeRecord in list: %v", p.list))
		}
		if record.node == n.node {
			p.list = append(p.list[:i], p.list[i+1:]...)
		}
	}
}

func (p *PathFindingList) len() int {
	return len(p.list)
}

func (p PathFindingList) contains(n Node) bool {
	for _, record := range p.list {
		if record.node.ID() == n.ID() {
			return true
		}
	}
	return false
}

func (p PathFindingList) find(n Node) *NodeRecord {
	for _, record := range p.list {
		if record.node.ID() == n.ID() {
			return record
		}
	}
	return nil
}
