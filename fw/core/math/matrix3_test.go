package math

import (
	"github.com/gen2brain/raylib-go/raylib"
	"math"
	"testing"
)

// TestMatrix3Translation tests the translation matrix creation
func TestMatrix3Translation(t *testing.T) {
	v := rl.Vector2{X: 10, Y: 20}
	m := Matrix3Translation(v)
	if m.m2 != v.X || m.m5 != v.Y {
		t.Errorf("Translation matrix does not translate correctly, got %v", m)
	}
}

// TestMatrix3Rotation tests the rotation matrix creation
func TestMatrix3Rotation(t *testing.T) {
	angle := float32(90) // 90 degrees
	m := Matrix3Rotation(angle)
	// Expected values for 90 degrees rotation
	expectedCos := float32(math.Cos(math.Pi / 2))
	expectedSin := float32(math.Sin(math.Pi / 2))

	if !almostEqualFloats(m.m0, expectedCos, 0.001) || !almostEqualFloats(m.m1, -expectedSin, 0.001) ||
		!almostEqualFloats(m.m3, expectedSin, 0.001) || !almostEqualFloats(m.m4, expectedCos, 0.001) {
		t.Errorf("Rotation matrix does not rotate correctly, got %v", m)
	}
}

// TestMatrix3Scale tests the scaling matrix creation
func TestMatrix3Scale(t *testing.T) {
	v := rl.Vector2{X: 2, Y: 3}
	m := Matrix3Scale(v)
	if m.m0 != v.X || m.m4 != v.Y {
		t.Errorf("Scaling matrix does not scale correctly, got %v", m)
	}
}

// TestFromTransformations tests the combination of translation, rotation, and scaling
func TestFromTransformations(t *testing.T) {
	position := rl.Vector2{X: 1, Y: 2}
	rotation := float32(45) // 45 degrees
	scale := rl.Vector2{X: 2, Y: 3}

	transformMatrix := FromTransformations(position, rotation, scale)

	// Manually combine matrices to test against
	tMatrix := Matrix3Translation(position)
	rMatrix := Matrix3Rotation(rotation)
	sMatrix := Matrix3Scale(scale)

	// The order should be scale, then rotation, then translation
	expectedMatrix := tMatrix.Multiply(rMatrix).Multiply(sMatrix)

	// Compare all elements
	if !almostEqualMatrix3(expectedMatrix, transformMatrix, 0.001) {
		t.Errorf("Generated matrix does not match expected matrix. Got %v, want %v", transformMatrix, expectedMatrix)
	}
}
