package physics

import (
	"math"
	"strconv"
	"strings"

	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// CollidersFromString parses a string of the form
// "[{x1 y1}{x2 y2}...][{x1 y1}{x2 y2}...]" into a slice of convex colliders.
func CollidersFromString(collider_string string, category CollisionCategory, callbacks map[CollisionCategory]CollisionCallback) []*Collider {
    colliders := []*Collider{}
    for _, poly := range parsePolygonString(collider_string) {
        colliders = append(colliders, 
            NewConvexColliderAbs(poly, BodyTypeStatic).
                SetCategory(category).
                SetCallbacks(callbacks).
                SetFixedRotation(true))
    }
    return colliders
}

func parsePolygonString(polyStr string) [][]rl.Vector2 {
	var polygons [][]rl.Vector2
	// Extract individual polygons
	polyStrs := strings.Split(polyStr, "][")
	for _, polyStr := range polyStrs {
		polyStr = strings.Trim(polyStr, "[] ")
		pts := strings.Split(polyStr, "}")
		var polygon []rl.Vector2
		for _, ptStr := range pts {
            if ptStr == "" {
                continue
            }
			xy := strings.Fields(strings.Trim(ptStr, "{} "))
			x, _ := strconv.Atoi(xy[0])
			y, _ := strconv.Atoi(xy[1])
			polygon = append(polygon, rl.NewVector2(float32(x), float32(y)))
		}
		polygons = append(polygons, polygon)
	}

	return polygons
}

func createInternalProbeCallback(results *[]*box2d.B2Fixture) box2d.B2BroadPhaseQueryCallback {
    return func(fixture *box2d.B2Fixture) bool { 
        *results = append(*results, fixture)
        return true // I don't know why we need to do this, or if this is correct.
    }
}

// ProbePoint checks if the given point intersects with any colliders, taking
// into account only the given collision categories. 
func ProbePoint(point rl.Vector2, categoriesToCheck CollisionCategory) []*Collider {
    point = pixelToSimulationScaleV(point)
    if categoriesToCheck == 0 {
        categoriesToCheck = math.MaxUint16 // set the bitmask to only 1s
    }

    // first, we accumulate all the fixtures where the point is inside the AABB
    aabb := box2d.MakeB2AABB()
    aabb.LowerBound.Set(float64(point.X), float64(point.Y))
    aabb.UpperBound.Set(float64(point.X), float64(point.Y))
    intersecting_fixtures := []*box2d.B2Fixture{}
    State.physicsWorld.QueryAABB(createInternalProbeCallback(&intersecting_fixtures), aabb)

    // then, we check if the point is actually inside any of the fixtures
    results := []*Collider{}
    for _, fix := range intersecting_fixtures {
        // if the fixture does not match any of the categories we're checking...
        if fix.GetFilterData().CategoryBits & uint16(categoriesToCheck) == 0 {
            continue
        }
        if fix.TestPoint(box2d.MakeB2Vec2(float64(point.X), float64(point.Y))) {
            results = append(results, fix.GetBody().GetUserData().(*Collider))
        }
    }

    return results
}

// GenerateMatrixMap generates a 2D slice indicating if the point at a certain
// coordinate intersects a collider.
// `resolution` divides the pixels determined by area.
// With an area of 10x10, a resolution of 1 will result in a 10x10 map, while
// a resolution of 2 will result in a 5x5 map.
//
// If map[x][y] == true, then this point intersects a collider.
//
// Note:
// When generating a map with a resolution of R, the point retrieved with
// map[x][y] corresponds to the real position (x*R, y*R).
func GenerateMatrixMap(area rl.Rectangle, resolution int32, categoriesToCheck CollisionCategory) [][]bool {
	width := int(area.Width) / int(resolution)
	height := int(area.Height) / int(resolution)

	navmap := make([][]bool, width)
	for i := range navmap {
		navmap[i] = make([]bool, height)
	}

	tileSize := float32(area.Width) / float32(width)

    // adjust probe_size to possibly fix tiles falsely marked as blocking.
    // careful though, values below 1 can quickly cause holes in the navmap.
    probe_size := 0.98
	circle_radius := (tileSize / 2)*float32(probe_size)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			P := rl.NewVector2(float32(i)*tileSize + area.X + circle_radius, float32(j)*tileSize + area.Y + circle_radius)
            hit_colliders := ProbePoint(P, categoriesToCheck)
            if len(hit_colliders) > 0 {
                navmap[i][j] = true
            }
		}
	}

	return navmap
}
