package math

import (
	"math"
	"testing"

	"github.com/gen2brain/raylib-go/raylib"
)

func TestGenerateMatrix(t *testing.T) {
	// Create a Transform2D with non-default values
	transform := NewTransform2D(rl.Vector2{X: 10, Y: 20}, 45, rl.Vector2{X: 2, Y: 2})

	// Generate the matrix
	matrix := transform.GenerateMatrix()

	// Test the matrix values
	expectedTranslation := Matrix3Translation(transform.position)
	expectedRotation := Matrix3Rotation(transform.rotation)
	expectedScale := Matrix3Scale(transform.scale)

	expectedMatrix := expectedTranslation.Multiply(expectedRotation).Multiply(expectedScale)

	if !almostEqualMatrix3(matrix, expectedMatrix, 0.001) {
		t.Errorf("Generated matrix does not match expected matrix. Got %v, want %v", matrix, expectedMatrix)
	}

	// Modify the transform and test the matrix again
	transform.SetPosition(rl.Vector2{X: 30, Y: 40})
	transform.SetRotation(90)
	transform.SetScale(rl.Vector2{X: 3, Y: 3})

	matrix = transform.GenerateMatrix()

	expectedTranslation = Matrix3Translation(transform.position)
	expectedRotation = Matrix3Rotation(transform.rotation)
	expectedScale = Matrix3Scale(transform.scale)
	expectedMatrix = expectedTranslation.Multiply(expectedRotation).Multiply(expectedScale)

	if !almostEqualMatrix3(matrix, expectedMatrix, 0.001) {
		t.Errorf("Generated matrix does not match expected matrix. Got %v, want %v", matrix, expectedMatrix)
	}
}

func TestNewTransform2DFromMatrix3(t *testing.T) {
	// Assume Matrix3Identity, Matrix3Translation, Matrix3Rotation, Matrix3Scale are available
	originalPosition := rl.Vector2{X: 30, Y: 40}
	originalRotation := float32(90) // degrees
	originalScale := rl.Vector2{X: 3, Y: 3}

	// Construct a combined transformation matrix
	matrix := Matrix3Identity()
	matrix = matrix.Multiply(Matrix3Translation(originalPosition))
	matrix = matrix.Multiply(Matrix3Rotation(originalRotation))
	matrix = matrix.Multiply(Matrix3Scale(originalScale))

	// Create a Transform2D from this matrix
	transform := NewTransform2DFromMatrix3(matrix)

	// Validate the transform properties
	if !almostEqualVector2(transform.position, originalPosition, 0.001) {
		t.Errorf("Position does not match. Got %v, want %v", transform.position, originalPosition)
	}
	if !almostEqualFloats(transform.rotation, originalRotation, 0.001) {
		t.Errorf("Rotation does not match. Got %v, want %v", transform.rotation, originalRotation)
	}
	if !almostEqualVector2(transform.scale, originalScale, 0.001) {
		t.Errorf("Scale does not match. Got %v, want %v", transform.scale, originalScale)
	}
}

// Helper functions to compare matrices, vectors and floats with a tolerance
func almostEqualMatrix3(m1, m2 Matrix3, tolerance float32) bool {
	return almostEqualFloats(m1.m0, m2.m0, tolerance) &&
		almostEqualFloats(m1.m1, m2.m1, tolerance) &&
		almostEqualFloats(m1.m2, m2.m2, tolerance) &&
		almostEqualFloats(m1.m3, m2.m3, tolerance) &&
		almostEqualFloats(m1.m4, m2.m4, tolerance) &&
		almostEqualFloats(m1.m5, m2.m5, tolerance) &&
		almostEqualFloats(m1.m6, m2.m6, tolerance) &&
		almostEqualFloats(m1.m7, m2.m7, tolerance) &&
		almostEqualFloats(m1.m8, m2.m8, tolerance)
}

func almostEqualVector2(v1, v2 rl.Vector2, tolerance float32) bool {
	return math.Abs(float64(v1.X-v2.X)) < float64(tolerance) && math.Abs(float64(v1.Y-v2.Y)) < float64(tolerance)
}

func almostEqualFloats(a, b, tolerance float32) bool {
	return math.Abs(float64(a-b)) < float64(tolerance)
}
