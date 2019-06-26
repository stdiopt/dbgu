package dbgu

import (
	"fmt"
	"reflect"
)

// TODO: walk through maps

type fieldWalkFunc func(name []string, val reflect.Value)

type fieldWalker struct{}

func (walk fieldWalker) Var(prefix []string, a reflect.Value, fn fieldWalkFunc) {
	ra := reflect.Indirect(a)
	if ra.Kind() == reflect.Struct {
		walk.Fields(prefix, ra, fn)
		return
	}
	if ra.Kind() == reflect.Slice {
		walk.Slice(prefix, ra, fn)
		return
	}
	fn(prefix, a)
}

func (walk fieldWalker) Fields(prefix []string, a reflect.Value, fn fieldWalkFunc) {
	for i := 0; i < a.NumField(); i++ {
		fname := a.Type().Field(i).Name
		if fname[0] < 'A' || fname[0] > 'Z' {
			continue
		}

		if len(prefix) > 0 {
			fname = "." + fname
		}
		walk.Var(append(prefix, fname), a.Field(i), fn)
	}
}
func (walk fieldWalker) Slice(prefix []string, a reflect.Value, fn fieldWalkFunc) {
	for i := 0; i < a.Len(); i++ {
		walk.Var(append(prefix, fmt.Sprint("[", i, "]")), a.Index(i), fn)
	}
}
