package datastructures

type Stack[T any] struct {
	items []T
}

// NewStack creates a new empty stack.
//
// prealloc specifies the preallocated capacity of the stack.
func NewStack[T any](prealloc int) *Stack[T] {
	return &Stack[T]{items: make([]T, 0, prealloc)}
}

// Push adds a new element to the top of the stack.
func (stack *Stack[T]) Push(key T) {
	stack.items = append(stack.items, key)
}

// Peek returns the top element of the stack without removing it.
// If the stack is empty, the second return value is false.
func (stack *Stack[T]) Peek() (T, bool) {
	var ret T
	if len(stack.items) > 0 {
		ret = stack.items[len(stack.items)-1]
		return ret, true
	}
	return ret, false
}

// Pop removes and returns the top element of the stack.
// If the stack is empty, the second return value is false.
func (stack *Stack[T]) Pop() (T, bool) {
	var ret T
	if len(stack.items) > 0 {
		ret = stack.items[len(stack.items)-1]
		stack.items = stack.items[:len(stack.items)-1]
		return ret, true
	}
	return ret, false
}

// IsEmpty returns true if the stack is empty.
func (stack *Stack[T]) IsEmpty() bool {
	return len(stack.items) == 0
}
