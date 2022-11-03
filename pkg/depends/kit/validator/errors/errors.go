package errors

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

type SyntaxError struct{ msg string }

func (e *SyntaxError) Error() string { return e.msg }

func NewSyntaxError(format string, args ...interface{}) *SyntaxError {
	return &SyntaxError{msg: "invalid syntax: " + fmt.Sprintf(format, args...)}
}

type UnsupportedTypeError struct{ msg string }

func (e *UnsupportedTypeError) Error() string { return e.msg }

func NewUnsupportedTypeError(typ string, rule string, msgs ...string) *UnsupportedTypeError {
	b := &strings.Builder{}
	b.WriteString("unsupported type: rule[" + rule + "] type[" + typ)
	for i, msg := range msgs {
		if i == 0 {
			b.WriteString(": ")
		} else {
			b.WriteString("; ")
		}
		b.WriteString(msg)
	}
	b.WriteString("]")
	return &UnsupportedTypeError{msg: b.String()}
}

type MissingRequiredFieldError struct{}

func (MissingRequiredFieldError) Error() string { return "missing required field" }

type NotMatchError struct {
	Target  string
	Current interface{}
	Pattern *regexp.Regexp
}

func (err *NotMatchError) Error() string {
	return fmt.Sprintf("%s %s not match %v", err.Target, err.Pattern, err.Current)
}

type MultipleOfError struct {
	Target     string
	Current    interface{}
	MultipleOf interface{}
}

func (e *MultipleOfError) Error() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(e.Target)
	buf.WriteString(fmt.Sprintf(" should be multiple of %v", e.MultipleOf))
	buf.WriteString(fmt.Sprintf(", but got invalid value %v", e.Current))
	return buf.String()
}

type NotInEnumError struct {
	Target  string
	Current interface{}
	Enums   []interface{}
}

func (e *NotInEnumError) Error() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(e.Target)
	buf.WriteString(" should be one of ")

	for i, v := range e.Enums {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%v", v))
	}

	buf.WriteString(fmt.Sprintf(", but got invalid value %v", e.Current))

	return buf.String()
}

type OutOfRangeError struct {
	Target           string
	Current          interface{}
	Minimum          interface{}
	Maximum          interface{}
	ExclusiveMaximum bool
	ExclusiveMinimum bool
}

func (e *OutOfRangeError) Error() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(e.Target)
	buf.WriteString(" should be")

	if e.Minimum != nil {
		buf.WriteString(" larger")
		if e.ExclusiveMinimum {
			buf.WriteString(" or equal")
		}

		buf.WriteString(fmt.Sprintf(" than %v",
			reflectx.Indirect(reflect.ValueOf(e.Minimum)).Interface()))
	}

	if e.Maximum != nil {
		if e.Minimum != nil {
			buf.WriteString(" and")
		}

		buf.WriteString(" less")
		if e.ExclusiveMaximum {
			buf.WriteString(" or equal")
		}

		buf.WriteString(fmt.Sprintf(" than %v",
			reflectx.Indirect(reflect.ValueOf(e.Maximum)).Interface()))
	}

	buf.WriteString(fmt.Sprintf(", but got invalid value %v", e.Current))

	return buf.String()
}
