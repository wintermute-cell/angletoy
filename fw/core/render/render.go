package render

import (
	input "gorl/fw/core/input/input_handling"
	"gorl/fw/core/math"
	"gorl/game/code/colorscheme"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Drawable is an interface that should be implemented by any object that
// is intended to be rendered.
type Drawable interface {
	ShouldDraw(layerFlags math.BitFlag) bool
	Draw()
	GetDrawIndex() int32
	AsInputReceiver() input.InputReceiver
}

// renderer is a set of cameras and a final render target that is used to
// render all drawables.
type renderer struct {
	cameras     []*Camera
	finalTarget rl.RenderTexture2D
}

// Init initializes the renderer with the given screen size.
func Init(screenSize rl.Vector2) {
	rendererInstance = renderer{
		cameras: []*Camera{},
		finalTarget: rl.LoadRenderTexture(
			int32(screenSize.X),
			int32(screenSize.Y),
		),
	}
}

// Deinit deinitializes the renderer.
func Deinit() {
	rl.UnloadRenderTexture(rendererInstance.finalTarget)
}

// SetScreenSize changes the size of the screen.
func SetScreenSize(screenSize rl.Vector2) {
	rl.UnloadRenderTexture(rendererInstance.finalTarget)
	rendererInstance.finalTarget = rl.LoadRenderTexture(
		int32(screenSize.X),
		int32(screenSize.Y),
	)
}

// rendererInstance is the global renderer instance.
var rendererInstance renderer

// Draw draws the given drawable to the screen, using all cameras.
// Returns the sorted list of drawables, back to front.
func Draw(drawables []Drawable) []input.InputReceiver {

	inputReceivers := []input.InputReceiver{}

	// sort the drawables by draw index
	slices.SortStableFunc(drawables, func(l, r Drawable) int {
		return int(l.GetDrawIndex() - r.GetDrawIndex())
	})

	for _, camera := range rendererInstance.cameras {
		rl.BeginTextureMode(camera.renderTarget.renderTexture)
		rl.BeginMode2D(*camera.rlcamera)
		rl.ClearBackground(rl.Blank)

		// Draw all drawables that should be drawn by this camera.
		for _, drawable := range drawables {
			if drawable.ShouldDraw(camera.drawFlags) {
				inputReceivers = append(inputReceivers, drawable.AsInputReceiver())
				drawable.Draw()
			}
		}

		rl.EndMode2D()
		rl.EndTextureMode()
	}

	// Draw all camera render targets to the final target.
	// Apply per camera shaders in the process.
	rl.BeginTextureMode(rendererInstance.finalTarget)
	rl.ClearBackground(colorscheme.Colorscheme.Color16.ToRGBA())
	for _, camera := range rendererInstance.cameras {
		applyShaders(camera)
		rl.DrawTexturePro(
			camera.renderTarget.renderTexture.Texture,
			rl.NewRectangle(0, 0, float32(camera.renderTarget.renderTexture.Texture.Width), -float32(camera.renderTarget.renderTexture.Texture.Height)),
			rl.NewRectangle(
				camera.renderTarget.DisplayPosition.X,
				camera.renderTarget.DisplayPosition.Y,
				camera.renderTarget.DisplaySize.X,
				camera.renderTarget.DisplaySize.Y,
			),
			rl.NewVector2(0, 0),
			0, rl.White,
		)
	}
	rl.EndTextureMode()

	// Draw the final target to the screen.
	// TODO: here we could apply final post processing shaders.

	rl.DrawTexturePro(
		rendererInstance.finalTarget.Texture,
		rl.NewRectangle(0, 0, float32(rendererInstance.finalTarget.Texture.Width), -float32(rendererInstance.finalTarget.Texture.Height)),
		rl.NewRectangle(0, 0, float32(rendererInstance.finalTarget.Texture.Width), float32(rendererInstance.finalTarget.Texture.Height)),
		rl.NewVector2(0, 0),
		0, rl.White,
	)

	return inputReceivers
}

// ApplyShaders applies the shaders of the camera to the cameras render target.
func applyShaders(camera *Camera) {
	currentSource := &camera.renderTarget.renderTexture
	currentTarget := &camera.bounceTexture

	for _, shader := range camera.shaders {
		rl.BeginTextureMode(*currentTarget)
		rl.BeginShaderMode(*shader)
		rl.DrawTexturePro(
			currentSource.Texture,
			rl.NewRectangle(0, 0, float32(currentSource.Texture.Width), -float32(currentSource.Texture.Height)),
			rl.NewRectangle(0, 0, float32(currentTarget.Texture.Width), float32(currentTarget.Texture.Height)),
			rl.NewVector2(0, 0),
			0, rl.White,
		)
		rl.EndShaderMode()
		rl.EndTextureMode()
		tmp := currentSource
		currentSource = currentTarget
		currentTarget = tmp
	}

	// if the last draw was to the bounce texture, draw it to the cameras target.
	if currentSource == &camera.bounceTexture {
		rl.DrawTexturePro(
			currentSource.Texture,
			rl.NewRectangle(0, 0, float32(currentSource.Texture.Width), -float32(currentSource.Texture.Height)),
			rl.NewRectangle(0, 0, float32(currentSource.Texture.Width), float32(currentSource.Texture.Height)),
			rl.NewVector2(0, 0),
			0, rl.White,
		)
	}
}
