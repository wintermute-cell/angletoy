package ik

import (
	"gorl/fw/util"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bone struct {
	position  rl.Vector2
	length    float32
	angle     float32 // radians
	min_angle float32
	max_angle float32
}

// NewBone creates a new bone.
func NewBone(position rl.Vector2, length, angle float32) *Bone {
	return &Bone{
		position: position,
		length:   length,
		angle:    angle,
	}
}

func (b *Bone) SetConstraints(min_angle, max_angle float32) {
	b.min_angle = min_angle
	b.max_angle = max_angle
}

// GetPosition returns the position of the bone.
func (b *Bone) GetPosition() rl.Vector2 {
	return b.position
}

// SetPosition updates the position of the bone.
func (b *Bone) SetPosition(position rl.Vector2) {
	b.position = position
}

// GetLength returns the length of the bone.
func (b *Bone) GetLength() float32 {
	return b.length
}

// GetAngle returns the angle of the bone in radians.
func (b *Bone) GetAngle() float32 {
	return b.angle
}

// SetAngle updates the angle of the bone.
func (b *Bone) SetAngle(angle float32) {
	b.angle = angle
}

// RotateZ returns the rotation matrix for a given angle in radians.
func RotateZ(theta float64) [2][2]float32 {
	return [2][2]float32{
		{float32(math.Cos(theta)), float32(-math.Sin(theta))},
		{float32(math.Sin(theta)), float32(math.Cos(theta))},
	}
}

// ApplyTransformation applies the transformation matrix to a vector.
func ApplyTransformation(matrix [2][2]float32, vector rl.Vector2) rl.Vector2 {
	return rl.Vector2{
		X: matrix[0][0]*vector.X + matrix[0][1]*vector.Y,
		Y: matrix[1][0]*vector.X + matrix[1][1]*vector.Y,
	}
}

// FK computes the end positions of the bones using forward kinematics.
func FK(bones []*Bone) []rl.Vector2 {
	positions := make([]rl.Vector2, len(bones)+1)
	positions[0] = bones[0].GetPosition()

	for i, bone := range bones {
		rotation := RotateZ(float64(bone.GetAngle()))
		direction := ApplyTransformation(rotation, rl.Vector2{X: bone.GetLength(), Y: 0})
		positions[i+1] = rl.Vector2Add(positions[i], direction)

		// If there's a next bone, update its starting position to be this bone's end position
		if i+1 < len(bones) {
			bones[i+1].SetPosition(positions[i+1])
		}
	}

	return positions
}

// IK applies inverse kinematics to adjust the bone angles to reach the target.
func IK(target rl.Vector2, bones []*Bone, maxIter int, errMin float32) (int, bool) {
	for loop := 0; loop < maxIter; loop++ {
		positions := FK(bones)

		for i := len(bones) - 1; i >= 0; i-- {
			curToEnd := rl.Vector2Subtract(positions[len(positions)-1], positions[i])
			curToTarget := rl.Vector2Subtract(target, positions[i])

			endTargetMag := rl.Vector2Length(curToEnd) * rl.Vector2Length(curToTarget)

			if endTargetMag <= 0.0001 {
				continue
			}

			cosRotAng := rl.Vector2DotProduct(curToEnd, curToTarget) / endTargetMag
			sinRotAng := (curToEnd.X*curToTarget.Y - curToEnd.Y*curToTarget.X) / endTargetMag

			rotAng := float32(math.Acos(float64(util.Clamp(cosRotAng, -1.0, 1.0))))
			if sinRotAng < 0.0 {
				rotAng = -rotAng
			}

			newAngle := bones[i].GetAngle() + rotAng
			// Normalize angle to [0, 360]
			newAngle = float32(math.Mod(float64(newAngle+2*math.Pi), 2*math.Pi))
			bones[i].SetAngle(newAngle)

			var parent *Bone = nil
			if i > 0 {
				parent = bones[i-1]
			}
			bones[i].enforceConstraints(parent)
			//rl.DrawText(
			//    fmt.Sprintf("%.1f || %.1f", relative_angle, newAngle),
			//    int32(bones[i].GetPosition().X),
			//    int32(bones[i].GetPosition().Y),
			//    4, rl.White)
		}
	}

	return maxIter, true
}

func (b *Bone) enforceConstraints(parent *Bone) float32 {
	if parent == nil {
		return b.angle
	}

	// Convert global angle to local angle relative to the parent
	localAngle := b.angle - parent.angle

	// Ensure local angle is within the range [0, 2π]
	for localAngle < 0 {
		localAngle += 2 * math.Pi
	}
	for localAngle >= 2*math.Pi {
		localAngle -= 2 * math.Pi
	}

	// do nothing if both constraints are 0
	if b.min_angle+b.max_angle == 0 {
		return localAngle
	}

	// Apply constraints to the local angle
	if b.min_angle <= b.max_angle {
		// Standard case, no wrapping
		if localAngle < b.min_angle {
			localAngle = b.min_angle
		} else if localAngle > b.max_angle {
			localAngle = b.max_angle
		}
	} else {
		// Wrapping case, such as [4, 0.2] but adjusted for radians
		if localAngle < b.min_angle && localAngle > b.max_angle {
			if math.Abs(float64(localAngle-b.min_angle)) < math.Abs(float64(localAngle-b.max_angle)) {
				localAngle = b.min_angle
			} else {
				localAngle = b.max_angle
			}
		}
	}

	// Convert local angle back to global
	b.angle = localAngle + parent.angle

	return localAngle
}
