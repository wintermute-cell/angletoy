package math

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Matrix3 is a 3x3 matrix
type Matrix3 struct {
	// | m0 m1 m2 |
	// | m3 m4 m5 |
	// | m6 m7 m8 |
	m0, m1, m2, m3, m4, m5, m6, m7, m8 float32
}

// String repr with padding
func (m Matrix3) String() string {
	return fmt.Sprintf("\n| %f %f %f |\n| %f %f %f |\n| %f %f %f |\n",
		m.m0, m.m1, m.m2,
		m.m3, m.m4, m.m5,
		m.m6, m.m7, m.m8,
	)
}

// ============================================================================
//		CONSTRUCTORS
// ============================================================================

// Matrix3Identity creates a new 3x3 matrix with the identity values.
func Matrix3Identity() Matrix3 {
	return Matrix3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
}

// NewMatrix3 creates a new 3x3 matrix with the given values.
func NewMatrix3(
	m0, m1, m2,
	m3, m4, m5,
	m6, m7, m8 float32,
) Matrix3 {
	return Matrix3{
		m0, m1, m2,
		m3, m4, m5,
		m6, m7, m8,
	}
}

// Matrix3Translation creates a new translation matrix.
func Matrix3Translation(v rl.Vector2) Matrix3 {
	return Matrix3{
		1, 0, v.X,
		0, 1, v.Y,
		0, 0, 1,
	}
}

// Matrix3Rotation creates a new rotation matrix.
func Matrix3Rotation(angle float32) Matrix3 {
	rad := angle * (math.Pi / 180)
	cos := float32(math.Cos(float64(rad)))
	sin := float32(math.Sin(float64(rad)))

	return Matrix3{
		cos, -sin, 0,
		sin, cos, 0,
		0, 0, 1,
	}
}

// Matrix3Scale creates a new scaling matrix.
func Matrix3Scale(v rl.Vector2) Matrix3 {
	return Matrix3{
		v.X, 0, 0,
		0, v.Y, 0,
		0, 0, 1,
	}
}

// FromTransformations creates a transformation matrix from position, rotation, and scale.
func FromTransformations(position rl.Vector2, rotation float32, scale rl.Vector2) Matrix3 {
	translationMatrix := Matrix3Translation(position)
	rotationMatrix := Matrix3Rotation(rotation)
	scaleMatrix := Matrix3Scale(scale)

	// Apply transformations in the order: scale -> rotation -> translation
	transformMatrix := translationMatrix.Multiply(rotationMatrix).Multiply(scaleMatrix)
	return transformMatrix
}

// ============================================================================
//		METHODS
// ============================================================================

// SetIdentity sets the matrix to the identity matrix.
func (m *Matrix3) SetIdentity() {
	m.m0 = 1
	m.m1 = 0
	m.m2 = 0
	m.m3 = 0
	m.m4 = 1
	m.m5 = 0
	m.m6 = 0
	m.m7 = 0
	m.m8 = 1
}

// Multiply multiplies two matrices
func (a Matrix3) Multiply(b Matrix3) Matrix3 {
	return Matrix3{
		a.m0*b.m0 + a.m1*b.m3 + a.m2*b.m6,
		a.m0*b.m1 + a.m1*b.m4 + a.m2*b.m7,
		a.m0*b.m2 + a.m1*b.m5 + a.m2*b.m8,
		a.m3*b.m0 + a.m4*b.m3 + a.m5*b.m6,
		a.m3*b.m1 + a.m4*b.m4 + a.m5*b.m7,
		a.m3*b.m2 + a.m4*b.m5 + a.m5*b.m8,
		a.m6*b.m0 + a.m7*b.m3 + a.m8*b.m6,
		a.m6*b.m1 + a.m7*b.m4 + a.m8*b.m7,
		a.m6*b.m2 + a.m7*b.m5 + a.m8*b.m8,
	}
}

// MultiplyV applies a matrix transformation to a Vector2.
func (m Matrix3) MultiplyV(v rl.Vector2) rl.Vector2 {
	return rl.Vector2{
		X: m.m0*v.X + m.m1*v.Y + m.m2,
		Y: m.m3*v.X + m.m4*v.Y + m.m5,
	}
}
