package internal

import (
	"fmt"
	"github.com/CorentinPtrl/cisconf/internal/utils"
	"reflect"
	"strconv"
	"strings"
)

func Generate(src any, dest any) (string, error) {
	var config string
	var err error

	config, err = GenerateDiff(src, dest, config)
	if err != nil {
		return "", err
	}
	t := reflect.TypeOf(dest)
	parent, err := GenerateTag(dest, "parent")
	if err != nil {
		return "", err
	}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("preg") != "" {
			config, _, err = GeneratePart(src, dest, i, t, config)
			if err != nil {
				return "", err
			}
		} else if t.Field(i).Tag.Get("parent") == "" {
			cmd, err := GenerateField(t.Field(i), reflect.ValueOf(&dest).Elem().Elem().Field(i), false)
			if err != nil {
				return "", err
			}
			if cmd == "" {
				continue
			}
			cmd = cmd + "\n"
			config = config + cmd
		}
	}
	if parent != "" {
		configSplit := strings.Split(config, "\n")
		config = ""
		for k, s := range configSplit {
			if strings.TrimSpace(s) == "" {
				continue
			}
			if k == len(configSplit)-1 {
				config = config + " " + s
				continue
			}
			config = config + " " + s + "\n"
		}
	}
	config = config + "!"
	return parent + "\n" + config + "\n", nil
}

func GenerateParent(src any) (string, error) {
	parent := ""
	srcType := reflect.TypeOf(src)
	for i := 0; i < srcType.NumField(); i++ {
		if srcType.Field(i).Tag.Get("parent") != "true" {
			continue
		}
		var err error
		parent, err = GenerateField(srcType.Field(i), reflect.ValueOf(&src).Elem().Elem().Field(i), false)
		if err != nil {
			return "", err
		}
	}
	return parent, nil
}

func GenerateDiff(src any, dest any, config string) (string, error) {
	diff, err := utils.DiffFields(dest, src)
	if err != nil {
		return config, err
	}
	for _, s := range diff {
		cmd, err := GenerateFieldByPath(src, s, true)
		if err != nil {
			return "", err
		}
		config = config + cmd + "\n"
	}
	for i := 0; i < reflect.TypeOf(src).NumField(); i++ {
		if reflect.TypeOf(src).Field(i).Tag.Get("preg") == "" {
			continue
		}
		if reflect.TypeOf(src).Field(i).Type.Kind() == reflect.Slice {
			for j := 0; j < reflect.ValueOf(src).Field(i).Len(); j++ {
				parentKey, err := GenerateParent(reflect.ValueOf(src).Field(i).Index(j).Interface())
				if err != nil {
					return "", err
				}
				if parentKey == "" {
					continue
				}
				parentFound := false
				for k := 0; k < reflect.ValueOf(dest).Field(i).Len(); k++ {
					destKey, err := GenerateParent(reflect.ValueOf(dest).Field(i).Index(k).Interface())
					if err != nil {
						continue
					}
					if destKey == "" {
						continue
					}
					if parentKey == destKey {
						parentFound = true
					}
				}
				if !parentFound {
					config = config + "no " + parentKey + "\n"
				}
			}
		}
	}
	return config, nil
}

func GenerateTag(parsed any, tag string) (string, error) {
	t := reflect.TypeOf(parsed)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get(tag) != "" {
			cmd, err := GenerateField(t.Field(i), reflect.ValueOf(&parsed).Elem().Elem().Field(i), false)
			if err != nil {
				return "", err
			}
			return cmd, nil
		}
	}
	return "", nil
}

func GeneratePart(src any, dest any, i int, t reflect.Type, config string) (string, string, error) {
	if reflect.ValueOf(dest).Field(i).Kind() == reflect.Slice {
		for j := 0; j < reflect.ValueOf(dest).Field(i).Len(); j++ {
			subConfig, _, err := GenerateSubConfig(src, dest, i, j, t)
			if err != nil {
				return "", "", err
			}
			for k, s := range strings.Split(subConfig, "\n") {
				if k == len(strings.Split(subConfig, "\n"))-1 {
					config += strings.Trim(s, " ")
					continue
				}
				config += s + "\n"
			}
		}
	} else {
		subConfig, err := Generate(reflect.ValueOf(src).Field(i).Interface(), reflect.ValueOf(dest).Field(i).Interface())
		if err != nil {
			return "", "", err
		}
		subConfig = strings.Join(strings.Split(subConfig, "\n")[:len(strings.Split(subConfig, "\n"))-2], "\n")
		subConfig += "\n" + utils.GetTag(t.Field(i).Type, "exit") + "\n!\n"
		config += subConfig
	}
	return config, "", nil
}

