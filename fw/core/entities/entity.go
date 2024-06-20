package entities

import (
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ IEntity = (*Entity)(nil)

// Base implementation of IEntity. Should be embedded in every custom entity.
type Entity struct {
	Name string

	enabled bool // If false, the entity will not be updated or drawn.
	visible bool // If false, the entity will not be drawn.

	// Transform2D is a struct that holds the position, rotation and scale of the entity.
	transform math.Transform2D

	// DrawIndex is used to determine the order in which entities are drawn.
	// Lower values are drawn behind higher values.
	drawIndex int32

	// LayerFlags is a bit flag that determines which layers the entity belongs to.
	// Cameras can selectively render entities based on their layer flags.
	// (Layer flags are not automatically inherited by children.)
	layerFlags math.BitFlag
}

// NewEntity creates a new base implementation of IEntity.
// This should be called by the constructor of a custom entity.
func NewEntity(name string, position rl.Vector2, rotation float32, scale rl.Vector2) *Entity {
	return &Entity{
		Name:       name,
		enabled:    true,
		visible:    true,
		transform:  math.NewTransform2D(position, rotation, scale),
		drawIndex:  0,
		layerFlags: math.Flag0,
	}
}

func (ent *Entity) Init()        {} // Should be overridden by the custom entity.
func (ent *Entity) Deinit()      {} // Should be overridden by the custom entity.
func (ent *Entity) Update()      {} // Should be overridden by the custom entity.
func (ent *Entity) FixedUpdate() {} // Should be overridden by the custom entity.
func (ent *Entity) Draw()        {} // Should be overridden by the custom entity.

// String returns the name of the entity.
func (ent *Entity) String() string {
	return ent.Name
}

// IsEnabled returns true if the entity is enabled.
func (ent *Entity) IsEnabled() bool {
	return ent.enabled
}

// SetEnabled sets the enabled state of the entity.
func (ent *Entity) SetEnabled(enabled bool) {
	ent.enabled = enabled
}

// IsVisible returns true if the entity is visible.
func (ent *Entity) IsVisible() bool {
	return ent.visible
}

// SetVisible sets the visible state of the entity.
func (ent *Entity) SetVisible(visible bool) {
	ent.visible = visible
}

// GetPosition returns the position of the entity.
func (ent *Entity) GetPosition() rl.Vector2 {
	return ent.transform.GetPosition()
}

// SetPosition sets the position of the entity.
func (ent *Entity) SetPosition(newPosition rl.Vector2) {
	ent.transform.SetPosition(newPosition)
}

// GetScale returns the scale of the entity.
func (ent *Entity) GetScale() rl.Vector2 {
	return ent.transform.GetScale()
}

// SetScale sets the scale of the entity.
func (ent *Entity) SetScale(newScale rl.Vector2) {
	ent.transform.SetScale(newScale)
}

// GetRotation returns the rotation of the entity.
func (ent *Entity) GetRotation() float32 {
	return ent.transform.GetRotation()
}

// SetRotation sets the rotation of the entity.
func (ent *Entity) SetRotation(newRotation float32) {
	ent.transform.SetRotation(newRotation)
}

// GetTransform returns the transform of the entity.
// Should generally not be modified directly.
// Use SetPosition, SetScale and SetRotation instead.
func (ent *Entity) GetTransform() *math.Transform2D {
	return &ent.transform
}

// SetTransform overwrites the transform of the entity.
// This includes the position, rotation and scale.
func (ent *Entity) SetTransform(newTransform math.Transform2D) {
	ent.transform = newTransform
}

// OnInputEvent is called when an input event is received.
// The entity must decide if it should handle the event or not.
// Return false if the event should not be propagated further.
func (ent *Entity) OnInputEvent(event *input.InputEvent) bool {
	return true
}

// GetDrawIndex returns the draw index of the entity.
// Lower values are drawn behind higher values.
func (ent *Entity) GetDrawIndex() int32 {
	return ent.drawIndex
}

// SetDrawIndex sets the draw index of the entity.
// Lower values are drawn behind higher values.
func (ent *Entity) SetDrawIndex(index int32) {
	ent.drawIndex = index
}

// GetName returns the name of the entity.
func (ent *Entity) GetName() string {
	return ent.Name
}

// GetLayerFlags returns the layer flags of the entity.
func (ent *Entity) GetLayerFlags() math.BitFlag {
	return ent.layerFlags
}

// SetLayerFlags sets the layer flags of the entity.
func (ent *Entity) SetLayerFlags(flags math.BitFlag) {
	ent.layerFlags = flags
}
