package input

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// An InputEvent describes any user input event.
// It has a TriggerType which describes what kind of trigger caused the event,
// and determines what data the event has.
type InputEvent struct {
	Action         Action
	cursorPosition rl.Vector2
	// TODO: there might be more fields necessary here, especially to support
	// other input devices such as gamepads.
}

func NewInputEvent(action Action, cursorPosition rl.Vector2) *InputEvent {
	return &InputEvent{Action: action, cursorPosition: cursorPosition}
}

func (e *InputEvent) GetScreenSpaceMousePosition() rl.Vector2 {
	// FIXME: this is called during update and not during render, so there is not active render stage.
	// Approaches:
	// 1. make sure this function is called during render.
	// 2. have only a single camera config for all stages, so we don't need to know the active stage.
	// 3. have each entity know its render stage. this could be accomplished by "learning", e.g. every time the entity is drawn, it learns the current render stage.
	return e.cursorPosition
}