func GenerateSubConfig(src any, dest any, fieldIndex int, destIndex int, t reflect.Type) (string, string, error) {
	destKey, err := GenerateTag(reflect.ValueOf(dest).Field(fieldIndex).Index(destIndex).Interface(), "parent")
	if err != nil {
		return "", "", err
	}
	srcIndex := -1
	for k := 0; k < reflect.ValueOf(src).Field(fieldIndex).Len(); k++ {
		parentKey, err := GenerateParent(reflect.ValueOf(src).Field(fieldIndex).Index(k).Interface())
		if err != nil {
			return "", "", err
		}
		if parentKey == "" {
			continue
		}
		if destKey == parentKey {
			srcIndex = k
			break
		}
	}
	if srcIndex == -1 {
		srcIndex = destIndex
		src = dest
	}
	subConfig, err := Generate(reflect.ValueOf(src).Field(fieldIndex).Index(srcIndex).Interface(), reflect.ValueOf(dest).Field(fieldIndex).Index(destIndex).Interface())
	if err != nil {
		return "", "", err
	}
	subConfig = strings.Join(strings.Split(subConfig, "\n")[:len(strings.Split(subConfig, "\n"))-2], "\n")
	subConfig += "\n" + utils.GetTag(t.Field(fieldIndex).Type, "exit") + "\n!\n"
	return subConfig, "", nil
}

func GenerateFieldByPath(parsed any, path string, reverse bool) (string, error) {
	val, structField, err := utils.GetValueAndField(parsed, path)
	if err != nil {
		return "", err
	}

	if structField.Type.Kind() == reflect.Bool {
		return "", nil
	}
	
	cmd, err := GenerateField(*structField, *val, reverse)
	if err != nil {
		return "", err
	}

	return cmd, nil
}

func GenerateField(field reflect.StructField, value reflect.Value, reverse bool) (string, error) {
	tag := field.Tag.Get("cmd")
	if tag != "" {
		defaultval := field.Tag.Get("default")
		cmd := generateCMD(value, tag, defaultval, reverse)
		return cmd, nil
	}
	return "", nil
}

func generateCMD(field reflect.Value, tag, defaultval string, reverse bool) string {
	var cmd string
	switch field.Type().Kind() {
	case reflect.Struct:
		cmds := tag
		for i := 0; i < field.NumField(); i++ {
			defaultvali := reflect.TypeOf(field.Interface()).Field(i).Tag.Get("default")
			gcmd := generateCMD(field.Field(i), reflect.TypeOf(field.Interface()).Field(i).Tag.Get("cmd"), defaultvali, false)
			if reverse {
				gcmd = "no " + gcmd
			}
			cmds = cmds + gcmd
		}
		cmd = cmds
	case reflect.String:
		value := field.String()
		if value == "" {
			return ""
		}
		cmd = fmt.Sprintf(tag, value)
		if reverse {
			cmd = "no " + cmd
		}
	case reflect.Int:
		value := field.Int()
		if value == 0 {
			return ""
		}
		cmd = fmt.Sprintf(tag, value)
		if reverse {
			cmd = "no " + cmd
		}
	case reflect.Bool:
		value := field.Bool()
		if reverse && value {
			return "no " + tag
		} else if reverse {
			return tag
		}
		concreteDefault, _ := strconv.ParseBool(defaultval)
		if value == concreteDefault {
			return ""
		}

		if !value && defaultval != "" && !reverse {
			return "no " + tag
		} else if !value {
			return ""
		}
		cmd = tag
	case reflect.Float64:
		value := field.Float()
		if value == 0.0 {
			return ""
		}
		cmd = fmt.Sprintf(tag, value)
		if reverse {
			cmd = "no " + cmd
		}
	case reflect.Slice:
		switch field.Type().String() {
		case "[]string":
			slice, _ := field.Interface().([]string)
			if len(slice) == 0 {
				return ""
			}
			cmds := ""
			for i2, s := range slice {
				cmd1 := fmt.Sprintf(tag, s)
				if reverse {
					cmd1 = "no " + cmd1
				}
				if i2 == 0 {
					cmds = cmds + cmd1
				} else {
					cmds = cmds + "\n" + cmd1
				}
			}
			cmd = cmd + cmds
		case "[]int":
			slice, _ := field.Interface().([]int)
			if len(slice) == 0 {
				return ""
			}
			var sliceStr []string
			for _, s := range slice {
				text := strconv.Itoa(s)
				sliceStr = append(sliceStr, text)
			}
			cmd = fmt.Sprintf(tag, strings.Join(sliceStr, ","))
			if reverse {
				cmd = "no " + cmd
			}
		default:
			if field.Type().Kind() == reflect.Slice && tag != "" {
				for i2 := 0; i2 < field.Len(); i2++ {
					cmds := tag
					for i1 := 0; i1 < field.Type().Elem().NumField(); i1++ {
						defaultvali := field.Index(0).Type().Field(i1).Tag.Get("default")
						gcmd := generateCMD(field.Index(i2).Field(i1), field.Index(0).Type().Field(i1).Tag.Get("cmd"), defaultvali, false)
						cmds = cmds + gcmd
					}
					if reverse {
						cmds = "no " + cmds
					}
					if i2 == 0 {
						cmd = cmd + "" + cmds
					} else {
						cmd = cmd + "\n" + "" + cmds
					}
				}
			}
		}
	default:
		panic(field.Type().Kind().String() + " not implemented")
	}

	return cmd
}
