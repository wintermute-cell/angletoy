package scenes

import (
	"gorl/fw/core/entities"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scene struct {
	rootNode *entities.Entity
}

// GetRoot returns the root node of the scene.
func (s *Scene) GetRoot() *entities.Entity {
	if s.rootNode == nil {
		s.rootNode = entities.NewEntity("scene_root", rl.Vector2Zero(), 0, rl.Vector2One())
	}
	return s.rootNode
}
