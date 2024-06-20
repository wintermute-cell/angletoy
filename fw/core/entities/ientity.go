package entities

import (
	input "gorl/fw/core/input/input_handling"
	"gorl/fw/core/math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// IEntity is an interface that every entity in the game should implement
type IEntity interface {
	// Lifecycle methods
	Init()
	Deinit()

	// Per-frame methods
	Update()
	FixedUpdate()
	Draw()

	// Transform
	GetPosition() rl.Vector2
	SetPosition(new_position rl.Vector2)
	GetScale() rl.Vector2
	SetScale(new_size rl.Vector2)
	SetRotation(new_rotation float32)
	GetRotation() float32
	math.Transformable // provides GetTransform()

	// OnInputEvent is called when an input event is received.
	// The entity must decide if it should handle the event or not.
	// Return false if the event should not be propagated further.
	input.InputReceiver // provides OnInputEvent()

	// Rendering
	GetDrawIndex() int32
	SetDrawIndex(index int32)
	IsEnabled() bool
	IsVisible() bool
	GetLayerFlags() math.BitFlag

	// Other
	GetName() string
}
