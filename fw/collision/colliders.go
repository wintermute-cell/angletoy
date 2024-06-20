package collision

import (
	"gorl/fw/core/logging"
	"gorl/fw/util"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv"
)

// -------------------------
//
//	CONVEX POLYGON COLLIDER
//
// -------------------------
// This checks at compile time if the interface is implemented
var _ Collider = (*ConvexCollider)(nil)

// RectCollider is a rectangular collider
type ConvexCollider struct {
	BaseCollider
    position *rl.Vector2
    points []rl.Vector2
	rotation     *float32
	resolv_shape resolv.IShape
    tags []string
}

// NewConvexCollider creates a new rect collider. The position will be taken as
// the rotation origin, all the other points are considered relative to the
// position. The position itself is not a point on the polygon.
//
// This example will create a 10x10 box at position (100, 100):
// NewConvexCollider(rl.NewVector2(100, 100), []rl.Vector2{{X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10}}, other arguments...)
//
func NewConvexCollider(
    position *rl.Vector2,
    points_clockwise []rl.Vector2,
    rotation *float32,
    tags []string,
    is_static bool,
) *ConvexCollider {
    if len(points_clockwise) < 2 {
        logging.Fatal("Failed to create Convex Polygon collider, less than 2 points were given!")
        return nil
    }
	c := &ConvexCollider{
		BaseCollider: BaseCollider{
			is_static: is_static,
		},
        position: position,
        points: points_clockwise,
		rotation:     rotation,
        tags: tags,
	}

    c.update_resolv_shape()
    logging.Debug("%v", c.resolv_shape)
	return c
}

// Like NewConvexCollider, but the passed points are not relative to the
// given position, but rather absolute world positions.
func NewConvexColliderAbs(
    position *rl.Vector2,
    points_clockwise []rl.Vector2,
    rotation *float32,
    tags []string,
    is_static bool,
) *ConvexCollider {
    if len(points_clockwise) < 2 {
        logging.Fatal("Failed to create Convex Polygon collider, less than 2 points were given!")
        return nil
    }
    points_clockwise_relative := []rl.Vector2{}
    for _, point := range points_clockwise {
        points_clockwise_relative = append(points_clockwise_relative, rl.Vector2Subtract(point, *position))
    }
	c := &ConvexCollider{
		BaseCollider: BaseCollider{
			is_static: is_static,
		},
        position: position,
        points: points_clockwise_relative,
		rotation:     rotation,
        tags: tags,
	}

    rotated_points := []rl.Vector2{}
    for _, point := range c.points {
        rotated_points = append(rotated_points, util.RotatePointAroundOrigin(point, *c.position, *c.rotation))
    }

    // transfer the raylib vectors into a slice of float64s which resolv understands
    points_as_one_list := []float64{}
    for _, point := range rotated_points {
        points_as_one_list = append(points_as_one_list, float64(point.X))
        points_as_one_list = append(points_as_one_list, float64(point.Y))
    }

    c.resolv_shape = resolv.NewConvexPolygon(
        float64(c.position.X), float64(c.position.Y),
        points_as_one_list...,
        )
	return c
}

// GetPosition returns the center of the rect collider
func (c *ConvexCollider) GetPosition() rl.Vector2 {
	return rl.NewVector2(c.position.X, c.position.Y)
}

// GetBounds returns the bounding box of the rect collider
func (c *ConvexCollider) GetBounds() rl.Rectangle {
	return calculate_convex_poly_bounds(*c.position, c.points)
}

