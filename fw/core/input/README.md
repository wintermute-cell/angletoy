# Core: input

TODO: this is a draft, expand.

- input_event defines triggers and maps them to an abstract event like "shoot" or "jump".

- input_handling checks for these events and passes them to the entities.

## Usage
- define action in action_map.go
- wait for action in entity.OnInputEvent(event) (if event.Action == input.MyAction)
- handle the action, for example bounds checking for cursor with rl.CheckCollision...(myShape, input.CursorPosition)
- return false to stop propagation or true to allow

## How does it work?
- input event defines types of physical triggers and maps a combination of triggers and keys to abstract actions.
- every frame, input_handling checks if any of these events have occurred. if so, it passes the fitting InputActions through all entities, in the order they were drawn in.
