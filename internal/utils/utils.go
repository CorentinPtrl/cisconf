package utils

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func DiffFields(src any, dest any) ([]string, error) {
	var r DiffReporter
	cmp.Diff(src, dest, cmp.Reporter(&r))
	return r.Fields(), nil
}

func GetTag(t reflect.Type, tag string) string {
	if t.Kind() == reflect.Slice {
		return GetTag(t.Elem(), tag)
	}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get(tag) != "" {
			return t.Field(i).Tag.Get(tag)
		}
	}
	return ""
}

func GetValueAndField(data interface{}, path string) (*reflect.Value, *reflect.StructField, error) {
	path = strings.TrimPrefix(path, ".")

	parts := strings.Split(path, ".")

	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	var structField *reflect.StructField

	for partIndex, part := range parts {
		re := regexp.MustCompile(`^([a-zA-Z0-9_]+)?\[(\d+)\]$`)
		matches := re.FindStringSubmatch(part)

		if matches != nil {
			fieldName := matches[1]
			index, _ := strconv.Atoi(matches[2])

			if fieldName != "" {
				/*defer func() {
					strName := fieldName
					valn := val
					pathn := path
					fmt.Println("defer fieldName", strName)
					fmt.Println("defer val", valn.Interface())
					fmt.Println("defer path", pathn)
					if r := recover(); r != nil {
						fmt.Println("recovered in GetValueAndField", r)
					}
				}()*/
				field, found := typ.FieldByName(fieldName)
				if !found {
					return nil, nil, fmt.Errorf("field %s not found", fieldName)
				}
				structField = &field
				val = val.FieldByName(fieldName)
				typ = val.Type()
			}

			if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
				return nil, nil, fmt.Errorf("field %s is not a slice or array", part)
			}

			if index < 0 || index >= val.Len() {
				return nil, nil, fmt.Errorf("index %d out of range", index)
			}

			if partIndex != len(parts)-1 {
				val = val.Index(index)
				typ = val.Type()
			} else {
				oldVal := val
				val = reflect.MakeSlice(reflect.SliceOf(val.Type().Elem()), 1, 1)
				val.Index(0).Set(oldVal.Index(index))
			}

		} else {
			field, found := typ.FieldByName(part)
			if !found {
				return nil, nil, fmt.Errorf("field %s not found", part)
			}

			structField = &field
			val = val.FieldByName(part)
			typ = val.Type()
		}

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
			typ = val.Type()
		}

		if !val.IsValid() {
			return nil, nil, fmt.Errorf("invalid path: %s", part)
		}
	}

	return &val, structField, nil
}
