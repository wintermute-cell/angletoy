package physics

import "math"

// CollisionCategory is a Bitmask.
//
// To combine multiple categories, use the or operator:
// CollisionCategoryOneAndTwo := CollisionCategoryOne | CollisionCategoryTwo
type CollisionCategory uint16
const (
    CollisionCategoryPlayer CollisionCategory = 1 << iota
    CollisionCategoryEnemy
    CollisionCategoryEnvironment
    CollisionCategoryBullet
)

const (
    // No collision with anything. All bits are set to 0.
    CollisionCategoryNone CollisionCategory = 0

    // Collision with everything. All bits are set to 1.
    CollisionCategoryAll CollisionCategory = math.MaxUint16
)
