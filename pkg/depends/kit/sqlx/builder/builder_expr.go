package builder

import (
	"bytes"
	"context"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

type SqlExpr interface {
	IsNil() bool
	Ex(ctx context.Context) *Ex
}

func IsNilExpr(e SqlExpr) bool { return e == nil || e.IsNil() }

func Expr(query string, args ...interface{}) *Ex {
	if query != "" {
		return &Ex{b: *bytes.NewBufferString(query), args: args}
	}
	return &Ex{args: args}
}

func ExprBy(build func(context.Context) *Ex) SqlExpr {
	return &by{build: build}
}

func ExactlyExpr(query string, args ...interface{}) *Ex {
	if query != "" {
		return &Ex{b: *(bytes.NewBufferString(query)), args: args, exactly: true}
	}
	return &Ex{args: args, exactly: true}
}

type ValueExpr interface {
	ValueEx() string
}

type Ex struct {
	b       bytes.Buffer
	args    []interface{}
	err     error
	exactly bool
}

func (e *Ex) IsNil() bool { return e == nil || e.b.Len() == 0 }

func (e *Ex) Query() string {
	if e == nil {
		return ""
	}
	return e.b.String()
}

func (e *Ex) Args() []interface{} {
	if e == nil || len(e.args) == 0 {
		return nil
	}
	return e.args
}

func (e *Ex) Err() error { return e.err }

func (e *Ex) AppendArgs(args ...interface{}) { e.args = append(e.args, args...) }

func (e *Ex) ArgsLen() int { return len(e.args) }

func (e *Ex) SetExactly(exactly bool) { e.exactly = exactly }

func (e *Ex) Grow(n int) {
	if n > 0 && cap(e.args)-len(e.args) < n {
		args := make([]interface{}, len(e.args), 2*cap(e.args)+n)
		copy(args, e.args)
		e.args = args
	}
}

func (e *Ex) WriteExpr(expr SqlExpr) {
	if !IsNilExpr(expr) {
		e.WriteHolder(0)
		e.AppendArgs(expr)
	}
}

func (e *Ex) WriteHolder(idx int) {
	if idx > 0 {
		e.b.WriteByte(',')
	}
	e.b.WriteByte('?')
}

func (e *Ex) WriteQuery(query string) { _, _ = e.b.WriteString(query) }

func (e *Ex) WriteQueryByte(b byte) { _ = e.b.WriteByte(b) }

func (e *Ex) WriteQueryRune(r rune) { _, _ = e.b.WriteRune(r) }

func (e *Ex) WriteGroup(f func(e *Ex)) {
	e.WriteQueryByte('(')
	f(e)
	e.WriteQueryByte(')')
}

func (e *Ex) WriteComments(comments []byte) {
	e.WriteQuery("/* ")
	_, _ = e.b.Write(comments)
	e.WriteQuery(" */")
}

func (e *Ex) WriteEnd() { e.WriteQueryByte(';') }

func (e *Ex) Ex(ctx context.Context) *Ex {
	if e.IsNil() {
		return nil
	}
	args, argc := e.args, len(e.args)

	er := Expr("")
	er.Grow(argc)

	query := e.Query()
	if e.exactly {
		er.WriteQuery(query)
		er.AppendArgs(args...)
		er.exactly = true
		return er
	}

	if shouldResolve := preprocessArgs(args); !shouldResolve {
		er.WriteQuery(query)
		er.AppendArgs(args...)
		er.SetExactly(true)
		return er
	}

	qc := 0
	for i := range query {
		switch c := query[i]; c {
		default:
			er.WriteQueryByte(c)
		case '?':
			if qc >= argc {
				panic(fmt.Errorf("missing arg %d of %s", qc, query))
			}
			switch arg := args[qc].(type) {
			case SqlExpr:
				if !IsNilExpr(arg) {
					sub := arg.Ex(ctx)
					if sub != er && !IsNilExpr(sub) {
						er.WriteQuery(sub.Query())
						er.AppendArgs(sub.Args()...)
					}
				}
			default:
				er.WriteHolder(0)
				er.AppendArgs(arg)
			}
			qc++
		}
	}
	er.SetExactly(true)
	return er
}

func preprocessArgs(args []interface{}) bool {
	shouldResolve := false

	sliceArgEx := func(vs []interface{}) *Ex {
		if n := len(vs); n > 0 {
			return ExactlyExpr(strings.Repeat(",?", n)[1:], vs...)
		}
		return ExactlyExpr("")
	}

	for i := range args {
		switch arg := args[i].(type) {
		case ValueExpr:
			args[i] = ExactlyExpr(arg.ValueEx(), arg)
			shouldResolve = true
		case SqlExpr:
			shouldResolve = true
		case driver.Valuer:
		case []interface{}:
			args[i] = sliceArgEx(arg)
			shouldResolve = true
		default:
			if t := reflect.TypeOf(arg); t.Kind() == reflect.Slice {
				if !reflectx.IsBytes(arg) {
					args[i] = sliceArgEx(toInterfaceSlice(arg))
					shouldResolve = true
				}
			}
		}
	}
	return shouldResolve
}

func toInterfaceSlice(v interface{}) []interface{} {
	switch x := (v).(type) {
	case []bool:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []string:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []float32:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []float64:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []int:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []int8:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []int16:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []int32:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []int64:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []uint:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []uint8:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []uint16:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []uint32:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []uint64:
		values := make([]interface{}, len(x))
		for i := range values {
			values[i] = x[i]
		}
		return values
	case []interface{}:
		return x
	}
	sliceRv := reflect.ValueOf(v)
	values := make([]interface{}, sliceRv.Len())
	for i := range values {
		values[i] = sliceRv.Index(i).Interface()
	}
	return values
}

func Multi(es ...SqlExpr) SqlExpr {
	return MultiWith(" ", es...)
}

func MultiWith(connector string, es ...SqlExpr) SqlExpr {
	return ExprBy(func(ctx context.Context) *Ex {
		e := Expr("")
		e.Grow(len(es))
		for i := range es {
			if i != 0 {
				e.WriteQuery(connector)
			}
			e.WriteExpr(es[i])
		}
		return e.Ex(ctx)
	})
}

func ResolveExpr(v interface{}) *Ex {
	return ResolveExprContext(context.Background(), v)
}

func ResolveExprContext(ctx context.Context, v interface{}) *Ex {
	switch e := v.(type) {
	case nil:
		return nil
	case SqlExpr:
		if IsNilExpr(e) {
			return nil
		}
		return e.Ex(ctx)
	}
	return nil
}
