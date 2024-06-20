package math

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ============================================================================
// Vector2Int defines an integer 2D vector.
// Operations are defined as methods and raylib-go style functions.
// ============================================================================

// Vector2Int represents a 2D vector with integer components.
type Vector2Int struct {
	X int
	Y int
}

// NewVector2Int creates a new Vector2Int.
func NewVector2Int(x, y int) Vector2Int {
	return Vector2Int{X: x, Y: y}
}

// Vector2IntZero returns a vector with both components set to 0.
func Vector2IntZero() Vector2Int {
	return Vector2Int{X: 0, Y: 0}
}

// Vector2IntOne returns a vector with both components set to 1.
func Vector2IntOne() Vector2Int {
	return Vector2Int{X: 1, Y: 1}
}

// Vector2IntFromRl returns a Vector2Int from a raylib Vector2.
// The components are truncated to integers.
func Vector2IntFromRl(v rl.Vector2) Vector2Int {
	return Vector2Int{X: int(v.X), Y: int(v.Y)}
}

// ToRl returns a raylib Vector2 from a Vector2Int.
func (v Vector2Int) ToRl() rl.Vector2 {
	return rl.NewVector2(float32(v.X), float32(v.Y))
}

// Add returns the sum of two vectors.
func (v Vector2Int) Add(v2 Vector2Int) Vector2Int {
	return Vector2Int{X: v.X + v2.X, Y: v.Y + v2.Y}
}

// Vector2IntAdd returns the sum of two vectors.
// This is the same as v1.Add(v2).
func Vector2IntAdd(v1, v2 Vector2Int) Vector2Int {
	return Vector2Int{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

// Sub returns the difference of two vectors.
func (v Vector2Int) Sub(v2 Vector2Int) Vector2Int {
	return Vector2Int{X: v.X - v2.X, Y: v.Y - v2.Y}
}

// Vector2IntSub returns the difference of two vectors.
// This is the same as v1.Sub(v2).
func Vector2IntSub(v1, v2 Vector2Int) Vector2Int {
	return Vector2Int{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

// Mul returns the vector scaled by a scalar value.
func (v Vector2Int) Mul(scalar int) Vector2Int {
	return Vector2Int{X: v.X * scalar, Y: v.Y * scalar}
}

// Vector2IntMul returns the vector scaled by a scalar value.
// This is the same as v.Mul(scalar).
func Vector2IntMul(v Vector2Int, scalar int) Vector2Int {
	return Vector2Int{X: v.X * scalar, Y: v.Y * scalar}
}

// Dot returns the dot product of two vectors.
func (v Vector2Int) Dot(v2 Vector2Int) int {
	return v.X*v2.X + v.Y*v2.Y
}

// Vector2IntDot returns the dot product of two vectors.
// This is the same as v1.Dot(v2).
func Vector2IntDot(v1, v2 Vector2Int) int {
	return v1.X*v2.X + v1.Y*v2.Y
}

// Magnitude returns the magnitude (length) of the vector.
func (v Vector2Int) Magnitude() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}

// Vector2IntMagnitude returns the magnitude (length) of the vector.
// This is the same as v.Magnitude().
func Vector2IntMagnitude(v Vector2Int) float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}

// Equals checks if two vectors are equal.
func (v Vector2Int) Equals(v2 Vector2Int) bool {
	return v.X == v2.X && v.Y == v2.Y
}

// Vector2IntEquals checks if two vectors are equal.
// This is the same as v1.Equals(v2).
func Vector2IntEquals(v1, v2 Vector2Int) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}

// ToString returns the string representation of the vector.
func (v Vector2Int) ToString() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}