func calculate_convex_poly_bounds(position rl.Vector2, points []rl.Vector2) rl.Rectangle {
	if len(points) == 0 {
		// Return a zero Rectangle if no points are provided.
		return rl.NewRectangle(0, 0, 0, 0)
	}
	
	// Initialize min and max coordinates with the first point.
	firstPoint := rl.Vector2Add(position, points[0])
	minX, minY, maxX, maxY := firstPoint.X, firstPoint.Y, firstPoint.X, firstPoint.Y
	
	// Iterate over all points, updating the min and max coordinates.
	for _, point := range points {
		transformedPoint := rl.Vector2Add(position, point)
		if transformedPoint.X < minX {
			minX = transformedPoint.X
		}
		if transformedPoint.Y < minY {
			minY = transformedPoint.Y
		}
		if transformedPoint.X > maxX {
			maxX = transformedPoint.X
		}
		if transformedPoint.Y > maxY {
			maxY = transformedPoint.Y
		}
	}
	
	// Create and return the bounding box as an rl.Rectangle.
	return rl.NewRectangle(minX, minY, maxX-minX, maxY-minY)
}

// GetRotation returns the rotation of the collider in degrees
func (c *ConvexCollider) GetRotation() float32 {
	return *c.rotation
}

// GetOrigin returns the origin of the collider
func (c *ConvexCollider) GetOrigin() rl.Vector2 {
	return *c.position
}

// GetType returns the type of the collider
func (c *ConvexCollider) GetType() ColliderType {
	return ColliderTypeConvex
}

// GetTags returns the tags of the collider
func (c *ConvexCollider) GetTags() []string {
    return c.tags
}

// GetResolvShape returns the resolv shape of the collider
func (c *ConvexCollider) GetResolvShape() resolv.IShape {
	return c.resolv_shape
}

// IsStatic returns whether the collider is static or not
func (c *ConvexCollider) IsStatic() bool {
	return c.is_static
}

func (c *ConvexCollider) DrawShape() {
    p := *c.position
    for i, point := range c.points {
        abs_point := rl.Vector2Add(point, p)
        if i < len(c.points) - 1 {
            rl.DrawLineV(abs_point, rl.Vector2Add(c.points[i+1], p), rl.Color{R: 255, G: 0, B: 0, A: 100})
        } else {
            rl.DrawLineV(abs_point, rl.Vector2Add(c.points[0], p), rl.Color{R: 255, G: 0, B: 0, A: 100})
        }
    }
}

// TODO: calling this function tanks performance
func (c *ConvexCollider) update_resolv_shape() {
    // rotate all points around the colliders position
    //rotated_points := []rl.Vector2{}
    //for _, point := range c.points {
    //    rotated_points = append(rotated_points, util.RotatePointAroundOrigin(point, *c.position, *c.rotation))
    //}

    //// transfer the raylib vectors into a slice of float64s which resolv understands
    //points_as_one_list := []float64{}
    //for _, point := range rotated_points {
    //    points_as_one_list = append(points_as_one_list, float64(point.X))
    //    points_as_one_list = append(points_as_one_list, float64(point.Y))
    //}

    //c.resolv_shape = resolv.NewConvexPolygon(
    //    float64(c.position.X), float64(c.position.Y),
    //    points_as_one_list...,
    //    )
}

// -----------------
//
//	RECT COLLIDER
//
// -----------------
// This checks at compile time if the interface is implemented
var _ Collider = (*RectCollider)(nil)

// RectCollider is a rectangular collider
type RectCollider struct {
	BaseCollider
	bounds       *rl.Rectangle
	rotation     *float32
	origin       *rl.Vector2
	resolv_shape resolv.IShape
    tags []string
}

// NewRectCollider creates a new rect collider
func NewRectCollider(bounds *rl.Rectangle, rotation *float32, origin *rl.Vector2, tags []string, is_static bool) *RectCollider {
	c := &RectCollider{
		BaseCollider: BaseCollider{
			is_static: is_static,
		},
		bounds:       bounds,
		rotation:     rotation,
		origin:       origin,
        tags: tags,
		resolv_shape: resolv.NewRectangle(float64(bounds.X), float64(bounds.Y), float64(bounds.Width), float64(bounds.Height)),
	}
	c.resolv_shape.SetRotation(-float64(*c.rotation * rl.Deg2rad))
	return c
}

