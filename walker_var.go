package dbgu

import (
	"fmt"
	"reflect"
)

type fieldWalkFunc func(name []string, val reflect.Value)

type fieldWalker struct{}

func (walk fieldWalker) Var(prefix []string, a reflect.Value, fn fieldWalkFunc) {
	ra := reflect.Indirect(a)
	switch ra.Kind() {
	case reflect.Struct:
		walk.Fields(prefix, ra, fn)
	case reflect.Slice:
		walk.Slice(prefix, ra, fn)
	case reflect.Map:
		walk.Map(prefix, ra, fn)
	default:
		fn(prefix, a)
	}
}

func (walk fieldWalker) Fields(prefix []string, a reflect.Value, fn fieldWalkFunc) {
	count := 0
	for i := 0; i < a.NumField(); i++ {
		fname := a.Type().Field(i).Name
		if fname[0] < 'A' || fname[0] > 'Z' {
			continue
		}

		if len(prefix) > 0 {
			fname = "." + fname
		}
		count++
		walk.Var(append(prefix, fname), a.Field(i), fn)
	}
	// It is a struct with no exported fields like (time.Time), we print it
	if count == 0 {
		fn(prefix, a)
	}
}
func (walk fieldWalker) Slice(prefix []string, a reflect.Value, fn fieldWalkFunc) {
	for i := 0; i < a.Len(); i++ {
		walk.Var(append(prefix, fmt.Sprint("[", i, "]")), a.Index(i), fn)
	}
}
func (walk fieldWalker) Map(prefix []string, a reflect.Value, fn fieldWalkFunc) {
	for _, k := range a.MapKeys() {
		walk.Var(append(prefix, fmt.Sprint("[", value(k), "]")), a.MapIndex(k), fn)
	}
}
