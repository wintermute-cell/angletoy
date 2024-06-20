package entities

import (
	"gorl/fw/core/datastructures"
	"gorl/fw/core/entities"
	"gorl/fw/core/gem"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/math"
	"gorl/fw/core/render"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that CameraEntity implements IEntity.
var _ entities.IEntity = &CameraEntity{}

// Camera Entity
type CameraEntity struct {
	*entities.Entity
	camera *render.Camera
	ctb    *cameraTransformationBuffer
}

func NewCameraEntity(
	camTarget, camOffset,
	displaySize, displayPosition rl.Vector2,
	drawFlags math.BitFlag,
) *CameraEntity {
	new_ent := &CameraEntity{
		Entity: entities.NewEntity("CameraEntity", camTarget, 0, rl.Vector2One()),
		camera: render.NewCamera(
			camTarget,
			camOffset,
			displaySize,
			displayPosition,
			drawFlags,
		),
		ctb: &cameraTransformationBuffer{},
	}
	return new_ent
}

// ============================================================================
// Utilities
// ============================================================================

// ScreenToWorld converts a screen position to a world position.
func (ent *CameraEntity) ScreenToWorld(screenPos rl.Vector2) rl.Vector2 {
	return ent.camera.ScreenToWorld(screenPos)
}

// WorldToScreen converts a world position to a screen position.
func (ent *CameraEntity) WorldToScreen(worldPos rl.Vector2) rl.Vector2 {
	return ent.camera.WorldToScreen(worldPos)
}

// ============================================================================
// IEntity
// ============================================================================

func (ent *CameraEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *CameraEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *CameraEntity) Update() {

	// 1. Apply the absolute transform of the camera entity to the render camera.
	absTransform := gem.GetAbsoluteTransform(ent)
	//ent.camera.SetTarget(absTransform.GetPosition())
	//ent.camera.SetRotation(absTransform.GetRotation())
	//ent.camera.SetZoom(absTransform.GetScale().X)
	ent.ctb.Position = datastructures.NewMaybe(absTransform.GetPosition())
	ent.ctb.Rotation = datastructures.NewMaybe(absTransform.GetRotation())
	ent.ctb.Zoom = datastructures.NewMaybe(absTransform.GetScale().X)

	// 2. Apply the cameraTransformationBuffer on top of that.
	ent.ctb.flushToCamera(ent.camera)

}

func (ent *CameraEntity) Draw() {
	// Draw logic for the entity
	// ...
}

func (ent *CameraEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...

	const moveSpeed = 100
	const zoomSpeed = 0.3

	if event.Action == input.ActionMoveLeft {
		ent.SetPosition(rl.Vector2Add(ent.GetPosition(), rl.NewVector2(-moveSpeed*rl.GetFrameTime(), 0)))
	}
	if event.Action == input.ActionMoveRight {
		ent.SetPosition(rl.Vector2Add(ent.GetPosition(), rl.NewVector2(moveSpeed*rl.GetFrameTime(), 0)))
	}
	if event.Action == input.ActionMoveUp {
		ent.SetPosition(rl.Vector2Add(ent.GetPosition(), rl.NewVector2(0, -moveSpeed*rl.GetFrameTime())))
	}
	if event.Action == input.ActionMoveDown {
		ent.SetPosition(rl.Vector2Add(ent.GetPosition(), rl.NewVector2(0, moveSpeed*rl.GetFrameTime())))
	}
	if event.Action == input.ActionZoomIn {
		ent.SetScale(rl.NewVector2(ent.GetScale().X+zoomSpeed*rl.GetFrameTime(), 1))
	}
	if event.Action == input.ActionZoomOut {
		ent.SetScale(rl.NewVector2(ent.GetScale().X-zoomSpeed*rl.GetFrameTime(), 1))
	}

	return true
}

// ============================================================================
// Camera Transformation Buffer
// ============================================================================

// cameraTransformationBuffer stores the transformation data for the camera
// within one frame. this allows us to compose a final transformation from
// multiple sources such as target tracking, screen shake, etc...
type cameraTransformationBuffer struct {
	Position       datastructures.Maybe[rl.Vector2]
	PositionChange []rl.Vector2
	Offset         datastructures.Maybe[rl.Vector2]
	OffsetChange   []rl.Vector2
	Rotation       datastructures.Maybe[float32]
	RotationChange []float32
	Zoom           datastructures.Maybe[float32]
	ZoomChange     []float32
}

// reset clears the transformation buffer without reallocation
func (ctb *cameraTransformationBuffer) reset() {
	ctb.Position.Unset()
	ctb.PositionChange = ctb.PositionChange[:0]
	ctb.Offset.Unset()
	ctb.OffsetChange = ctb.OffsetChange[:0]
	ctb.Rotation.Unset()
	ctb.RotationChange = ctb.RotationChange[:0]
	ctb.Zoom.Unset()
	ctb.ZoomChange = ctb.ZoomChange[:0]
}

func (ctb *cameraTransformationBuffer) flushToCamera(camera *render.Camera) {
	if position, ok := ctb.Position.Get(); ok {
		camera.SetTarget(position)
	}
	for _, positionChange := range ctb.PositionChange {
		camera.SetTarget(rl.Vector2Add(camera.GetTarget(), positionChange))
	}
	if offset, ok := ctb.Offset.Get(); ok {
		camera.SetOffset(offset)
	}
	for _, offsetChange := range ctb.OffsetChange {
		camera.SetOffset(rl.Vector2Add(camera.GetOffset(), offsetChange))
	}
	if rotation, ok := ctb.Rotation.Get(); ok {
		camera.SetRotation(rotation)
	}
	for _, rotationChange := range ctb.RotationChange {
		camera.SetRotation(camera.GetRotation() + rotationChange)
	}
	if zoom, ok := ctb.Zoom.Get(); ok {
		camera.SetZoom(zoom)
	}
	for _, zoomChange := range ctb.ZoomChange {
		camera.SetZoom(camera.GetZoom() + zoomChange)
	}
	ctb.reset()
}
