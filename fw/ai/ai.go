package ai

import rl "github.com/gen2brain/raylib-go/raylib"

type AiTarget struct {
	Position    rl.Vector2
	IsCompleted bool
}

type AiController interface {
	Update()
	SetSteeringForce(new_force rl.Vector2)
	GetSteeringForce() rl.Vector2
	GetControllable() AiControllable
}

type AiControllable interface {
	SetAiTarget(target *AiTarget)
	SetAiSteeringForce(force rl.Vector2)
	GetPosition() rl.Vector2
	GetPlayerPosition() rl.Vector2
	CanSeePlayer() bool
	CanSeePoint(point rl.Vector2) bool
}
