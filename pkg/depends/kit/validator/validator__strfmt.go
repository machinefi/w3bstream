package validator

import (
	"context"
	"reflect"
	"regexp"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
)

type ValidateFunc func(interface{}) error

type StrFmt struct {
	names []string
	fn    ValidateFunc
}

func NewStrfmtValidator(f ValidateFunc, name string, aliases ...string) *StrFmt {
	return &StrFmt{append([]string{name}, aliases...), f}
}

func NewRegexpStrfmtValidator(expr string, name string, aliases ...string) *StrFmt {
	regexpStrfmtPattern := regexp.MustCompile(expr)
	f := func(v interface{}) error {
		if !regexpStrfmtPattern.MatchString(v.(string)) {
			return &errors.NotMatchError{
				Target:  name,
				Current: v,
				Pattern: regexpStrfmtPattern,
			}
		}
		return nil
	}
	return NewStrfmtValidator(f, name, aliases...)
}

func (v *StrFmt) String() string { return "@" + v.names[0] }

func (v *StrFmt) Names() []string { return v.names }

func (v StrFmt) New(_ context.Context, r *Rule) (Validator, error) {
	return &v, v.TypeCheck(r)
}

func (v *StrFmt) TypeCheck(r *Rule) error {
	if r.Type.Kind() == reflect.String {
		return nil
	}
	return errors.NewUnsupportedTypeError(r.String(), v.String())
}

func (v *StrFmt) Validate(val interface{}) error {
	if rv, ok := val.(reflect.Value); ok && rv.CanInterface() {
		val = rv.Interface()
	}
	str, ok := val.(string)
	if !ok {
		return errors.NewUnsupportedTypeError(
			reflect.TypeOf(val).String(), v.String(),
		)
	}
	return v.fn(str)
}
