/*
*   LEGAL NOTICE
*   -------------------------------------------------------------------
*   This file contains modified code from https://github.com/beefsack/go-astar,
*   a project by Michael Charles Alexander (MIT License, 2014). Full license
*   details are provided at the bottom of this file.
*
 */

// astar is an A* pathfinding implementation.

package navigation

import (
	"container/heap"
)

// A priorityQueue implements heap.Interface and holds Nodes.  The
// priorityQueue is used to track open nodes by rank.
type priorityQueue []*node

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].rank < pq[j].rank
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	no := x.(*node)
	no.index = n
	*pq = append(*pq, no)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	no := old[n-1]
	no.index = -1
	*pq = old[0 : n-1]
	return no
}

// node is a wrapper to store A* data for a Pathable node.
type node struct {
	pather Pathable
	cost   float32
	rank   float32
	parent *node
	open   bool
	closed bool
	index  int
}

// nodeMap is a collection of nodes keyed by Pathable nodes for quick reference.
type nodeMap map[Pathable]*node

// get gets the Pathable object wrapped in a node, instantiating if required.
func (nm nodeMap) get(p Pathable) *node {
	n, ok := nm[p]
	if !ok {
		n = &node{
			pather: p,
		}
		nm[p] = n
	}
	return n
}

// astar_path calculates a short path and the distance between the two Pathable nodes.
//
// If no path is found, found will be false.
func astar_path(from, to Pathable) (path []Pathable, distance float32, found bool) {
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	fromNode := nm.get(from)
	fromNode.open = true
	heap.Push(nq, fromNode)
	for {
		if nq.Len() == 0 {
			// There's no path, return found false.
			return
		}
		current := heap.Pop(nq).(*node)
		current.open = false
		current.closed = true

		if current == nm.get(to) {
			// Found a path to the goal.
			p := []Pathable{}
			curr := current
			for curr != nil {
				p = append(p, curr.pather)
				curr = curr.parent
			}
			return p, current.cost, true
		}

		for _, neighbour := range current.pather.path_neighbours() {
			cost := current.cost + current.pather.path_neighbour_cost(neighbour)
			neighbourNode := nm.get(neighbour)
			if cost < neighbourNode.cost {
				if neighbourNode.open {
					heap.Remove(nq, neighbourNode.index)
				}
				neighbourNode.open = false
				neighbourNode.closed = false
			}
			if !neighbourNode.open && !neighbourNode.closed {
				neighbourNode.cost = cost
				neighbourNode.open = true
				neighbourNode.rank = cost + neighbour.path_estimated_cost(to)
				neighbourNode.parent = current
				heap.Push(nq, neighbourNode)
			}
		}
	}
}

/*
*   FULL LICENSING INFO
*   ----------------------------------------------------------
*   The MIT License (MIT)
*
*   Copyright (c) 2014 Michael Charles Alexander
*
*   Permission is hereby granted, free of charge, to any person obtaining a
*   copy of this software and associated documentation files (the "Software"),
*   to deal in the Software without restriction, including without limitation
*   the rights to use, copy, modify, merge, publish, distribute, sublicense,
*   and/or sell copies of the Software, and to permit persons to whom the
*   Software is furnished to do so, subject to the following conditions:
*
*   The above copyright notice and this permission notice shall be included in
*   all copies or substantial portions of the Software.
*
*   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
*   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
*   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
*   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
*   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
*   FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
*   DEALINGS IN THE SOFTWARE.
*   ----------------------------------------------------------
 */
