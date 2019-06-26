package dbgu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/alecthomas/chroma/quick"
)

// Dumper settings
type Dumper struct {
	JSON  bool
	Color bool
}

// Var a Var (struct)
func (d Dumper) Var(out io.Writer, v interface{}) {
	if d.JSON {
		d.dumpJSON(out, v)
		return
	}
	d.dump(out, v)
}

// VarString returns a string
func (d Dumper) VarString(v interface{}) string {
	out := &bytes.Buffer{}
	d.Var(out, v)
	return out.String()
}

//Diff Dumps the different fields from b to a
func (d Dumper) Diff(out io.Writer, a, b interface{}) {
	d.dumpDiff(out, a, b)
}

// DiffString returns a string
func (d Dumper) DiffString(a, b interface{}) string {
	out := &bytes.Buffer{}
	d.Diff(out, a, b)
	return out.String()
}

func (d Dumper) dumpJSON(out io.Writer, v interface{}) {
	if d.Color == false {
		enc := json.NewEncoder(out)
		enc.SetIndent("  ", "  ")
		enc.Encode(v)
		return
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetIndent("  ", "  ")
	enc.Encode(v)

	err := quick.Highlight(out, buf.String(), "json", "terminal16m", "monokai")
	if err != nil {
		fmt.Fprint(out, "highlight error", err)
	}
	return
}

func (d Dumper) dump(out io.Writer, v interface{}) {
	c := colors(d.Color)
	val := reflect.Indirect(reflect.ValueOf(v))

	walk := fieldWalker{}

	fmt.Fprintln(out)

	walk.Var([]string{}, val, func(names []string, val reflect.Value) {
		fmt.Fprintf(out, "%s: %v\n", c.hl(strings.Join(names, "")), c.green(value(val)))
	})
}

func (d Dumper) dumpDiff(out io.Writer, a, b interface{}) {
	// Just to be safe from invalid types here
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprint(out, "diff error: ", err)
			debug.PrintStack()
		}
	}()
	diff := []oldnew{}

	va := reflect.Indirect(reflect.ValueOf(a))
	vb := reflect.Indirect(reflect.ValueOf(b))

	pad := 0

	walk := diffWalker{}
	walk.Compare([]string{}, va, vb, func(name []string, fa, fb reflect.Value) {
		fullName := strings.Trim(strings.Join(name, ""), ".")
		if pad < len(fullName) {
			pad = len(fullName)
		}
		diff = append(diff, oldnew{fullName, value(fa), value(fb)})
	})
	c := colors(d.Color)
	// Render
	if len(diff) == 0 {
		fmt.Fprint(out, "\n", c.hl("no difference"))
	}

	for _, v := range diff {
		fmt.Fprintf(out, "\n%s From %v <- %v",
			c.hl(fmt.Sprintf("%-*s", pad, v.name)),
			c.red(v.oldVal),
			c.green(v.newVal),
		)
	}
	fmt.Fprintln(out)
}
