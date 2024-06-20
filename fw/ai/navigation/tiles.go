package navigation

/*  TILES
*   --------------------------------------------------------
*   Tiles is a basic imlementation of the Pathable interface.
*   It should suffice for most tasks.
*   The Structure consists of a fixed size 2D grid World, with each tile
*   of the grid being a navigatable node with a certain traversal cost.
*
*   USAGE
*   --------------------------------------------------------
*   1. Create a World using NewPathableWorld()
*   2. Set custom tile costs using world.SetCost() (optional)
*   3. Retrieve a pair of tiles t1, t2 using world.GetTile()
*   4. Find a path using FindPath()
*
*   EXAMPLE
*   --------------------------------------------------------
*   pathable_world = navigation.NewPathableWorld(rl.NewRectangle(0, 0, 10, 10))
*   pathable_world.SetCost(rl.NewVector2(3, 2), -1)
*   pathable_world.SetCost(rl.NewVector2(4, 2), -1)
*   pathable_world.SetCost(rl.NewVector2(3, 3), -1)
*   start_tile = pathable_world.GetTile(rl.Vector2Zero())
*   end_tile = pathable_world.GetTile(rl.NewVector2(5, 5))
*   path, dist, found := navigation.FindPath(start_tile, end_tile, navigation.NavigationMethodAstar)
*
 */

import (
	"gorl/fw/core/logging"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PathableWorld struct {
	Bounds     rl.Rectangle
	Tiles      [][]*PathableTile
	Resolution int32
}

func NewPathableWorld(bounds rl.Rectangle, resolution int32) *PathableWorld {
	world := &PathableWorld{Bounds: bounds, Resolution: resolution}

	width := int(bounds.Width) - int(bounds.X)
	height := int(bounds.Height) - int(bounds.Y)

	tiles := make([][]*PathableTile, height)
	for i := range tiles {
		tiles[i] = make([]*PathableTile, width)
		for j := range tiles[i] {
			tiles[i][j] = &PathableTile{
				Position: rl.Vector2{X: float32(j + int(bounds.X)), Y: float32(i + int(bounds.Y))},
				Cost:     0,
				World:    world,
			}
		}
	}

	world.Tiles = tiles

	return world
}

func (pw *PathableWorld) GetTile(position rl.Vector2) *PathableTile {
	scaled_pos := rl.NewVector2(position.X/float32(pw.Resolution), position.Y/float32(pw.Resolution))

	// Bounds check
	if scaled_pos.X < pw.Bounds.X ||
		scaled_pos.Y < pw.Bounds.Y ||
		scaled_pos.X > pw.Bounds.X+pw.Bounds.Width ||
		scaled_pos.Y > pw.Bounds.Y+pw.Bounds.Height {
        logging.Error("Called GetTile on position %v (real: %v) that is out of bounds for %v", scaled_pos, position, pw.Bounds)
		scaled_pos = rl.NewVector2(0, 0)
	}

	t := pw.get_tile_internal(scaled_pos)
	return t
}

func (pw *PathableWorld) get_tile_internal(position rl.Vector2) *PathableTile {
	// Convert world coordinates to array indices
	arrayX := int32(position.X) - int32(pw.Bounds.X)
	arrayY := int32(position.Y) - int32(pw.Bounds.Y)

	// Ensure the indices are within bounds
	if arrayX < 0 ||
		arrayX >= int32(len(pw.Tiles[0])) ||
		arrayY < 0 ||
		arrayY >= int32(len(pw.Tiles)) {
		return nil
	}

	return pw.Tiles[arrayY][arrayX]
}

// SetCost sets the cost of the tile at the given position.
func (pw *PathableWorld) SetCost(position rl.Vector2, cost float32) {
	scaled_pos := rl.NewVector2(position.X/float32(pw.Resolution), position.Y/float32(pw.Resolution))
	tile := pw.get_tile_internal(scaled_pos)
	if tile != nil {
		tile.Cost = cost
	} else {
		logging.Error("Tried to set the cost for a tile at an invalid position: %v", position)
	}
}

// PathableTile is a basic implementation of the Pathable interface.
// It should suffice for most tasks.
type PathableTile struct {
	Position rl.Vector2
	Cost     float32
	World    *PathableWorld
}

func (t *PathableTile) path_neighbours() []Pathable {
	neighbours := []Pathable{}
	for _, offset := range [][]float32{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		n := t.World.get_tile_internal(rl.NewVector2(t.Position.X+offset[0], t.Position.Y+offset[1]))
		// only if the requested neighbour is within bounds and has a positive cost
		if n != nil && n.Cost >= 0 {
			neighbours = append(neighbours, n)
		}
	}
	return neighbours
}

// path_neighbour_cost returns the movement cost of the directly neighbouring tile.
func (t *PathableTile) path_neighbour_cost(to Pathable) float32 {
	to_tile := to.(*PathableTile)
	return float32(to_tile.Cost)
}

// path_estimated_cost uses Manhattan distance to estimate orthogonal distance
// between non-adjacent nodes.
func (t *PathableTile) path_estimated_cost(to Pathable) float32 {
	to_tile := to.(*PathableTile)
	absX := to_tile.Position.X - t.Position.X
	if absX < 0 {
		absX = -absX
	}
	absY := to_tile.Position.Y - t.Position.Y
	if absY < 0 {
		absY = -absY
	}
	return float32(absX + absY)
}

func (t *PathableTile) GetPosition() rl.Vector2 {
	return t.Position
}
