package astar

import (
	"gorl/fw/core/datastructures"
	"gorl/fw/core/math"
	"slices"
)

// http://theory.stanford.edu/~amitp/GameProgramming/Heuristics.html#manhattan-distance
func manhattanHeuristic(a, b math.Vector2Int) float64 {
	delta := math.Vector2IntSub(a, b)
	D := 1.0
	return D * (float64(delta.X) + float64(delta.Y))
}

// http://theory.stanford.edu/~amitp/GameProgramming/ImplementationNotes.html
// and
// https://www.redblobgames.com/pathfinding/a-star/introduction.html
//
// A* algorithm, obstacles are marked as true in the grid.
func AstarPath(start, goal math.Vector2Int, grid [][]bool) []math.Vector2Int {

	// check that both start and goal are within the grid and are not obstacles

	frontier := datastructures.NewMinPriorityQueue[math.Vector2Int, float64]()
	frontier.Push(start, 0)

	cameFrom := make(map[math.Vector2Int]math.Vector2Int)
	costSoFar := make(map[math.Vector2Int]float64)

	cameFrom[start] = start
	costSoFar[start] = 0

	for !frontier.Empty() {
		current, _, _ := frontier.Pop()

		// stop when we reach the goal
		if current == goal {
			break
		}

		neighbors := neighbors(current, grid)
		for _, next := range neighbors {
			newCost := costSoFar[current] + 1
			if _, ok := costSoFar[next]; !ok || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				priority := newCost + 0 //  manhattanHeuristic(goal, next)
				frontier.Push(next, priority)
				cameFrom[next] = current
			}
		}
	}

	// reconstruct the path from goal to start
	path := make([]math.Vector2Int, 0)
	current := goal
	for current != start {
		path = append(path, current)
		current = cameFrom[current]
	}

	// return the reversed the path
	slices.Reverse(path)
	//logging.Debug("A* path found: %v", path)
	return path
}

// pointIsValid checks if a point is within the grid and is not an obstacle.
func pointIsValid(point math.Vector2Int, grid [][]bool) bool {
	return point.X >= 0 &&
		point.X < len(grid) &&
		point.Y >= 0 &&
		point.Y < len(grid[0]) &&
		!grid[point.X][point.Y]
}

// neighbors returns the 4-connected neighbors of a point that are within the
// grid and are not obstacles.
func neighbors(current math.Vector2Int, grid [][]bool) []math.Vector2Int {
	// 4-connected neighbors
	neighbors := []math.Vector2Int{
		{X: current.X - 1, Y: current.Y},
		{X: current.X + 1, Y: current.Y},
		{X: current.X, Y: current.Y - 1},
		{X: current.X, Y: current.Y + 1},
	}

	// filter out neighbors that are outside the grid or are obstacles
	validNeighbors := make([]math.Vector2Int, 0, len(neighbors))
	for _, n := range neighbors {
		if pointIsValid(n, grid) {
			validNeighbors = append(validNeighbors, n)
		}
	}

	return validNeighbors
}
