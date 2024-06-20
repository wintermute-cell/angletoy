package scenes

import (
	"gorl/fw/core/entities"
)

// IScene is an interface that every scene in the game should implement.
type IScene interface {
	// Init and Deinit are implemented by the user.
	Init()
	Deinit()

	// GetRoot is provided by the base implementation.
	GetRoot() *entities.Entity
}
