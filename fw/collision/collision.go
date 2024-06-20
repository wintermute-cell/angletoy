package collision

import (
	"gorl/fw/core/logging"
	"gorl/fw/util"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

type ColliderType int32

const (
	ColliderTypeCircle ColliderType = iota
	ColliderTypeRect
    ColliderTypeConvex
)

// Collider is an interface for objects that can collide
type Collider interface {
	GetPosition() rl.Vector2
	GetBounds() rl.Rectangle
	GetOrigin() rl.Vector2
	GetRotation() float32
	GetType() ColliderType
    GetTags() []string
	GetResolvShape() resolv.IShape
    DrawShape()
	update_resolv_shape()
	IsStatic() bool
}

type CollisionEvent struct {
	SelfType        ColliderType
	OtherType       ColliderType
	SelfCollider    Collider
	OtherCollider   Collider
	Resolution      rl.Vector2
	OtherResolution rl.Vector2
}

type BaseCollider struct {
	is_static bool
}

type CollisionCallback func(event CollisionEvent)

type collider_set struct {
	collider Collider
	callback CollisionCallback
}

type collision_state struct {
	collision_layers map[string][]*collider_set
    queued_for_destruction map[Collider][]string
}

var c collision_state

var collision_matrix map[ColliderType]map[ColliderType]collision_resolver

func InitCollision() {
	c = collision_state{
		collision_layers: make(map[string][]*collider_set),
        queued_for_destruction: make(map[Collider][]string),
	}

	collision_matrix = map[ColliderType]map[ColliderType]collision_resolver{
		ColliderTypeCircle: {
			ColliderTypeCircle: ResolveCircleToCircle,
			ColliderTypeRect:   ResolveCircleToRectangle,
            ColliderTypeConvex: ResolveGeneric,
		},
		ColliderTypeRect: {
			ColliderTypeRect: ResolveRectangleToRectangle, // TODO: this is quite fucky. tested using rect x rotated rect and it did not work well.
			ColliderTypeCircle:   ResolveCircleToRectangle,
            ColliderTypeConvex: ResolveGeneric,
		},
        ColliderTypeConvex: {
            ColliderTypeConvex: ResolveGeneric,
            ColliderTypeCircle: ResolveGeneric,
            ColliderTypeRect: ResolveGeneric,
        },
	}

	logging.Info("Initialized Collision System")
}

func DeinitCollision() {
}

func Update() {
    // Destroy every collider that is queued for destruction.
    // We use this queueing system to avoid modifying the collider list as the
    // Update loop us running through it.
    for col, layers := range c.queued_for_destruction {
        destroy_collider_internal(col, layers)
    }

	// Update the underlying resolve shapes to the colliders current values.
	for _, layer := range c.collision_layers {
		for _, col_set := range layer {
			col_set.collider.update_resolv_shape()
		}
	}

	// Test and resolve collisions
	for _, layer := range c.collision_layers {
		for i, col_set := range layer {
			// for every collider in this layer, loop over all other colliders in this layer
			for j := i + 1; j < len(layer); j++ {
				other_set := layer[j]
				a := col_set
				b := other_set

				// if both are static, nothing to do
				if a.collider.IsStatic() && b.collider.IsStatic() {
					continue
				}

                // if the colliders are further apart than their bounds, nothing to do
                a_range := util.Max(a.collider.GetBounds().Width, a.collider.GetBounds().Height)
                b_range := util.Max(b.collider.GetBounds().Width, b.collider.GetBounds().Height)
                if rl.Vector2Distance(a.collider.GetPosition(), b.collider.GetPosition()) > a_range + b_range {
                    continue
                }

				resolver_func, ok := collision_matrix[a.collider.GetType()][b.collider.GetType()]
				swapped := false
				if !ok {
					resolver_func, ok = collision_matrix[b.collider.GetType()][a.collider.GetType()]
					swapped = true
					if !ok {
						logging.Error("Tried to handle collision check, but did not find a resolver in the collision matrix. Combination: %v - %v", a.collider.GetType(), b.collider.GetType())
						return
					}
				}

				var a_res rl.Vector2
				var b_res rl.Vector2
				if !swapped {
					a_res, b_res = resolver_func(a.collider, b.collider)
				} else {
					b_res, a_res = resolver_func(b.collider, a.collider)
				}

				// if there is no resolution, we're not colliding
				if a_res == rl.Vector2Zero() && b_res == rl.Vector2Zero() {
					continue
				}

				a.callback(
					CollisionEvent{
						SelfType:        ColliderTypeCircle,
						OtherType:       ColliderTypeRect,
						SelfCollider:    a.collider,
						OtherCollider:   b.collider,
						Resolution:      a_res,
						OtherResolution: b_res,
					})
				b.callback(
					CollisionEvent{
						SelfType:        ColliderTypeRect,
						OtherType:       ColliderTypeCircle,
						SelfCollider:    a.collider,
						OtherCollider:   b.collider,
						Resolution:      b_res,
						OtherResolution: a_res,
					})
			}
		}
	}
}

// Draw all registered colliders. Useful for debugging purposes.
func DrawColliders() {
	for _, layer := range c.collision_layers {
		for _, col_set := range layer {
			col := col_set.collider
            col.DrawShape()
		}
	}
}

// RegisterCollider registers a collider to a layer. When a collision occurs,
// the callback function is called with the other collider as an argument.
// Only colliders on the same layer can collide.
func RegisterCollider(collider Collider, layers []string, callback CollisionCallback) Collider {
	set := &collider_set{collider: collider, callback: callback}
	for _, layer := range layers {
		c.collision_layers[layer] = append(
			c.collision_layers[layer],
			set)
	}
	return set.collider
}

// DestroyCollider removes a collider from the given layers.
func DestroyCollider(collider Collider, layers []string) {
    c.queued_for_destruction[collider] = append(c.queued_for_destruction[collider], layers...)
}

func destroy_collider_internal(collider Collider, layers []string) {
	for _, layer := range layers {
		for i, set := range c.collision_layers[layer] {
			if set.collider == collider {
				c.collision_layers[layer] = util.SliceDelete(c.collision_layers[layer], i, i+1)
				break
			}
		}
	}
}

// GenerateNavmap generates a navigation map for the given layers.
// `resolution` divides the pixels determined by area.
// With an area of 10x10, a resolution of 1 will result in a 10x10 map, while
// a resolution of 2 will result in a 5x5 map.
//
// Note:
// When generating a map with a resolution of R, the point retrieved with
// map[x][y] corresponds to the real position (x*R, y*R).
func GenerateNavmap(layers []string, area rl.Rectangle, resolution int32) [][]bool {
	width := int(area.Width) / int(resolution)
	height := int(area.Height) / int(resolution)

	navmap := make([][]bool, width)
	for i := range navmap {
		navmap[i] = make([]bool, height)
		for j := range navmap[i] {
			navmap[i][j] = true
		}
	}

	tileSize := float32(area.Width) / float32(width)

    // adjust probe_size to possibly fix tiles falsely marked as blocking.
    // careful though, values below 1 can quickly cause holes in the navmap.
    probe_size := 0.98
	circle_radius := (tileSize / 2)*float32(probe_size)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			P := rl.NewVector2(float32(i)*tileSize + area.X + circle_radius, float32(j)*tileSize + area.Y + circle_radius)

			for _, layer := range layers {
				for _, set := range c.collision_layers[layer] {
					collider := set.collider
					contact := collider.GetResolvShape().Intersection(0, 0, resolv.NewCircle(float64(P.X), float64(P.Y), float64(circle_radius)))
					if contact != nil {
						navmap[i][j] = false
					}
				}
			}
		}
	}

	return navmap
}

func CollidersFromString(collider_string string) []Collider {
    z := float32(0)
    colliders := []Collider{}
    for _, poly := range parsePolygonString(collider_string) {
        colliders = append(colliders, 
            NewConvexColliderAbs(
                &poly[0],
                poly,
                &z,
                []string{},
                true,
                ),
            )
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