// GetPosition returns the center of the rect collider
func (c *RectCollider) GetPosition() rl.Vector2 {
	return rl.NewVector2(c.bounds.X, c.bounds.Y)
}

// GetBounds returns the bounding box of the rect collider
func (c *RectCollider) GetBounds() rl.Rectangle {
	return *c.bounds
}

// GetRotation returns the rotation of the collider in degrees
func (c *RectCollider) GetRotation() float32 {
	return *c.rotation
}

// GetOrigin returns the origin of the collider
func (c *RectCollider) GetOrigin() rl.Vector2 {
	return *c.origin
}

// GetType returns the type of the collider
func (c *RectCollider) GetType() ColliderType {
	return ColliderTypeRect
}

// GetTags returns the tags of the collider
func (c *RectCollider) GetTags() []string {
    return c.tags
}

// GetResolvShape returns the resolv shape of the collider
func (c *RectCollider) GetResolvShape() resolv.IShape {
	return c.resolv_shape
}

// IsStatic returns whether the collider is static or not
func (c *RectCollider) IsStatic() bool {
	return c.is_static
}

func (c *RectCollider) DrawShape() {
    rl.DrawRectanglePro(
        *c.bounds,
        rl.Vector2Zero(),
        float32(*c.rotation),
        rl.Color{R: 255, G: 0, B: 0, A: 100},
    )
}

func (c *RectCollider) update_resolv_shape() {
	c.resolv_shape = resolv.NewRectangle(float64(c.bounds.X), float64(c.bounds.Y), float64(c.bounds.Width), float64(c.bounds.Height))
	c.resolv_shape.SetRotation(-float64(*c.rotation * rl.Deg2rad))
}

// -----------------
//  CIRCLE COLLIDER
// -----------------

// This checks at compile time if the interface is implemented
var _ Collider = (*CircleCollider)(nil)

// CircleCollider is a circular collider
type CircleCollider struct {
	BaseCollider
	center       *rl.Vector2
	radius       float32
    tags []string
	resolv_shape resolv.IShape
}

// NewCircleCollider creates a new circle collider
func NewCircleCollider(center *rl.Vector2, radius float32, tags []string, is_static bool) *CircleCollider {
	c := &CircleCollider{
		center:       center,
		radius:       radius,
        tags: tags,
		resolv_shape: resolv.NewCircle(float64(center.X), float64(center.Y), float64(radius)),
	}
	c.is_static = is_static
	return c
}

// GetPosition returns the center of the circle collider
func (c *CircleCollider) GetPosition() rl.Vector2 {
	return *c.center
}

// GetBounds returns the bounding box of the circle collider
func (c *CircleCollider) GetBounds() rl.Rectangle {
	return rl.NewRectangle(
		c.center.X-c.radius,
		c.center.Y-c.radius,
		c.radius*2, c.radius*2)
}

// GetRotation returns the rotation of the collider in degrees
func (c *CircleCollider) GetRotation() float32 {
	return 0
}

// GetOrigin returns the origin of the collider
func (c *CircleCollider) GetOrigin() rl.Vector2 {
	return rl.Vector2Zero()
}

// GetType returns the type of the collider
func (c *CircleCollider) GetType() ColliderType {
	return ColliderTypeCircle
}

func (c *CircleCollider) GetTags() []string {
    return c.tags
}

func (c *CircleCollider) GetResolvShape() resolv.IShape {
	return c.resolv_shape
}

// IsStatic returns whether the collider is static or not
func (c *CircleCollider) IsStatic() bool {
	return c.is_static
}

func (c *CircleCollider) DrawShape() {
    rl.DrawCircleV(*c.center, c.radius, rl.Color{R: 255, G: 0, B: 0, A: 100})
}

func (c *CircleCollider) update_resolv_shape() {
	c.resolv_shape.SetPosition(float64(c.center.X), float64(c.center.Y))
}
