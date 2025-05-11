package cisconf

import (
	"fmt"
	"github.com/CorentinPtrl/cisconf/internal"
	"reflect"
)

func Unmarshal(data string, v any) error {
	vValue := reflect.ValueOf(v)

	if vValue.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer to struct, got %s", vValue.Kind())
	}
	return internal.Parse(data, v)
}

func Marshal(v any) (string, error) {
	vValue := reflect.ValueOf(v)

	if vValue.Kind() == reflect.Ptr {
		vValue = vValue.Elem()
	}

	return internal.Generate(vValue.Interface(), vValue.Interface())
}

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
