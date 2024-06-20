package input

import rl "github.com/gen2brain/raylib-go/raylib"

// Actions
type Action int32

const (
	ActionMoveUp Action = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionClickDown
	ActionClickHeld
	ActionClickUp
	ActionMouseHover
	ActionEscape
	// Add other actions as needed
	ActionZoomIn
	ActionZoomOut
	ActionNextAnimation
)

var ActionMap = map[Action][]Trigger{
	ActionMoveUp: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyW},
	},
	ActionMoveDown: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyS},
	},
	ActionMoveLeft: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyA},
	},
	ActionMoveRight: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyD},
	},
	ActionClickDown: {
		{InputType: InputTypeMouse, TriggerType: TriggerTypePressed, MouseButton: rl.MouseLeftButton},
	},
	ActionClickHeld: {
		{InputType: InputTypeMouse, TriggerType: TriggerTypeDown, MouseButton: rl.MouseLeftButton},
	},
	ActionClickUp: {
		{InputType: InputTypeMouse, TriggerType: TriggerTypeReleased, MouseButton: rl.MouseLeftButton},
	},
	ActionMouseHover: {
		{InputType: InputTypeMouse, TriggerType: TriggerTypePassive},
	},
	ActionEscape: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyEscape},
	},

	ActionZoomIn: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyQ},
	},
	ActionZoomOut: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyE},
	},
	ActionNextAnimation: {
		{InputType: InputTypeKey, TriggerType: TriggerTypePressed, Key: rl.KeyN},
	},
	// Add other action-trigger mappings
}
