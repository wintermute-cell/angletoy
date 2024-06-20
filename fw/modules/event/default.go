package event

// Default event dispatcher instance
var defaultDispatcher = NewDispatcher()

// Listen registers a new event listener to the given event name.
// (Uses the default dispatcher)
func Listen(name string, fn EventHandler) error {
	return defaultDispatcher.Listen(name, fn)
}

// Trigger fires an event by name and passes the given parameters to the listeners
// (Uses the default dispatcher)
func Trigger(name string, params ...interface{}) error {
	return defaultDispatcher.Trigger(name, params...)
}

// HasEvent returns true if a event with the given name exists.
// (Uses the default dispatcher)
func HasEvent(name string) bool {
	return defaultDispatcher.HasEvent(name)
}

// ListEvents returns a list of all registered events
// (Uses the default dispatcher)
func ListEvents() []string {
	return defaultDispatcher.ListEvents()
}

// RemoveEvents deletes events from the event list.
// (Uses the default dispatcher)
func RemoveEvents(names ...string) {
	defaultDispatcher.RemoveEvents(names...)
}
