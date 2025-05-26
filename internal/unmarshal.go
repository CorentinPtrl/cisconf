package internal

import (
	"fmt"
	"github.com/CorentinPtrl/cisconf/internal/utils"
	"github.com/mcuadros/go-defaults"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Parse(data string, v any) error {
	defaults.SetDefaults(v)
	vValue := reflect.ValueOf(v)
	t := reflect.TypeOf(v).Elem()
	re := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	fullConfigNoNew := re.ReplaceAllString(data, "")

	re = regexp.MustCompile(`(?m)^!`)
	splitRun := re.Split(fullConfigNoNew, -1)

	for _, part := range splitRun {
		if utils.GetTag(t, "reg") != "" {
			if err := ProcessParse(part, v); err != nil {
				return err
			}
		}
		err := UnmarshalPart(t, vValue, part)
		if err != nil {
			return err
		}
	}
	return nil
}

func parsePart(part string, field reflect.Value, structField reflect.StructField) error {
	switch structField.Type.Kind() {
	case reflect.Struct:
		fieldValue := reflect.New(structField.Type).Elem()
		if err := Parse(part, fieldValue.Interface()); err != nil {
			return err
		}
		field.Set(fieldValue)
	case reflect.Slice:
		if field.IsNil() {
			field.Set(reflect.MakeSlice(reflect.SliceOf(structField.Type.Elem()), 0, 0))
		}
		fieldValue := reflect.New(structField.Type.Elem()).Elem()
		if err := Parse(part, fieldValue.Addr().Interface()); err != nil {
			return err
		}
		field.Set(reflect.Append(field, fieldValue))
	default:
		return fmt.Errorf("unsupported type: %s", structField.Type)
	}
	return nil
}

func UnmarshalPart(t reflect.Type, vValue reflect.Value, part string) error {
	fistLineArr := strings.Split(part, "\n")
	if len(fistLineArr) == 1 {
		return nil
	}
	for i := 0; i < t.NumField(); i++ {
		for k := 0; k < len(fistLineArr); k++ {
			firstLine := fistLineArr[k]
			preg := t.Field(i).Tag.Get("preg")
			if preg == "" {
				continue
			}
			re, err := regexp.Compile(preg)
			if err != nil {
				return err
			}
			if !re.MatchString(firstLine) {
				continue
			}
			if t.Field(i).Tag.Get("preg") != "" && utils.GetTag(reflect.TypeOf(vValue.Elem().Field(i).Interface()), "exit") != "" {
				re, _ := regexp.Compile(t.Field(i).Tag.Get("preg") + "\\n(?:\\s+.*\\n)*?\\s*" + utils.GetTag(reflect.TypeOf(vValue.Elem().Field(i).Interface()), "exit"))
				parts := re.FindAllString(part, -1)
				for _, subPart := range parts {
					field := vValue.Elem().Field(i)
					if err := parsePart(subPart, field, t.Field(i)); err != nil {
						return err
					}
				}
				break
			} else {
				field := vValue.Elem().Field(i)
				err = parsePart(part, field, t.Field(i))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ProcessParse takes a part of the config and a pointer to a struct with tags on how to parse the config part.
// ProcessParse can also take a reflect.Value to call itself.
func ProcessParse(part string, parsed any) error {
	// Check if type is already a "reflect.Value", to let the function call itself in case of a struct in a struct
	var tmp, rv reflect.Value
	// In the normal case the function takes a struct.
	//In case of a struct as parsed parameter it creates a copy of the struct to write to the copy of its reflection
	// When a reflect.Value is already used as parsed parameter this is not needed and cause an error.
	// The Value still needs to be a reflect of a struct value to function. This is needed to work with recursion e.g.
	// structs in structs.
	// I don't really now how I got all this reflect stuff together and working so note to further self:
	// Make this more stable as soon I understand this shit...
	if reflect.TypeOf(parsed).String() != "reflect.Value" {
		// Generate a copy of the "parsed" interface to fill it with values
		v := reflect.Indirect(reflect.ValueOf(&parsed)).Elem()
		// tmp is the copy and can be written to.
		// parsed, rv and rt are read only. parsed can be completely overwritten but no values can be set on by one
		tmp = reflect.New(v.Elem().Type()).Elem()
		tmp.Set(v.Elem())

		rv = reflect.Indirect(reflect.ValueOf(parsed))
	} else {
		// write any value as reflect.Value to all the copied and real value
		// in this case parsed is writeable and no copy needed
		tmp, rv = parsed.(reflect.Value), parsed.(reflect.Value)
	}
	rt := rv.Type()

	// for through all field of the struct, get the regex tag and fill it with the found data
	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get("reg")

		if field.Tag.Get("reg") == "" && field.Tag.Get("preg") == "" {
			if field.Type.Kind() != reflect.Struct {
				continue
			}
			err := ProcessParse(part, tmp.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
		}

		if tag != "" {
			re := regexp.MustCompile("(?:no\\s+)?" + tag)
			// @todo check if no is with the command!
			data := re.FindAllStringSubmatch(part, -1)
			if len(data) == 0 {
				continue
			}

			// Check the Line of the struct field
			switch field.Type.Kind() {
			case reflect.String:
				tmp.Field(i).SetString(data[0][1])
			case reflect.Int:
				value, err := strconv.ParseInt(data[0][1], 10, 64)
				if err != nil {
					return err
				}
				tmp.Field(i).SetInt(value)
			case reflect.Bool:
				if strings.HasPrefix(data[0][0], "no ") {
					tmp.Field(i).SetBool(false)
					break
				}
				tmp.Field(i).SetBool(true)
			case reflect.Float64:
				float, err := strconv.ParseFloat(data[0][1], 64)
				if err != nil {
					return err
				}
				tmp.Field(i).SetFloat(float)
			case reflect.Struct:
				fieldValue := reflect.New(field.Type).Elem()
				if err := ProcessParse(part, fieldValue.Addr().Interface()); err != nil {
					return err
				}
				tmp.Field(i).Set(fieldValue)
			case reflect.Slice:
				switch field.Type.String() {
				case "[]string":
					for _, d := range data {
						value := tmp.Field(i)
						value.Set(reflect.Append(value, reflect.ValueOf(d[1])))
					}
				case "[]int":
					value := tmp.Field(i)
					for _, d := range data {
						separated := strings.Split(d[2], ",")
						for _, number := range separated {
							if strings.Contains(number, "-") {
								vlanSplit := strings.Split(number, "-")
								from, _ := strconv.Atoi(vlanSplit[0])
								to, _ := strconv.Atoi(vlanSplit[1])
								for j := from; j <= to; j++ {
									value.Set(reflect.Append(value, reflect.ValueOf(j)))
								}
								continue
							}
							vlanI, _ := strconv.Atoi(number)
							value.Set(reflect.Append(value, reflect.ValueOf(vlanI)))
						}
					}
				default:
					values := tmp.Field(i)
					if values.Type().Kind() == reflect.Slice {
						for _, t := range data {
							tmp2 := reflect.New(values.Type().Elem()).Elem()
							err := Parse(strings.Join(t, " "), tmp2.Addr().Interface())
							if err != nil {
								return err
							}
							values.Set(reflect.Append(values, tmp2))
						}
						continue
					}
					panic(field.Type.String() + " not implemented!")
				}
			default:
				panic(field.Type.String() + " not implemented!")
			}
		}
	}

	// Overwrite parsed with tmp to get the values back to the caller.
	// In case parsed is a reflect.Value it needs to write into the element and not the Value itself
	if reflect.TypeOf(parsed).String() != "reflect.Value" {
		reflect.ValueOf(parsed).Elem().Set(tmp)
		return nil
	}

	parsed.(reflect.Value).Set(tmp)
	return nil
}
