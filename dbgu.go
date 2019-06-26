package dbgu

import (
	"fmt"
	"reflect"
)

// Opt option type
type Opt int

// optionConstants
const (
	OptColor = iota
	OptNoColor
	OptJSON
	OptNoJSON
)

// Dump initiator
func Dump(opts ...Opt) Dumper {
	// Defaults
	enableColor := true
	enableJSON := false
	for _, o := range opts {
		switch o {
		case OptColor:
			enableColor = true
		case OptNoColor:
			enableColor = false
		case OptJSON:
			enableJSON = true
		case OptNoJSON:
			enableJSON = false
		}
	}
	return Dumper{enableJSON, enableColor}
}

func value(v reflect.Value) string {
	if !v.IsValid() {
		return "nil"
	}
	if v.Kind() == reflect.Struct {
		return fmt.Sprintf("%v", v.Interface())
	}
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		return fmt.Sprintf("&(%#v)", v.Elem().Interface())
	}
	return fmt.Sprintf("%#v", v.Interface())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ANSI Terminal colors

type colors bool

func (c colors) hl(s interface{}) string {
	if !c {
		return fmt.Sprint(s)
	}
	return fmt.Sprint("\033[01;37m", s, "\033[0m")
}
func (c colors) red(s interface{}) string {
	if !c {
		return fmt.Sprint(s)
	}
	return fmt.Sprint("\033[01;31m", s, "\033[0m")
}
func (c colors) green(s interface{}) string {
	if !c {
		return fmt.Sprint(s)
	}
	return fmt.Sprint("\033[01;32m", s, "\033[0m")
}
func (c colors) blue(s interface{}) string {
	if !c {
		return fmt.Sprint(s)
	}
	return fmt.Sprint("\033[01;34m", s, "\033[0m")
}
