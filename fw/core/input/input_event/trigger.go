package input

// TriggerType defines the type of event, e.g. down, pressed, released.
type TriggerType int32

const (
	TriggerTypeDown TriggerType = iota
	TriggerTypePressed
	TriggerTypeReleased
	TriggerTypePassive // Passive triggers are always active
)

// InputType defines the physical type of trigger. It can be a key, mouse
// button, or gamepad button.
type InputType int32

const (
	InputTypeKey InputType = iota
	InputTypeMouse
	InputTypeGamepad
)

// A Trigger is a definition of an input trigger that can cause an action. We
// use it to map specific triggers to abstract actions.
type Trigger struct {
	InputType   InputType
	TriggerType TriggerType
	Key         int32
	MouseButton int32
}
