package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/settings"
	"gorl/fw/util"
	"gorl/game/code/colorscheme"
	"math"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that AngleShowcaserEntity implements IEntity.
var _ entities.IEntity = &AngleShowcaserEntity{}

// AngleShowcaser Entity
type AngleShowcaserEntity struct {
	*entities.Entity // Required!

	angleFuncs              []AngleFunc
	showcaseCircleRadius    float32
	showcaseCirclePositions []rl.Vector2
	pointerPositions        []rl.Vector2
	calculatedAngles        []float32
}

// NewAngleShowcaserEntity creates a new instance of the AngleShowcaserEntity.
func NewAngleShowcaserEntity() *AngleShowcaserEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	radius := float32(100.0)
	new_ent := &AngleShowcaserEntity{
		Entity: entities.NewEntity("AngleShowcaserEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		angleFuncs: []AngleFunc{
			DotProductAngleFunc,
			CrossProduct2DAngleFunc,
			Atan2AngleFunc,
		},
		showcaseCircleRadius: radius,
		showcaseCirclePositions: []rl.Vector2{
			{
				X: float32(settings.CurrentSettings().RenderWidth / 4),
				Y: float32(settings.CurrentSettings().RenderHeight / 2),
			},
			{
				X: float32(settings.CurrentSettings().RenderWidth / 4 * 2),
				Y: float32(settings.CurrentSettings().RenderHeight / 2),
			},
			{
				X: float32(settings.CurrentSettings().RenderWidth / 4 * 3),
				Y: float32(settings.CurrentSettings().RenderHeight / 2),
			},
		},
		calculatedAngles: []float32{0, 0, 0},
	}

	for _, pos := range new_ent.showcaseCirclePositions {
		new_ent.pointerPositions = append(
			new_ent.pointerPositions,
			rl.NewVector2(pos.X, pos.Y-radius),
		)
	}
	return new_ent
}

func (ent *AngleShowcaserEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *AngleShowcaserEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *AngleShowcaserEntity) Update() {
	// Update logic for the entity per frame
	// ...
}

func (ent *AngleShowcaserEntity) Draw() {
	// Draw logic for the entity
	// ...
	rl.DrawText("Calculating angles between two vectors can have some unforseen quirks.", 510, 200, 20, rl.White)

	approaches := []string{
		"Dot Product",
		"Cross Product 2D",
		"Atan2 of Cross and Dot",
	}
	formulas := []string{
		"angle = arccos((A . B) / (|A| * |B|))",
		"angle = (A x B) / (|A| * |B|)",
		"angle = atan2((A x B), (A . B))",
	}
	infoText := []string{
		`Angle Range: [0, π] radians or [0, 180] degrees.
Quirks:
  - Only provides the smallest angle between vectors.
  - Ignores the direction of the vectors relative to each other (cannot
	determine if one vector is clockwise or counterclockwise to the other).`,
		`Angle Range: [-π, π] radians or [-180, 180] degrees.
Quirks:
  - Provides signed angles, indicating direction
	(positive for counterclockwise, negative for clockwise).
  - Can result in smaller magnitude angles for obtuse 
	angles unless adjusted for a full range with additional checks.`,
		`Angle Range: [-π, π] radians or [-180, 180] degrees.
Quirks:
  - Returns the full angle between vectors, considering direction.
  - Handles cases where vectors are aligned with axes well.
  - Particularly useful in 2D for obtaining a signed angle and
	determining the rotation direction.`,
	}
	for i := range ent.showcaseCirclePositions {
		rl.DrawCircleV(ent.showcaseCirclePositions[i], ent.showcaseCircleRadius, colorscheme.Colorscheme.Color02.ToRGBA())
		textWidth := rl.MeasureText("Approach: "+approaches[i], 20)
		//rl.DrawText("Approach: "+approaches[i], int32(ent.showcaseCirclePositions[i].X-150), int32(ent.showcaseCirclePositions[i].Y+ent.showcaseCircleRadius+40), 20, rl.White)
		rl.DrawText("Approach: "+approaches[i],
			int32(ent.showcaseCirclePositions[i].X-float32(textWidth)/2),
			int32(ent.showcaseCirclePositions[i].Y+ent.showcaseCircleRadius+40),
			20, rl.White)
		// formulas
		textWidth = rl.MeasureText(formulas[i], 20)
		rl.DrawText(formulas[i],
			int32(ent.showcaseCirclePositions[i].X-float32(textWidth)/2),
			int32(ent.showcaseCirclePositions[i].Y+ent.showcaseCircleRadius+70),
			20, rl.White)

		// info text with angle range and quirks
		rl.DrawText(infoText[i],
			int32(ent.showcaseCirclePositions[i].X-200),
			int32(ent.showcaseCirclePositions[i].Y+ent.showcaseCircleRadius+120),
			10, rl.White)

		// reference up line
		rl.DrawLineEx(
			ent.showcaseCirclePositions[i],
			rl.NewVector2(ent.showcaseCirclePositions[i].X, ent.showcaseCirclePositions[i].Y-ent.showcaseCircleRadius),
			3,
			colorscheme.Colorscheme.Color01.ToRGBA(),
		)
	}

	for i := range ent.pointerPositions {
		mDir := util.Vector2NormalizeSafe(rl.Vector2Subtract(ent.pointerPositions[i], ent.showcaseCirclePositions[i]))
		lineEnd := rl.Vector2Add(ent.showcaseCirclePositions[i], rl.Vector2Scale(mDir, ent.showcaseCircleRadius))
		//rl.DrawLineV(
		//	ent.showcaseCirclePositions[i],
		//	lineEnd,
		//	colorscheme.Colorscheme.Color10.ToRGBA(),
		//)
		rl.DrawLineEx(
			ent.showcaseCirclePositions[i],
			lineEnd,
			3,
			colorscheme.Colorscheme.Color10.ToRGBA(),
		)
	}

	for i, angle := range ent.calculatedAngles {
		mDir := util.Vector2NormalizeSafe(rl.Vector2Subtract(ent.pointerPositions[i], ent.showcaseCirclePositions[i]))
		lineEnd := rl.Vector2Add(ent.showcaseCirclePositions[i], rl.Vector2Scale(mDir, ent.showcaseCircleRadius))
		rl.DrawText(
			"Angle: "+strconv.FormatFloat(float64(angle), 'f', 2, 32)+" radians",
			int32(lineEnd.X+10),
			int32(lineEnd.Y),
			20,
			rl.White,
		)
		rl.DrawText(
			"Angle: "+strconv.FormatFloat(float64(angle*rl.Rad2deg), 'f', 2, 32)+" degrees",
			int32(lineEnd.X+10),
			int32(lineEnd.Y+20),
			20,
			rl.White,
		)
	}
}

