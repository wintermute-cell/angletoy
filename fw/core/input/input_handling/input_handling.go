package input

import (
	input "gorl/fw/core/input/input_event"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// InputReceiver is an interface that should be implemented by any object that
// is intended to receive input events.
type InputReceiver interface {
	OnInputEvent(event *input.InputEvent) bool
}

// HandleInputEvents checks for input events and propagates them to the entities.
// Receives a sorted slice of layers, each containing a slice of entities.
// Both must be sorted from back to front (from far away to close to camera).
func HandleInputEvents(inputReceivers []InputReceiver) {
	// TODO: since therea re no more layers, ths makes no sense. rewrite.
	events := checkForInputs()
	for _, event := range events {
		// walk backwards so that the front-most entities receive the input first
		for i := len(inputReceivers) - 1; i >= 0; i-- {
			shouldContinue := inputReceivers[i].OnInputEvent(event)
			if !shouldContinue {
				break
			}
		}
	}
}

func checkForInputs() []*input.InputEvent {

	events := []*input.InputEvent{}
	mousePosition := rl.GetMousePosition()

	for action, triggers := range input.ActionMap {
		for _, trigger := range triggers {
			switch trigger.InputType {
			case input.InputTypeKey:
				switch trigger.TriggerType {
				case input.TriggerTypeDown:
					if rl.IsKeyDown(trigger.Key) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				case input.TriggerTypePressed:
					if rl.IsKeyPressed(trigger.Key) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				case input.TriggerTypeReleased:
					if rl.IsKeyReleased(trigger.Key) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				}
			case input.InputTypeMouse:
				switch trigger.TriggerType {
				case input.TriggerTypeDown:
					if rl.IsMouseButtonDown(trigger.MouseButton) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				case input.TriggerTypePressed:
					if rl.IsMouseButtonPressed(trigger.MouseButton) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				case input.TriggerTypeReleased:
					if rl.IsMouseButtonReleased(trigger.MouseButton) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				case input.TriggerTypePassive:
					events = append(events, input.NewInputEvent(action, mousePosition))
				}
			case input.InputTypeGamepad:
				// Implement the checks for gamepad buttons using a similar pattern
			}
		}
	}

	return events
}
