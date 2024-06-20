package event

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// dispatcher is a dispatcher with a map of events and corresponding listeners
type dispatcher struct {
	sync.RWMutex

	events map[string][]EventHandler
}

// EventHandler is a function that can be registered as an event listener.
// It must be a function that returns an error and nothing else.
type EventHandler any

type Dispatcher interface {
	Listen(name string, fn EventHandler) error
	Trigger(name string, params ...any) error
	HasEvent(name string) bool
	ListEvents() []string
	RemoveEvents(names ...string)
}

// NewDispatcher returns a new event dispatcher.
//
// Use this to create a custom dispatcher. Otherwise use the default dispatcher
// by directly calling the functions on the event package.
func NewDispatcher() Dispatcher {
	return &dispatcher{
		events: make(map[string][]EventHandler),
	}
}

// Listen registers a new event listener to the given event name
func (e *dispatcher) Listen(name string, fn EventHandler) error {
	e.Lock()
	defer e.Unlock()

	if fn == nil {
		return errors.New("the function (fn) provided is nil")
	}

	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		return errors.New("the provided function (fn) is not of type 'func'")
	}
	if t.NumOut() != 1 || t.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("the function (fn) must have exactly one return value of type error; got %d return(s): %s, ...", t.NumOut(), t.Out(0))
	}
	if list, ok := e.events[name]; ok && len(list) > 0 {
		tt := reflect.TypeOf(list[0])
		if tt.NumIn() != t.NumIn() {
			return fmt.Errorf("function (fn) signature mismatch, expected %d parameters, got %d", tt.NumIn(), t.NumIn())
		}
		for i := 0; i < tt.NumIn(); i++ {
			if tt.In(i) != t.In(i) {
				return fmt.Errorf("function (fn) parameter type mismatch at position %d: expected %s, got %s", i+1, tt.In(i), t.In(i))
			}
		}
	}

	e.events[name] = append(e.events[name], fn)
	return nil
}

// Trigger fires an event by name and passes the given parameters to the listeners
func (e *dispatcher) Trigger(name string, params ...any) error {
	e.RLock()
	defer e.RUnlock()

	fns := e.events[name]
	for i := len(fns) - 1; i >= 0; i-- {
		stopped, err := e.call(fns[i], params...)
		if err != nil {
			return err
		}
		if stopped {
			break
		}
	}

	return nil
}

func (e *dispatcher) call(fn EventHandler, params ...any) (stopped bool, err error) {
	var (
		f     = reflect.ValueOf(fn)
		t     = f.Type()
		numIn = t.NumIn()
		in    = make([]reflect.Value, 0, numIn)
	)

	if t.IsVariadic() {
		n := numIn - 1
		if len(params) < n {
			return stopped, fmt.Errorf("insufficient parameters for variadic function: expected at least %d, got %d; expected types: %v", n, len(params), paramTypes(t, n))
		}
		for _, param := range params[:n] {
			in = append(in, reflect.ValueOf(param))
		}
		s := reflect.MakeSlice(t.In(n), 0, len(params[n:]))
		for _, param := range params[n:] {
			s = reflect.Append(s, reflect.ValueOf(param))
		}
		in = append(in, s)

		result := f.CallSlice(in)[0].Interface()
		if err, ok := result.(error); ok && err != nil {
			return stopped, err
		}
		return stopped, nil
	}

	if len(params) != numIn {
		return stopped, fmt.Errorf("parameter count mismatch: expected %d, got %d; expected types: %v, received types: %v", numIn, len(params), paramTypes(t, numIn), receivedParamTypes(params))
	}
	for _, param := range params {
		in = append(in, reflect.ValueOf(param))
	}

	result := f.Call(in)[0].Interface()
	if err, ok := result.(error); ok && err != nil {
		return stopped, err
	}
	return stopped, nil
}

// Helper function to list expected parameter types
func paramTypes(t reflect.Type, numIn int) []string {
	types := make([]string, numIn)
	for i := 0; i < numIn; i++ {
		types[i] = t.In(i).String()
	}
	return types
}

// Helper function to list received parameter types
func receivedParamTypes(params []any) []string {
	types := make([]string, len(params))
	for i, p := range params {
		types[i] = reflect.TypeOf(p).String()
	}
	return types
}

// HasEvent returns true if an event with the given name exists .
func (e *dispatcher) HasEvent(name string) bool {
	e.RLock()
	defer e.RUnlock()
	_, ok := e.events[name]
	return ok
}

// ListEvents returns a list of all registered events.
func (e *dispatcher) ListEvents() []string {
	e.RLock()
	defer e.RUnlock()
	list := make([]string, 0, len(e.events))
	for name := range e.events {
		list = append(list, name)
	}
	return list
}

// RemoveEvents deletes events from the event list.
func (e *dispatcher) RemoveEvents(names ...string) {
	e.Lock()
	defer e.Unlock()
	if len(names) > 0 {
		for _, name := range names {
			delete(e.events, name)
		}
		return
	}
	e.events = make(map[string][]EventHandler)
}
