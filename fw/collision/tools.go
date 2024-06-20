package collision

import (
	"gorl/fw/util"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

type RaycastHit struct {
	HitCollider        Collider
	IntersectionPoints []rl.Vector2
}

// Raycast casts a ray from origin to direction, returning a list of all
// colliders that were hit. The ignore_list is a list of colliders that will be
// ignored in the cast.
//
// Note:
// A max_range of 0 will result in Infinite length.
func Raycast(layer_name string, origin, direction rl.Vector2, ignore_list []Collider, max_range float32) []*RaycastHit {
	dir := util.Vector2NormalizeSafe(direction)
	ret := []*RaycastHit{}
    if max_range == 0 {
        max_range = 10000000
    }
	for _, col_set := range c.collision_layers[layer_name] {
		col := col_set.collider
		if util.SliceContains(ignore_list, col) {
			continue
		}

        endpoint := rl.Vector2Add(origin, rl.Vector2Scale(dir, max_range))
        ray := resolv.NewLine(float64(origin.X), float64(origin.Y), float64(endpoint.X), float64(endpoint.Y))
        contact := col.GetResolvShape().Intersection(0, 0, ray)

        if contact != nil {
            points := []rl.Vector2{}
            for _, p := range contact.Points {
                points = append(points, rl.NewVector2(float32(p.X()), float32(p.Y())))
            }
            hit := &RaycastHit{
                HitCollider: col,
                IntersectionPoints: points,
            }
            ret = append(ret, hit)
        }
	}
	return ret
}
