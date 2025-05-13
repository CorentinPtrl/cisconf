package utils

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type DiffReporter struct {
	path   cmp.Path
	diffs  []string
	fields []string
}

func printPath(path cmp.Path) string {
	var result string
	for i := 1; i < len(path); i++ {
		if path.Index(i).Type().Kind() == reflect.Slice {
			if path.Index(i+1).String() == "<nil>" {
				result += path.Index(i).String()
				i++
				continue
			}
			result += fmt.Sprintf("%s[%d]", path.Index(i).String(), addressCountIndex(path.Index(i+1).String()))
			i++
		} else {
			result += path.Index(i).String()
		}
	}
	return result
}

func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	vx, vy := ps.Values()
	if vx.Kind() == reflect.Slice && vy.Kind() == reflect.Slice {
		fvx, _ := r.path.Index(0).Values()
		if isTagInPath(fvx, r.path, "preg") {
			r.path = append(r.path, ps)
			return
		}

		lenX, lenY := vx.Len(), vy.Len()
		if lenX != lenY && lenY > lenX {
			for i := lenX; i < lenY; i++ {
				r.fields = append(r.fields, fmt.Sprintf("%s[%d]", printPath(append(r.path, ps)), i))
			}
		}
	}
	r.path = append(r.path, ps)
}

func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		fvx, _ := r.path.Index(0).Values()
		if vy.IsValid() {

			if r.path.Index(1).Type().Kind() == reflect.Slice {
				_, dest, err := GetValueAndField(fvx.Interface(), fmt.Sprintf("%s[%d]", r.path.Index(1).String(), addressCountIndex(r.path.Index(2).String())))
				if err != nil {
					return
				}
				if dest.Tag.Get("preg") != "" {
					return
				}
				r.fields = append(r.fields, fmt.Sprintf("%s[%d]", r.path.Index(1).String(), addressCountIndex(r.path.Index(2).String())))
			} else {
				_, dest, err := GetValueAndField(fvx.Interface(), r.path.Index(1).String())
				if err != nil {
					return
				}
				if dest.Tag.Get("preg") != "" {
					return
				}
				r.fields = append(r.fields, r.path.Index(1).String())
			}
			r.diffs = append(r.diffs, fmt.Sprintf("%#v:\n\t-: %+v\n\t+: %+v\n", r.path, vx, vy))
		}
	}
}

func isTagInPath(val reflect.Value, path cmp.Path, tag string) bool {
	pathStr := printPath(path)
	if pathStr == "" {
		return false
	}
	splitPath := strings.Split(pathStr, ".")
	for i := 0; i < len(splitPath); i++ {
		_, dest, err := GetValueAndField(val.Interface(), strings.Join(splitPath[i:], "."))
		if err != nil {
			continue
		}
		if dest.Tag.Get(tag) != "" {
			return true
		}
	}
	return false
}

func addressCountIndex(addr string) int {
	r := regexp.MustCompile(`\d+`)
	m := r.FindStringSubmatch(addr)

	if len(m) > 0 {
		i, _ := strconv.Atoi(m[0])

		return i
	}

	return -1
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r DiffReporter) Fields() []string {
	return r.fields
}

func (r *DiffReporter) String() string {
	return strings.Join(r.fields, "\n")
}