func (ent *AngleShowcaserEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.

	if event.Action == input.ActionClickHeld {
		mPos := event.GetScreenSpaceMousePosition()
		for idx, pos := range ent.showcaseCirclePositions {
			if rl.CheckCollisionPointCircle(mPos, pos, ent.showcaseCircleRadius) {
				toMouse := rl.Vector2Subtract(mPos, pos)
				ent.pointerPositions[idx] = mPos
				upDirection := rl.NewVector2(0, -1)
				angle := ent.angleFuncs[idx](upDirection, toMouse)
				ent.calculatedAngles[idx] = angle
			}
		}
	}

	return true
}

type AngleFunc func(a, b rl.Vector2) float32

// Angle Range: [0, π] radians or [0, 180] degrees.
// Quirks:
//   - Only provides the smallest angle between vectors.
//   - Ignores the direction of the vectors relative to each other (cannot
//     determine if one vector is clockwise or counterclockwise to the other).
func DotProductAngleFunc(a, b rl.Vector2) float32 {
	lenA := rl.Vector2Length(a)
	lenB := rl.Vector2Length(b)
	if lenA == 0 || lenB == 0 {
		return 0
	}

	dot := rl.Vector2DotProduct(a, b)
	cosTheta := dot / (lenA * lenB)
	return float32(math.Acos(float64(cosTheta)))
}

// Angle Range: [-π, π] radians or [-180, 180] degrees.
// Quirks:
//   - Provides signed angles, indicating direction (positive for counterclockwise, negative for clockwise).
//   - Can result in smaller magnitude angles for obtuse angles unless adjusted for a full range with additional checks.
func CrossProduct2DAngleFunc(a, b rl.Vector2) float32 {
	lenA := rl.Vector2Length(a)
	lenB := rl.Vector2Length(b)
	if lenA == 0 || lenB == 0 {
		return 0
	}

	cross := rl.Vector2CrossProduct(a, b)
	sinTheta := cross / (lenA * lenB)
	return float32(math.Asin(float64(sinTheta)))
}

// Angle Range: [-π, π] radians or [-180, 180] degrees.
// Quirks:
//   - Returns the full angle between vectors, considering direction.
//   - Handles cases where vectors are aligned with axes well.
//   - Particularly useful in 2D for obtaining a signed angle and determining the rotation direction.
func Atan2AngleFunc(a, b rl.Vector2) float32 {
	//return float32(math.Atan2(float64(b.Y), float64(b.X)) - math.Atan2(float64(a.Y), float64(a.X)))
	cross := rl.Vector2CrossProduct(a, b)
	dot := rl.Vector2DotProduct(a, b)
	return float32(math.Atan2(float64(cross), float64(dot)))
}
