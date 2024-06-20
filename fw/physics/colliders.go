package physics

import (
	"gorl/fw/core/logging"

	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// -------------------
//  TYPES
// -------------------

type CollisionCallback func()
type BodyType uint8

const (
	BodyTypeStatic BodyType = iota
	BodyTypeKinematic
	BodyTypeDynamic
)

// -------------------
//
//	GENERIC COLLIDER
//
// -------------------
type Collider struct {
	body      *box2d.B2Body
	callbacks map[CollisionCategory]CollisionCallback
}

func newCollider(body *box2d.B2Body) *Collider {
	c := &Collider{
		body:      body,
		callbacks: make(map[CollisionCategory]CollisionCallback),
	}

	// link the collider to the body, so we can retrieve it on collisions
	c.body.SetUserData(c)

	return c
}

// Returns a pointer to the Box2D body of the collider.
func (c *Collider) GetB2Body() *box2d.B2Body {
	return c.body
}

func (c *Collider) GetCallbacks() map[CollisionCategory]CollisionCallback {
	return c.callbacks
}

// GetPosition returns the current position of the given collider.
func (c *Collider) GetPosition() rl.Vector2 {
	b2v := c.GetB2Body().GetPosition()
	v := rl.NewVector2(float32(b2v.X), float32(b2v.Y))
	v = simulationToPixelScaleV(v)
	return v
}

// SetPosition sets the position of the given collider. This may cause
// unexpected behavior if the collider is a dynamic body.
func (c *Collider) SetPosition(position rl.Vector2) {
	position = pixelToSimulationScaleV(position)
	c.GetB2Body().SetTransform(box2d.MakeB2Vec2(float64(position.X), float64(position.Y)), c.GetB2Body().GetAngle())
}

// GetVelocity returns the current velocity of the given collider.
func (c *Collider) GetVelocity() rl.Vector2 {
	b2v := c.GetB2Body().GetLinearVelocity()
	v := rl.NewVector2(float32(b2v.X), float32(b2v.Y))
	return v
}

// GetVertices returns the vertices of the collider. This is only valid for
// polygon or chain shape colliders.
func (c *Collider) GetVertices() []rl.Vector2 {
	vertices := make([]rl.Vector2, 0)
	for f := c.GetB2Body().GetFixtureList(); f != nil; f = f.GetNext() {
		switch f.GetShape().GetType() {
		case box2d.B2Shape_Type.E_polygon:
			shape := f.GetShape().(*box2d.B2PolygonShape)
			for i := 0; i < shape.M_count; i++ {
				v := shape.M_vertices[i]
				rlVec := simulationToPixelScaleV(rl.NewVector2(float32(v.X), float32(v.Y)))
				vertices = append(vertices, rlVec)
			}
		case box2d.B2Shape_Type.E_chain:
			shape := f.GetShape().(*box2d.B2ChainShape)
			for i := 0; i < shape.M_count; i++ {
				v := shape.M_vertices[i]
				rlVec := simulationToPixelScaleV(rl.NewVector2(float32(v.X), float32(v.Y)))
				vertices = append(vertices, rlVec)
			}
		}
	}
	return vertices
}

// ------------------------
//  Config Chain Functions
// ------------------------

// Set the density of the collider. The default density is 1.0
func (col *Collider) SetDensity(density float32) *Collider {
	col.SetDensity(density)
	return col
}

// Set the linear damping of the collider. The default damping is 1.0
func (col *Collider) SetLinearDamping(damping float32) *Collider {
	col.GetB2Body().SetLinearDamping(float64(damping))
	return col
}

// Set the category of the collider. The default category is CollisionCategoryAll.
func (col *Collider) SetCategory(category CollisionCategory) *Collider {
	col.GetB2Body().GetFixtureList().M_filter.CategoryBits = uint16(category)
	return col
}

// Set the collision mask of the collider. The default mask is CollisionCategoryAll.
func (col *Collider) SetMask(mask CollisionCategory) *Collider {
	col.GetB2Body().GetFixtureList().M_filter.MaskBits = uint16(mask)
	return col
}

// Set the callbacks of the collider. The default callbacks are none.
func (col *Collider) SetCallbacks(callbacks map[CollisionCategory]CollisionCallback) *Collider {
	if callbacks == nil {
		logging.Warning("Attempted to set callbacks map to nil. Refusing.")
		return col
	}
	col.callbacks = callbacks
	return col
}

// Set the fixed rotation flag of the collider. The default is false.
func (col *Collider) SetFixedRotation(fixed_rotation bool) *Collider {
	col.GetB2Body().SetFixedRotation(fixed_rotation)
	return col
}

// Set the is bullet flag of the collider. The default is false.
func (col *Collider) SetIsBullet(is_bullet bool) *Collider {
	col.GetB2Body().SetBullet(is_bullet)
	return col
}

// Set the is sensor flag of the collider. The default is false.
func (col *Collider) SetIsSensor(is_sensor bool) *Collider {
	col.GetB2Body().GetFixtureList().M_isSensor = is_sensor
	return col
}

// -------------------
//  Generic Functions
// -------------------

// DestroyCollider removes the given collider from the physics world.
func DestroyCollider(collider *Collider) {
	State.destructionQueue = append(State.destructionQueue, collider.GetB2Body())
}

// ---------------------
// CHAIN SHAPE COLLIDER
// ---------------------
//
// NewChainShapeCollider creates a new chain shape collider given its vertices.
func NewChainShapeCollider(
	position rl.Vector2,
	vertices []rl.Vector2,
	body_type BodyType,
) *Collider {
	// Convert the given position and vertices to the simulation scale
	position = pixelToSimulationScaleV(position)
	for i := range vertices {
		vertices[i] = pixelToSimulationScaleV(vertices[i])
	}

	// body definition
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(float64(position.X), float64(position.Y))
	bd.Type = uint8(body_type)
	bd.FixedRotation = false
	bd.AllowSleep = true
	bd.LinearDamping = float64(1.0)

	// shape
	shape := box2d.MakeB2ChainShape()
	b2vertices := make([]box2d.B2Vec2, len(vertices))
	for i, v := range vertices {
		b2vertices[i] = box2d.MakeB2Vec2(float64(v.X), float64(v.Y))
	}
	shape.CreateChain(b2vertices, len(vertices))

	// fixture definition
	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = float64(1.0)
	fd.Filter.CategoryBits = uint16(CollisionCategoryAll)

	// creating the body
	body := State.physicsWorld.CreateBody(&bd)
	body.CreateFixtureFromDef(&fd)

	// create the wrapping collider
	c := newCollider(body)

	return c
}

// -------------------
//
//	CIRCLE COLLIDER
//
// -------------------
//
// NewCircleCollider creates a new circle collider at the given position with
// the given radius, density and static flag.
func NewCircleCollider(
	position rl.Vector2,
	radius float32,
	body_type BodyType,
) *Collider {
	position = pixelToSimulationScaleV(position)
	radius = pixelToSimulationScale(radius)

	// body definition
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(float64(position.X), float64(position.Y))
	bd.Type = uint8(body_type)
	bd.FixedRotation = false
	bd.AllowSleep = true
	bd.LinearDamping = float64(1.0)

	// shape
	shape := box2d.MakeB2CircleShape()
	shape.M_radius = float64(radius)

	// fixture definition
	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = float64(1.0)
	fd.Filter.CategoryBits = uint16(CollisionCategoryAll)

	// creating the body
	body := State.physicsWorld.CreateBody(&bd)
	body.CreateFixtureFromDef(&fd)

	// create the wrapping collider
	c := newCollider(body)

	return c
}

// -------------------
//
//	CONVEX COLLIDER
//
// -------------------
//
// NewConvexCollider creates a new convex collider given its vertices (in clockwise order).
// Vertices are specified in pixels and will be converted to the simulation scale.
//
// The vertices are assumed to be in relative space to the given position.
func NewConvexCollider(
	position rl.Vector2,
	vertices []rl.Vector2,
	body_type BodyType,
) *Collider {
	// Convert the given position and vertices to the simulation scale
	position = pixelToSimulationScaleV(position)
	for i := range vertices {
		vertices[i] = pixelToSimulationScaleV(vertices[i])
	}

	b2vertices := make([]box2d.B2Vec2, len(vertices))
	for i, v := range vertices {
		b2vertices[i] = box2d.MakeB2Vec2(float64(v.X), float64(v.Y))
	}

	// Create the body definition
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(float64(position.X), float64(position.Y))
	bd.Type = uint8(body_type)
	bd.FixedRotation = false
	bd.AllowSleep = true
	bd.LinearDamping = float64(1.0)

	// Create the shape
	shape := box2d.MakeB2PolygonShape()
	shape.Set(b2vertices, len(vertices))

	// Fixture definition
	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = float64(1.0)
	fd.Filter.CategoryBits = uint16(CollisionCategoryAll)

	// Create the body in the world
	body := State.physicsWorld.CreateBody(&bd)
	body.CreateFixtureFromDef(&fd)

	// create the wrapping collider
	c := newCollider(body)

	return c
}

// Like NewConvexCollider, but takes in the vertices as absolute coordinates in
// world space, rather than position relative. This is helpful for definining
// static game-map colliders for example.
func NewConvexColliderAbs(
	vertices []rl.Vector2,
	body_type BodyType,
) *Collider {
	// Convert the given vertices to the simulation scale
	for i := range vertices {
		vertices[i] = pixelToSimulationScaleV(vertices[i])
	}

	b2vertices := make([]box2d.B2Vec2, len(vertices))
	for i, v := range vertices {
		b2vertices[i] = box2d.MakeB2Vec2(float64(v.X), float64(v.Y))
	}

	// Create the body definition
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(0, 0) // the position is (0, 0) since were getting absolute world coordinates
	bd.Type = uint8(body_type)
	bd.FixedRotation = false
	bd.AllowSleep = true
	bd.LinearDamping = float64(1.0)

	// Create the shape
	shape := box2d.MakeB2PolygonShape()
	shape.Set(b2vertices, len(vertices))

	// Fixture definition
	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = float64(1.0)
	fd.Filter.CategoryBits = uint16(CollisionCategoryAll)

	// Create the body in the world
	body := State.physicsWorld.CreateBody(&bd)
	body.CreateFixtureFromDef(&fd)

	// create the wrapping collider
	c := newCollider(body)

	return c
}
