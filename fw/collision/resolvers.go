package collision

import (
	"gorl/fw/core/logging"
	"gorl/fw/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type collision_resolver func(Collider, Collider) (rl.Vector2, rl.Vector2)

// --------------------
// RESOLVING FUNCTIONS
// --------------------

// ResolveCircleToCircle resolves a collision between two circles, returning
// a suggested resolution vector for each collider. (Taking into account any
// static colliders.)
func ResolveCircleToCircle(a Collider, b Collider) (rl.Vector2, rl.Vector2) {
	sum_of_radii := a.GetBounds().Width/2 + b.GetBounds().Width/2
	center_dist := rl.Vector2Distance(a.GetPosition(), b.GetPosition())

	overlap := sum_of_radii - center_dist
	if overlap <= 0 {
		return rl.Vector2Zero(), rl.Vector2Zero()
	}

	a_to_b := rl.Vector2Scale(util.Vector2NormalizeSafe(
		rl.Vector2Subtract(b.GetPosition(), a.GetPosition())), overlap)
	b_to_a := rl.Vector2Scale(a_to_b, -1.0)

	if !a.IsStatic() && !b.IsStatic() {
		// each collider moves half the way
		a_half := rl.Vector2Scale(b_to_a, 0.5) // a moves away from b (in the direction b->a)
		b_half := rl.Vector2Scale(a_to_b, 0.5) // b moves away from a (in the direction a->b)
		return a_half, b_half
	} else if !a.IsStatic() {
		// a must move, and will move in the direction b->a (away from b)
		return b_to_a, rl.Vector2Zero()
	} else if !b.IsStatic() {
		return rl.Vector2Zero(), a_to_b
	} else {
		logging.Warning("Tried to resolve faulty collision between two static colliders!")
		return rl.Vector2Zero(), rl.Vector2Zero()
	}
}

// DEPRECATED: Use ResolveGeneric instead.
// ResolveCircleToRectangle resolves a collision between a circle and a rectangle,
// returning a suggested resolution vector for each collider. (Taking into account
// any static colliders.)
func ResolveCircleToRectangle(a Collider, b Collider) (rl.Vector2, rl.Vector2) {
	self_inter := a.GetResolvShape().Intersection(0, 0, b.GetResolvShape())
	if self_inter == nil {
		return rl.Vector2Zero(), rl.Vector2Zero()
	}
	self_mtv := self_inter.MTV
	b_to_a := rl.NewVector2(float32(self_mtv.X()), float32(self_mtv.Y()))

	other_inter := b.GetResolvShape().Intersection(0, 0, a.GetResolvShape())
	other_mtv := other_inter.MTV
	a_to_b := rl.NewVector2(float32(other_mtv.X()), float32(other_mtv.Y()))

	if !a.IsStatic() && !b.IsStatic() {
		// each collider moves half the way
		a_half := rl.Vector2Scale(b_to_a, 0.5) // a moves away from b (in the direction b->a)
		b_half := rl.Vector2Scale(a_to_b, 0.5) // b moves away from a (in the direction a->b)
		return a_half, b_half
	} else if !a.IsStatic() {
		// a must move, and will move in the direction b->a (away from b)
		return b_to_a, rl.Vector2Zero()
	} else if !b.IsStatic() {
		return rl.Vector2Zero(), a_to_b
	} else {
		logging.Warning("Tried to resolve faulty collision between two static colliders!")
		return rl.Vector2Zero(), rl.Vector2Zero()
	}
}

// DEPRECATED: Use ResolveGeneric instead.
// ResolveRectangleToRectangle resolves a collision between two rectangles,
// returning a suggested resolution vector for each collider. (Taking into
// account any static colliders.)
func ResolveRectangleToRectangle(a Collider, b Collider) (rl.Vector2, rl.Vector2) {
	self_inter := a.GetResolvShape().Intersection(0, 0, b.GetResolvShape())
	if self_inter == nil {
		return rl.Vector2Zero(), rl.Vector2Zero()
	}

	self_mtv := self_inter.MTV
	b_to_a := rl.NewVector2(float32(self_mtv.X()), float32(self_mtv.Y()))

	other_inter := b.GetResolvShape().Intersection(0, 0, a.GetResolvShape())
	if other_inter == nil {
		return rl.Vector2Zero(), rl.Vector2Zero()
	}
	other_mtv := other_inter.MTV
	a_to_b := rl.NewVector2(float32(other_mtv.X()), float32(other_mtv.Y()))

	if !a.IsStatic() && !b.IsStatic() {
		// each collider moves half the way
		a_half := rl.Vector2Scale(b_to_a, 0.5) // a moves away from b (in the direction b->a)
		b_half := rl.Vector2Scale(a_to_b, 0.5) // b moves away from a (in the direction a->b)
		return a_half, b_half
	} else if !a.IsStatic() {
		// a must move, and will move in the direction b->a (away from b)
		return b_to_a, rl.Vector2Zero()
	} else if !b.IsStatic() {
		return rl.Vector2Zero(), a_to_b
	} else {
		logging.Warning("Tried to resolve faulty collision between two static colliders!")
		return rl.Vector2Zero(), rl.Vector2Zero()
	}
}

// ResolveGeneric resolves a collision between two colliders, returning a suggested
// resolution vector for each collider. (Taking into account any static colliders.)
func ResolveGeneric(a Collider, b Collider) (rl.Vector2, rl.Vector2) {
	self_inter := a.GetResolvShape().Intersection(0, 0, b.GetResolvShape())
	if self_inter == nil {
		return rl.Vector2Zero(), rl.Vector2Zero()
	}

	self_mtv := self_inter.MTV
	b_to_a := rl.NewVector2(float32(self_mtv.X()), float32(self_mtv.Y()))

	other_inter := b.GetResolvShape().Intersection(0, 0, a.GetResolvShape())
	if other_inter == nil {
		return rl.Vector2Zero(), rl.Vector2Zero()
	}
	other_mtv := other_inter.MTV
	a_to_b := rl.NewVector2(float32(other_mtv.X()), float32(other_mtv.Y()))

	if !a.IsStatic() && !b.IsStatic() {
		// each collider moves half the way
		a_half := rl.Vector2Scale(b_to_a, 0.5) // a moves away from b (in the direction b->a)
		b_half := rl.Vector2Scale(a_to_b, 0.5) // b moves away from a (in the direction a->b)
		return a_half, b_half
	} else if !a.IsStatic() {
		// a must move, and will move in the direction b->a (away from b)
		return b_to_a, rl.Vector2Zero()
	} else if !b.IsStatic() {
		return rl.Vector2Zero(), a_to_b
	} else {
		logging.Warning("Tried to resolve faulty collision between two static colliders!")
		return rl.Vector2Zero(), rl.Vector2Zero()
	}
}
