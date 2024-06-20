package navigation

/*  NAVIGATION
*   --------------------------------------------------------
*   The navigation package consists of three main parts:
*   - The `Pathable` interface
*   - `Pathable` implementations
*   - Solvers
*
*   USAGE
*   --------------------------------------------------------
*   1. Pick a `Pathable` implementation, or create your own.
*   2. Set up an environment to path through with that implementation.
*   3. Run FindPath() on two Pathables to obtain a Path with the chosen method.
 */

import (
	"gorl/fw/core/logging"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Pathable interface {
	path_neighbours() []Pathable
	path_neighbour_cost(to Pathable) float32
	path_estimated_cost(to Pathable) float32
	GetPosition() rl.Vector2
}

type NavigationMethod int32

const (
	NavigationMethodAstar NavigationMethod = iota
)

// FindPath finds a path between the start and end Pathables using the specified method.
// If no path is found, the found value will be false.
func FindPath(start Pathable, end Pathable, method NavigationMethod) (path []Pathable, distance float32, found bool) {
	if start == nil || end == nil {
		logging.Error("Either start: %v or end: %v (or both) are nil! Unable to find path.", start, end)
		return nil, 0, false
	}

	switch method {
	case NavigationMethodAstar:
		return astar_path(start, end)
	}

	logging.Error("Specified an unknown Pathfinding method, failed to find path! Method: %v", method)
	return nil, 0, false
}
