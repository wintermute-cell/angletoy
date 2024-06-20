package util

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Vector2Clamp restricts a vector within the limits specified by min and max vectors.
func Vector2Clamp(input, min, max rl.Vector2) rl.Vector2 {
	if input.X < min.X {
		input.X = min.X
	} else if input.X > max.X {
		input.X = max.X
	}
	if input.Y < min.Y {
		input.Y = min.Y
	} else if input.Y > max.Y {
		input.Y = max.Y
	}
	return input
}

type number interface {
	int | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type signed_number interface {
	int | int16 | int32 | int64 | float32 | float64
}

// Sign returns the sign of x (-1 if x < 0, 1 if x > 0, 0 if x == 0)
func Sign[T signed_number](x T) T {
	if x < 0 {
		return T(-1)
	} else if x > 0 {
		return T(1)
	} else {
		return T(0)
	}
}

// Abs returns the absolute value of x
func Abs[T number](x T) T {
	ret := x
	if x < 0 {
		ret = -x
	}
	return ret
}

// Max will return the maximum value between x and y
func Max[T number](x, y T) T {
	return T(math.Max(float64(x), float64(y)))
}

// Min will return the minimum value between x and y.
func Min[T number](x, y T) T {
	return T(math.Min(float64(x), float64(y)))
}

// Clamps x between lower_bound and upper_bound, both inclusive.
// (Clamp will return at least lower_bound and at most upper_bound)
func Clamp[T number](x, lower_bound, upper_bound T) T {
	v := Min[T](x, upper_bound)
	v = Max[T](v, lower_bound)
	return v
}

// Round x to the nearest integer, either down if x < .5 or up if x >= .5.
func Round[T number](x T) T {
	integer, fraction := math.Modf(float64(x))
	v := x
	if fraction >= 0.5 {
		v = T(integer) + T(1.0)
	} else {
		v = T(integer)
	}
	return v
}

// Vector2NormalizeSafe returns the normalized vector. If the input vector is
// zero, it returns a zero vector instead of (NaN, NaN).
func Vector2NormalizeSafe(v rl.Vector2) rl.Vector2 {
	if v == rl.Vector2Zero() {
		return rl.Vector2Zero()
	} else {
		return rl.Vector2Normalize(v)
	}
}

// Vector2Lerp returns the linear interpolation between two numbers.
func Lerp[T number](a, b T, factor float32) T {
	return T(float32(a)*(1.0-factor) + (float32(b) * factor))
}

// ShortestLerp returns the shortest linear interpolation between two numbers.
func ShortestLerp(current, target, factor float32) float32 {
	// Calculate the difference
	difference := target - current

	// Calculate possible wrapped differences
	wrappedDifferencePlus := float64(difference + 360)
	wrappedDifferenceMinus := float64(difference - 360)

	// Check which one is the smallest in terms of absolute value
	if math.Abs(wrappedDifferencePlus) < math.Abs(float64(difference)) {
		difference = float32(wrappedDifferencePlus)
	} else if math.Abs(wrappedDifferenceMinus) < math.Abs(float64(difference)) {
		difference = float32(wrappedDifferenceMinus)
	}

	// Compute the lerped value
	lerped := current + difference*factor

	// Adjust the lerped value to be within the 0-360 range
	for lerped < 0 {
		lerped += 360
	}
	for lerped >= 360 {
		lerped -= 360
	}

	return lerped
}

// Vector2Angle returns the angle of a vector in degrees.
func Vector2Angle(v rl.Vector2) float32 {
	const RadToDeg = 180.0 / math.Pi

	// Get the angle in radians
	radian := math.Atan2(float64(v.Y), float64(v.X))

	// Convert the angle to degrees
	degree := float32(radian) * RadToDeg

	// Normalize the degree to be in [0, 360]
	if degree < 0 {
		degree += 360
	}

	return degree
}

// RotatePointAroundOrigin rotates a point around an origin by a given angle.
func RotatePointAroundOrigin(point, origin rl.Vector2, angle_deg float32) rl.Vector2 {
	angleRad := angle_deg * (math.Pi / 180) // Convert angle to radians

	cosAngle := float32(math.Cos(float64(angleRad)))
	sinAngle := float32(math.Sin(float64(angleRad)))

	// Translate point back to origin
	translated := rl.Vector2Subtract(point, origin)

	// Rotate point
	rotatedX := translated.X*cosAngle - translated.Y*sinAngle
	rotatedY := translated.X*sinAngle + translated.Y*cosAngle

	// Translate point back
	finalPoint := rl.NewVector2(rotatedX+origin.X, rotatedY+origin.Y)

	return finalPoint
}

// ChildToWorldSpace converts a relative child_position to world space given
// the root's position and rotation.
func ChildToWorldSpace(root_position, child_relative_position rl.Vector2, root_rotation_degrees float32) rl.Vector2 {
	// Convert rotation to radians
	root_rotation := float32(math.Pi) * root_rotation_degrees / 180.0

	// Compute the rotation matrix elements
	sinTheta := float32(math.Sin(float64(root_rotation)))
	cosTheta := float32(math.Cos(float64(root_rotation)))

	// Rotate the relative child position using the rotation matrix
	rotatedX := child_relative_position.X*cosTheta - child_relative_position.Y*sinTheta
	rotatedY := child_relative_position.X*sinTheta + child_relative_position.Y*cosTheta

	// Translate the rotated position by the root's position to get the world space position
	worldX := rotatedX + root_position.X
	worldY := rotatedY + root_position.Y

	return rl.NewVector2(worldX, worldY)
}

func Vector2MoveTowards(current, target rl.Vector2, step float32) rl.Vector2 {
	// Vector from current to target
	delta := rl.Vector2Subtract(target, current)

	distanceToTarget := rl.Vector2Length(delta)

	// If the distance to move is less than the step, just return the target position.
	if distanceToTarget < step {
		return target
	}

	// Normalize the delta vector
	dir := Vector2NormalizeSafe(delta)

	// Move the vector by step towards target
	moved := rl.Vector2Add(current, rl.Vector2Scale(dir, step))

	return moved
}
