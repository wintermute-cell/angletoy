package store

import (
	"reflect"
)

// store holds values keyed by their type.
type store struct {
	data map[reflect.Type]interface{}
}

// globalStore is the default store instance.
var globalStore = newStore()

// newStore creates a new instance of Store.
func newStore() *store {
	s := &store{
		data: make(map[reflect.Type]interface{}),
	}

	// Add premade stored types.
	// These should always be present in the store.
	addPremade(NewAppState(), s)

	return s
}

// Add adds or replaces a value in the store keyed by its type.
func Add[T any](value T) {
	t := reflect.TypeOf(value)
	globalStore.data[t] = value
}

// Get retrieves a value from the store by its type. It returns the value and a
// boolean indicating if it was found.
func Get[T any]() (T, bool) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	v, ok := globalStore.data[t]
	if !ok {
		return *new(T), false
	}
	return v.(T), true
}

func addPremade[T any](value T, s *store) {
	t := reflect.TypeOf(value)
	s.data[t] = value
}
