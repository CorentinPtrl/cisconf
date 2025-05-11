package cisconf

import (
	"fmt"
	"github.com/CorentinPtrl/cisconf/internal"
	"reflect"
)

// Unmarshal parses the given config into the provided struct pointer.
// It expects the struct to be a pointer to a struct type.
// If the provided value is not a pointer to a struct, it will return an error.
func Unmarshal(data string, v any) error {
	vValue := reflect.ValueOf(v)

	if vValue.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer to struct, got %s", vValue.Kind())
	}
	return internal.Parse(data, v)
}

// Marshal generates a config string from the provided struct.
// It can handle both struct and pointer to struct types.
func Marshal(v any) (string, error) {
	vValue := reflect.ValueOf(v)

	if vValue.Kind() == reflect.Ptr {
		vValue = vValue.Elem()
	}

	return internal.Generate(vValue.Interface(), vValue.Interface())
}

// Diff compares two structs of the same type and returns a string representation of the differences.
// It expects both src and dest to be of the same type.
// If they are not, it will return an error.
func Diff(src any, dest any) (string, error) {
	if src == nil || dest == nil {
		return "", fmt.Errorf("src or dest is nil")
	}

	if reflect.TypeOf(src) != reflect.TypeOf(dest) {
		return "", fmt.Errorf("src and dest must be of the same type")
	}

	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)

	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
		destValue = destValue.Elem()
	}

	return internal.Generate(srcValue.Interface(), destValue.Interface())
}
