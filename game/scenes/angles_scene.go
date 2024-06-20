package scenes

import (
	"gorl/fw/core/gem"
	"gorl/fw/modules/scenes"
	"gorl/game/entities"
)

// This checks at compile time if the interface is implemented
var _ scenes.IScene = (*AnglesScene)(nil)

// Angles Scene
type AnglesScene struct {
	scenes.Scene // Required!

	// Custom Fields
	// Add fields here for any state that the scene should keep track of
	// ...
}

func (scn *AnglesScene) Init() {
	cam := entities.NewCameraEntity()
	gem.Append(scn.GetRoot(), cam)

	showcaser := entities.NewAngleShowcaserEntity()
	gem.Append(scn.GetRoot(), showcaser)
}

func (scn *AnglesScene) Deinit() {
	// De-initialization logic for the scene
	// ...
}
