package datastructures

// Maybe is a container that may or may not contain a value. Such a type is
// sometimes also called Option.
type Maybe[T any] struct {
	value T
	has   bool
}

// NewMaybe creates a new Maybe container with the given value.
func NewMaybe[T any](value T) Maybe[T] {
	return Maybe[T]{value, true}
}

// EmptyMaybe creates a new Maybe container without a value.
func EmptyMaybe[T any]() Maybe[T] {
	return Maybe[T]{has: false}
}

// Get returns the value and a bool indicating if the value is ok (true) or considered
// "empty" (false).
func (m *Maybe[T]) Get() (T, bool) {
	return m.value, m.has
}

// Has returns true if the Maybe contains a value, false otherwise.
func (m *Maybe[T]) Has() bool {
	return m.has
}

// Yank returns the value of the maybe without checking if the value is
// considered ok. The value might be nil or some other invalid value.
func (m *Maybe[T]) Yank() T {
	return m.value
}

// Unset marks the contained value as invalid.
func (m *Maybe[T]) Unset() {
	m.has = false
}
