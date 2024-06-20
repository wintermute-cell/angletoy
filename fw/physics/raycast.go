package physics

import (
	"gorl/fw/core/logging"
	"gorl/fw/util"
	"sort"

	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RaycastHit struct {
	HitCollider        *Collider
	IntersectionPoint  rl.Vector2
    HitNormal rl.Vector2
}

func createInternalRaycastCallback(results *[]RaycastHit, filter CollisionCategory) box2d.B2RaycastCallback {
    return func(fixture *box2d.B2Fixture, point, normal box2d.B2Vec2, fraction float64) float64 {
        // check if the filter mask contains the category of the collider
        if (fixture.GetFilterData().CategoryBits & uint16(filter)) > 0 {
            // Create a RaycastHit for the hit fixture
            hit := RaycastHit{
                HitCollider: fixture.GetBody().GetUserData().(*Collider),
                IntersectionPoint: simulationToPixelScaleV(rl.Vector2{X: float32(point.X), Y: float32(point.Y)}),
                HitNormal: simulationToPixelScaleV(rl.Vector2{X: float32(normal.X), Y: float32(normal.Y)}),
            }
            *results = append(*results, hit)
            
            return 1  // continue the raycast to get all fixtures in its path
        } else {
            return -1
        }
    }
}


// Raycast casts a ray from origin to direction, returning a list of all
// colliders that were hit.
func Raycast(origin, direction rl.Vector2, length float32, categoriesToHit CollisionCategory) []RaycastHit {
    if length == 0 {
        logging.Warning("Attempted zero length raycast.")
        return []RaycastHit{}
    }
    // translate input values to simulation scale
    oOrigin := origin
    origin = pixelToSimulationScaleV(origin)
    length = pixelToSimulationScale(length)

    // calculate the endpoint
	normalized_direction := util.Vector2NormalizeSafe(direction)
    max_range := length
    endpoint := rl.Vector2Add(origin, rl.Vector2Scale(normalized_direction, max_range))

    // translate to box2d data types
    b2origin := box2d.MakeB2Vec2(float64(origin.X), float64(origin.Y))
    b2endpoint := box2d.MakeB2Vec2(float64(endpoint.X), float64(endpoint.Y))

    // do the raycast
    var results []RaycastHit
    callback := createInternalRaycastCallback(&results, categoriesToHit)
    State.physicsWorld.RayCast(callback, b2origin, b2endpoint)

    // sort results by distance to origin
    sort.Slice(results, func(i, j int) bool {
        return rl.Vector2Distance(oOrigin, results[i].IntersectionPoint) < rl.Vector2Distance(oOrigin, results[j].IntersectionPoint)
    })
    
    return results
}
