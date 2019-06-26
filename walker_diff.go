package dbgu

import (
	"fmt"
	"reflect"
)

type oldnew struct {
	name   string
	oldVal string
	newVal string
}

type diffWalkFunc func(name []string, fa, fb reflect.Value)

type diffWalker struct{}

func (walk diffWalker) Fields(prefix []string, a, b reflect.Value, fn diffWalkFunc) {
	count := 0
	for i := 0; i < b.NumField(); i++ {
		fname := b.Type().Field(i).Name
		if fname[0] < 'A' || fname[0] > 'Z' {
			continue
		}
		var fa reflect.Value
		if a.IsValid() {
			fa = a.Field(i)
		}
		fb := b.Field(i)
		walk.Compare(append(prefix, "."+fname), fa, fb, fn)
		count++
	}
	if count == 0 {
		if !eq(a, b) {
			fn(prefix, a, b)
		}
	}
}
func (walk diffWalker) Map(prefix []string, a, b reflect.Value, fn diffWalkFunc) {
	searched := map[string]struct{ a, b reflect.Value }{}
	if a.IsValid() {
		if a.IsNil() {
			a = reflect.Indirect(reflect.New(a.Type()))
		}
		for _, k := range a.MapKeys() {
			sk := fmt.Sprint(k)
			v := searched[sk]
			v.a = k
			searched[sk] = v
		}
	}
	if b.IsValid() {
		if b.IsNil() {
			b = reflect.Indirect(reflect.New(b.Type()))
		}
		for _, k := range b.MapKeys() {
			sk := fmt.Sprint(k)
			v := searched[sk]
			v.b = k
			searched[sk] = v
		}
	}

	// combined
	for k, v := range searched {
		fname := fmt.Sprintf("[%#v]", k)
		var fa, fb reflect.Value
		if v.a.IsValid() {
			fa = a.MapIndex(v.a)
		}
		if v.b.IsValid() {
			fb = b.MapIndex(v.b)
			if fb.IsValid() && fb.Kind() == reflect.Interface {
				fb = fb.Elem()
			}
		}
		walk.Compare(append(prefix, fname), fa, fb, fn)
	}
}

// Takes into account of orderer items
func (walk diffWalker) Slice(prefix []string, a, b reflect.Value, fn diffWalkFunc) {
	searched := map[int]struct{ a, b reflect.Value }{}
	if a.IsValid() {
		for i := 0; i < a.Len(); i++ {
			v := searched[i]
			v.a = a.Index(i)
			searched[i] = v
		}
	}

	if b.IsValid() {
		for i := 0; i < b.Len(); i++ {
			v := searched[i]
			v.b = b.Index(i)
			searched[i] = v
		}
	}
	// combined
	for k, v := range searched {
		fname := fmt.Sprintf("[%#v]", k)
		var fa, fb reflect.Value
		if v.a.IsValid() {
			fa = a.Index(k)
		}
		if v.b.IsValid() {
			fb = b.Index(k)
			if fb.IsValid() && fb.Kind() == reflect.Interface {
				fb = fb.Elem()
			}
		}
		walk.Compare(append(prefix, fname), fa, fb, fn)
	}

}
func (walk diffWalker) Compare(prefix []string, a, b reflect.Value, fn diffWalkFunc) {
	// struct
	if bt := reflect.Indirect(b); bt.Kind() == reflect.Struct {
		walk.Fields(prefix, reflect.Indirect(a), bt, fn)
		return
	}
	// map
	if bt := reflect.Indirect(b); bt.Kind() == reflect.Map {
		walk.Map(prefix, reflect.Indirect(a), bt, fn)
		return
	}
	// slice
	if bt := reflect.Indirect(b); bt.Kind() == reflect.Slice {
		walk.Slice(prefix, reflect.Indirect(a), bt, fn)
		return
	}
	var ia, ib interface{}
	if a.IsValid() && (a.Kind() != reflect.Ptr || !a.IsNil()) {
		ia = a.Interface()
	}
	if b.IsValid() && (b.Kind() != reflect.Ptr || !b.IsNil()) {
		ib = b.Interface()
	}

	if !eq(ia, ib) {
		fn(prefix, a, b)
	}
}
func eq(a, b interface{}) bool {
	if a == nil || b == nil {
		return reflect.DeepEqual(a, b)
	}
	// String comparison avoids deep comparing of fields like time which
	// sometimes contains differences while being the 'same'
	ia := reflect.Indirect(reflect.ValueOf(a)).Interface()
	ib := reflect.Indirect(reflect.ValueOf(b)).Interface()
	return fmt.Sprint(ia) == fmt.Sprint(ib)
}
