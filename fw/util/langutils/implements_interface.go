package langutils

import (
	"fmt"
	"reflect"
)

// FIXME: this function is pretty useless, since only checking for an interface and not just a type is very limited,
// and it interacts badly with the way we have to handle fetching entities from the Gem.
// We should probably just remove this function and let the user handle filtering entities themselves.
//
// ImplementsInterface checks if the provided struct pointer implements the interface T.
// T should be an interface type, and structPtr should be a pointer to a struct.
func ImplementsInterface[T any](structPtr any) bool {
	structType := reflect.TypeOf(structPtr)
	if structType.Kind() != reflect.Ptr {
		fmt.Println("Error: structPtr must be a pointer to a struct")
		return false
	}
	// Get the type of the interface from the generic type T.
	interfaceType := reflect.TypeOf((*T)(nil)).Elem()
	// Check if the struct type implements the interface type.
	return structType.Implements(interfaceType)
}
