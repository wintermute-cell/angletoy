package render

import (
	"gorl/fw/core/math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// renderTarget is a render texture that is drawn to the screen at a specific
// position and size.
type renderTarget struct {
	DisplayPosition rl.Vector2
	DisplaySize     rl.Vector2
	renderTexture   rl.RenderTexture2D
}

// Camera represents a raylib camera together with a render target, a set
// of draw flags that determine which drawables should be drawn by this camera
// and a set of shaders that should be applied to the render target.
type Camera struct {
	rlcamera     *rl.Camera2D
	renderTarget *renderTarget
	drawFlags    math.BitFlag
	shaders      []*rl.Shader

	// a render texture used when applying the shader stack.
	bounceTexture rl.RenderTexture2D
}

// NewCamera creates a new camera with the given target, offset, display size,
// display position and draw flags. The camera is added to the global renderer
// instance.
func NewCamera(camTarget, camOffset, displaySize, displayPosition rl.Vector2, drawFlags math.BitFlag) *Camera {
	rlCamera := rl.NewCamera2D(camOffset, camTarget, 0, 1)
	camera := &Camera{
		rlcamera:      &rlCamera,
		renderTarget:  &renderTarget{displayPosition, displaySize, rl.LoadRenderTexture(int32(displaySize.X), int32(displaySize.Y))},
		drawFlags:     drawFlags,
		shaders:       make([]*rl.Shader, 0),
		bounceTexture: rl.LoadRenderTexture(int32(displaySize.X), int32(displaySize.Y)),
	}
	rendererInstance.cameras = append(rendererInstance.cameras, camera)
	return camera
}

// Destroy destroys the camera and removes it from the global renderer instance.
func (c *Camera) Destroy() {
	for i, camera := range rendererInstance.cameras {
		if camera == c {
			rendererInstance.cameras = append(rendererInstance.cameras[:i], rendererInstance.cameras[i+1:]...)
			break
		}
	}
	rl.UnloadRenderTexture(c.renderTarget.renderTexture)
}

// ScreenToWorld converts a screen position to a world position.
func (c *Camera) ScreenToWorld(screenPos rl.Vector2) rl.Vector2 {
	return rl.GetScreenToWorld2D(screenPos, *c.rlcamera)
}

// WorldToScreen converts a world position to a screen position.
func (c *Camera) WorldToScreen(worldPos rl.Vector2) rl.Vector2 {
	return rl.GetWorldToScreen2D(worldPos, *c.rlcamera)
}

// SetTarget sets the target (position) of the camera.
func (c *Camera) SetTarget(target rl.Vector2) {
	c.rlcamera.Target = target
}

// GetTarget returns the target (position) of the camera.
func (c *Camera) GetTarget() rl.Vector2 {
	return c.rlcamera.Target
}

// SetOffset sets the offset of the camera.
func (c *Camera) SetOffset(offset rl.Vector2) {
	c.rlcamera.Offset = offset
}

// GetOffset returns the offset of the camera.
func (c *Camera) GetOffset() rl.Vector2 {
	return c.rlcamera.Offset
}

// SetRotation sets the rotation of the camera.
func (c *Camera) SetRotation(rotation float32) {
	c.rlcamera.Rotation = rotation
}

// GetRotation returns the rotation of the camera.
func (c *Camera) GetRotation() float32 {
	return c.rlcamera.Rotation
}

// SetZoom sets the zoom of the camera.
func (c *Camera) SetZoom(zoom float32) {
	c.rlcamera.Zoom = zoom
}

// GetZoom returns the zoom of the camera.
func (c *Camera) GetZoom() float32 {
	return c.rlcamera.Zoom
}

// SetDrawFlags sets the draw flags of the camera.
func (c *Camera) SetDrawFlags(drawFlags math.BitFlag) {
	c.drawFlags = drawFlags
}

// GetDrawFlags returns the draw flags of the camera.
func (c *Camera) GetDrawFlags() math.BitFlag {
	return c.drawFlags
}

// AddShader adds a shader to the camera.
func (c *Camera) AddShader(shader *rl.Shader) {
	c.shaders = append(c.shaders, shader)
}

// RemoveShader removes a shader from the camera.
func (c *Camera) RemoveShader(shader *rl.Shader) {
	for i, s := range c.shaders {
		if s == shader {
			c.shaders = append(c.shaders[:i], c.shaders[i+1:]...)
			break
		}
	}
}

//
//type cameraTransformationBuffer struct {
//	Position       []rl.Vector2
//	PositionChange []rl.Vector2
//	Offset         []rl.Vector2
//	OffsetChange   []rl.Vector2
//	Rotation       []float32
//	RotationChange []float32
//	Zoom           []float32
//	ZoomChange     []float32
//}
//
//// Resets all slices in cameraTransformationBuffer to empty without reallocation
//func (ctb *cameraTransformationBuffer) reset() {
//	ctb.Position = ctb.Position[:0]
//	ctb.PositionChange = ctb.PositionChange[:0]
//	ctb.Offset = ctb.Offset[:0]
//	ctb.OffsetChange = ctb.OffsetChange[:0]
//	ctb.Rotation = ctb.Rotation[:0]
//	ctb.RotationChange = ctb.RotationChange[:0]
//	ctb.Zoom = ctb.Zoom[:0]
//	ctb.ZoomChange = ctb.ZoomChange[:0]
//}
//
//var cameraTransformations cameraTransformationBuffer
//
//func applyCameraTransformations(camera *rl.Camera2D) {
//	for _, position := range cameraTransformations.Position {
//		camera.Target = position
//	}
//	for _, positionChange := range cameraTransformations.PositionChange {
//		camera.Target = rl.Vector2Add(camera.Target, positionChange)
//	}
//	for _, offset := range cameraTransformations.Offset {
//		camera.Offset = offset
//	}
//	for _, offsetChange := range cameraTransformations.OffsetChange {
//		camera.Offset = rl.Vector2Add(camera.Offset, offsetChange)
//	}
//	for _, rotation := range cameraTransformations.Rotation {
//		camera.Rotation = rotation
//	}
//	for _, rotationChange := range cameraTransformations.RotationChange {
//		camera.Rotation += rotationChange
//	}
//	for _, zoom := range cameraTransformations.Zoom {
//		camera.Zoom = zoom
//	}
//	for _, zoomChange := range cameraTransformations.ZoomChange {
//		camera.Zoom += zoomChange
//	}
//	cameraTransformations.reset()
//}
//
//// TODO: these transforms have to be applied to the correct camera / render stage
//func SetCameraPosition(position rl.Vector2) {
//	cameraTransformations.Position = append(cameraTransformations.Position, position)
//}
